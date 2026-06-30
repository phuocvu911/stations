# stations-pathfinder — Data Flow & Implementation Guide

This document explains how the solution works and how to build it yourself,
stage by stage. Each stage says **what goes in, what comes out, why it works,
and how to test it before moving on**. The code references point at the real
files in this repo so you can compare your attempt with a working version.

---

## 1. The big picture

The program is a pipeline of five stages. Data only flows forward — each
stage produces a value the next one consumes, and any stage can abort with an
error.

```
 CLI args                 file on disk
    │                          │
    ▼                          ▼
┌─────────────┐         ┌─────────────┐
│ 1. validate │         │ 2. parse &  │
│    args     │────────▶│   validate  │
│  (main.go)  │  path   │  (parse.go) │
└─────────────┘         └──────┬──────┘
                               │ *Network  (stations, connections)
                               ▼
                        ┌─────────────┐
                        │ 3. find     │   min-cost max-flow on a
                        │   routes    │   node-split graph
                        │  (flow.go,  │
                        │   plan.go)  │
                        └──────┬──────┘
                               │ [][]int   (vertex-disjoint routes)
                               ▼
                        ┌─────────────┐
                        │ 4. choose   │   how many turns? how many
                        │   the plan  │   trains per route?
                        │  (plan.go)  │
                        └──────┬──────┘
                               │ paths + counts
                               ▼
                        ┌─────────────┐
                        │ 5. simulate │   one printed line per turn
                        │ (schedule.go│
                        └──────┬──────┘
                               ▼
                            stdout
```

Errors at any stage go to **stderr** as `Error: ...` and exit code 1
(`main.go` wraps everything in one `run() error` so there is exactly one
place that prints errors — a pattern worth copying).

### The data shapes that flow between stages

| Stage boundary | Type | Meaning |
|---|---|---|
| 1 → 2 | `string`, `int` | map path, start/end names, train count |
| 2 → 3 | `*Network` | `names []string`, `index map[string]int`, `conns [][2]int`, `adj [][]int` |
| 3 → 4 | `[][]int` | each route as a list of station ids, e.g. `[0 1 3]` |
| 4 → 5 | `[][]int` + `[]int` | routes sorted shortest-first + trains assigned to each |
| 5 → out | text | `T1-victoria T2-euston` per turn |

A station is referred to by **integer id** (its index into `names`)
everywhere after parsing. Strings are only used at the boundary (parsing in,
printing out). This makes the graph algorithms simple and fast.

---

## 2. Stage 1 — Arguments (`main.go`)

**In:** `os.Args[1:]` · **Out:** validated `mapPath, start, end, numTrains`

Nothing clever here, but get the rules exact:

- exactly 4 arguments, otherwise the "incorrect number of arguments" error;
- the train count must be a *positive* integer: reject `0`, `-3`, `abc`.

**Checkpoint:** run with too few args, with `0` trains, with `abc` trains.
All three must print an `Error:` line on stderr (test with `2>&1 >/dev/null`
so you only see stderr) and exit non-zero (`echo $?`).

---

## 3. Stage 2 — Parsing (`parse.go`)

**In:** file contents · **Out:** `*Network` or an error

### 3.1 Line cleaning comes first

Every line goes through the same two steps *before* any interpretation:

```
"waterloo  , 3 , 1   # south"  ──strip comment──▶ "waterloo  , 3 , 1   "
                                ──trim space────▶ "waterloo  , 3 , 1"
```

In code: cut at the first `#`, then `strings.TrimSpace`. Blank results are
skipped. Doing this in one place means comments and whitespace are handled
identically everywhere — section headers, stations, connections.

### 3.2 A tiny state machine

The parser keeps a `section` variable: `"" → "stations" → "connections"`.
A line is interpreted according to the *current* section:

- `stations:` / `connections:` switch the section (and seeing one twice is
  an error);
- in `stations`: split on `,`, expect exactly 3 fields, trim each;
- in `connections`: split on `-`, expect exactly 2 fields, trim each;
- any data line *before* the first section header is an error.

### 3.3 Validate as you go vs. validate after

Stations are validated immediately (name charset `[a-z0-9_]+`, positive
coordinates, duplicate name, duplicate coordinates, >10 000 stations).

Connections are only *collected* during the scan and resolved **after** the
whole file is read. Why: a connection refers to stations by name, and
resolving afterwards means the parser doesn't care which section came first
and gives the right error ("station does not exist") instead of a confusing
one. When resolving, normalize each pair so `a-b` and `b-a` collide in the
duplicate check:

```go
key := [2]int{min(ia, ib), max(ia, ib)}   // direction-independent identity
```

### 3.4 What `Network` stores

```go
type Network struct {
    names []string          // id -> name        (printing)
    index map[string]int    // name -> id        (lookup)
    conns [][2]int          // each track once   (building the flow graph)
    adj   [][]int           // id -> neighbours  (degree bounds later)
}
```

**Checkpoint:** feed it every broken map you can think of — duplicate
station, duplicate reversed connection, station at same coordinates, unknown
station in a connection, missing `connections:` section, a connection
`a-a`. Then feed it the London example from the subject and dump the parsed
struct. Only move on when all of these behave.

---

## 4. Stage 3 — Finding routes. The heart of the project.

**In:** `*Network`, start, end · **Out:** sets of vertex-disjoint routes

### 4.1 Why "shortest path" is not enough

Take the London map. The shortest route `waterloo→victoria→st_pancras` can
pipeline one train per turn. But with 4 trains it is faster to *also* use
`waterloo→euston→st_pancras`, even though a single train wouldn't need it.
And sometimes routes can't just be added: the obvious shortest route can
**block** all alternatives. Try this map (in the repo's test as "the trap"):

```
a ── b ── c ── z        shortest route: a-b-c-z  (length 3)
│    └─── e ── z
└─── d ── c
```

If train traffic claims `a-b-c-z`, no second disjoint route exists (`d`
dead-ends at the occupied `c`, `e` is unreachable without `b`). But the pair
`a-b-e-z` **and** `a-d-c-z` coexist. A correct algorithm must be able to
*revise its earlier choice*. That is exactly what min-cost flow gives you.

### 4.2 Stations have capacity → split each node in two

Flow networks limit *edges*, but our constraint is "one train per
*station*". The standard trick: split every station `v` into `v_in` and
`v_out` joined by an internal edge whose capacity is the station capacity:

```
                 cap 1 (intermediate)
        ───▶ v_in ────────▶ v_out ───▶
                 cap ∞ (start/end)
```

Each track `u-v` becomes two directed edges `u_out→v_in` and `v_out→u_in`,
capacity 1, **cost 1** (cost = "this route gets one station longer").
Now *one unit of flow from start to end is exactly one route*, and two units
can never share an intermediate station because its internal edge has
capacity 1. (`plan.go`, graph construction at the top of `planMovements`.)

### 4.3 Successive shortest augmentations

Repeat: find the **cheapest** residual path from `start_out` to `end_in`,
push 1 unit along it. After k rounds the flow is the cheapest possible set
of k disjoint routes (total length minimal) — that's the classic min-cost
flow guarantee.

The "revising earlier choices" magic lives in the **residual edges**: when
an edge carries flow, the reverse edge appears with cost −1 ("refund a step
by undoing it"). In the trap map, round 2 finds:

```
a → d → c  →(reverse of b→c, cost −1)→  b → e → z
```

which *cancels* the b→c segment of round 1. Net result after cancellation:
`a-b-e-z` and `a-d-c-z`. Because of those negative residual costs you cannot
use plain Dijkstra; use **SPFA / Bellman-Ford** (`flow.go`, `augmentOne`) —
it tolerates negative edges as long as there are no negative cycles, which
min-cost flow guarantees.

Implementation notes that save you debugging hours (`flow.go`):

- store edges in one flat slice, pairs at indices `2i, 2i+1`, so the reverse
  edge of `e` is always `e^1` — no bookkeeping;
- record `parent[v] = edge index used to reach v` during SPFA, then walk
  backwards from the sink doing `cap--` / `cap++` on the pair.

### 4.4 Reading the routes back out (decomposition)

After k augmentations, the routes exist only implicitly as edge flows. To
extract them (`plan.go`, `decompose`):

1. collect every directed track edge with flow (`cap == 0` since cap was 1);
2. if both `u→v` and `v→u` carry flow, they cancel — drop both;
3. starting from `start`, repeatedly follow flow edges until `end`,
   consuming them. Vertex capacity 1 means every intermediate station has
   exactly one outgoing flow edge, so the walk can't branch wrongly.

**Checkpoint:** print the routes for the London map with k=2 (expect the two
2-hop routes) and the trap map with k=2 (expect `a-b-e-z` + `a-d-c-z`). If
the trap map gives you only one route, your residual edges aren't working.

---

## 5. Stage 4 — Choosing the plan (`plan.go`)

**In:** route sets for k = 1, 2, 3… · **Out:** the best routes + a train
count per route

### 5.1 The pipeline formula — derive it once, use it everywhere

Trains on one route enter one per turn and march in lockstep, always one
block apart (so the station-capacity rule holds automatically). On a route
of length `L` (edges), the train that enters on turn `i` arrives on turn
`L + i − 1`. Therefore:

> **A route of length L delivers `T − L + 1` trains within T turns** (if T ≥ L).

Sanity check against the subject: London, two routes of length 2, T=3 turns
→ each delivers 3−2+1 = 2 trains → 4 trains total in 3 turns. Matches.

### 5.2 Minimal turns for a given route set

For routes of lengths `L₁…Lₖ` and `n` trains, find the smallest `T` with

```
Σ max(0, T − Lᵢ + 1) ≥ n
```

The left side grows with T, so binary search works (`minTurns`). Then assign
each route its capacity `T − Lᵢ + 1` and trim the surplus off the *longest*
routes (`assignTrains`) — a route may legitimately get 0 trains.

### 5.3 Why try every k, and when to stop

More routes is not always better: with 1 train, the shortest route alone is
optimal; the cheapest *pair* of routes might both be longer. So: after each
augmentation, decompose, compute `minTurns`, keep the best
(`planMovements`'s loop). Two cheap bounds keep this fast:

- k never needs to exceed `min(trains, degree(start), degree(end))`;
- augmentation costs only ever increase, so once a new route's length ≥ the
  best turn count so far, no future route can help — stop.

**Checkpoint:** `maps/mixed.map` has routes of lengths 1, 2, 4. With 1 train
the answer is 1 turn (direct route only); with 9 trains it is 5 turns and
the length-4 route gets 0 trains. If your answer uses the long route, your
"try every k / trim surplus" logic is off.

---

## 6. Stage 5 — Simulation & output (`schedule.go`)

**In:** routes (sorted shortest first) + per-route counts · **Out:** stdout

Keep a list of trains, each with `(route, position)`. Per turn:

1. advance every en-route train one station, recording `T<n>-<station>`;
2. then let each route with remaining quota admit **one** new train
   (the start-station track can be used once per turn);
3. print the collected movements as one space-separated line; stop when a
   turn produces no movements.

Numbering falls out for free: trains are created in dispatch order, so
iterating the train list in index order prints `T1 T2 T3…` naturally, and
matches the subject's example exactly.

Why this never violates the rules (worth convincing yourself!):

- *station capacity*: trains on one route stay exactly one block apart, and
  routes share no intermediate stations (vertex-disjoint!);
- *track once per turn*: on one route, consecutive trains use *different*
  tracks in the same turn (train at block i uses track i→i+1, its follower
  uses track i−1→i);
- *every train arrives within T*: that's exactly what `assignTrains`
  guaranteed.

**Checkpoint:** reproduce the subject's 4-train London output, byte for
byte, and check `| wc -l` = 3.

---

## 7. Suggested build order

Build it in this order — every step is runnable and testable on its own:

| Step | Build | Test with |
|---|---|---|
| 1 | arg validation + error plumbing (`run() error`) | wrong arg counts, bad train numbers |
| 2 | line cleaner + section state machine | maps with comments/whitespace |
| 3 | station & connection validation | all 12 error cases from the subject |
| 4 | plain BFS shortest path (temporary!) | 1 train end-to-end — your first real output |
| 5 | simulation for a single route | n trains on a line-shaped map: expect L+n−1 turns |
| 6 | node-split flow graph + SPFA augmentation | trap map must yield 2 routes |
| 7 | decomposition | print the routes as station names |
| 8 | `minTurns` / `assignTrains` / try-every-k | mixed.map cases above |
| 9 | full simulation over multiple routes | the subject's example, byte for byte |
| 10 | stress test | 10k-station generated map (see README) |

Step 4's BFS is throwaway, but it gets the whole pipeline (parse → route →
simulate → print) working before the hard algorithm, which makes debugging
the flow code far easier — when something breaks later, you know which stage
to blame.

### How to know it's correct

Don't eyeball output — verify it. A validator that replays your stdout
against the map and checks every rule (moves use real connections, station
capacity, track reuse, all trains arrive) is ~60 lines and catches what eyes
miss. The one used to validate this repo's outputs checks exactly those four
rules; writing your own is an excellent exercise because it forces you to
restate the movement rules precisely.

---

## 8. Concepts to look up if you want to go deeper

- **Max-flow / min-cut**, augmenting paths (Ford–Fulkerson) — why pushing
  flow one path at a time finds a global optimum.
- **Min-cost max-flow**, successive shortest paths, why residual costs are
  negative, and why SPFA/Bellman-Ford (or Dijkstra with potentials) is used.
- **Menger's theorem** — number of vertex-disjoint paths = minimum vertex
  cut; this is why `min(degree(start), degree(end))` bounds k.
- **Suurballe's algorithm** — the specialised two-disjoint-shortest-paths
  algorithm; our flow approach generalises it to k paths.
- The scheduling part is a tiny case of **makespan minimisation** —
  the pipeline formula is the whole theory you need here, but the term will
  lead you to related problems.
