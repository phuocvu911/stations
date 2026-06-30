# stations-pathfinder

A command line pathfinder which moves trains across a rail network in the
minimum number of movement turns. Only one train may occupy an intermediate
station at a time, and each track may be used once per turn.

## Usage

```
go run . [path to network map] [start station] [end station] [number of trains]
```

Example:

```
$ go run . maps/london.map waterloo st_pancras 4
T1-victoria T2-euston
T1-st_pancras T2-st_pancras T3-victoria T4-euston
T3-st_pancras T4-st_pancras
```

Each output line is one turn; `T1-victoria` means train T1 moves to
victoria during that turn. All errors are reported on stderr with a
message starting with `Error:` and a non-zero exit code.

## Bonus: -draw

```
go run . -draw maps/london.map waterloo st_pancras 4
```

`-draw` (or `--draw`) renders the network as ASCII art in the subject's
style — an x-axis header, y coordinates down the left, stations as
`X <- name`, tracks drawn with `- | / \` — and then replays the schedule
turn by turn with train positions overlaid on the map:

```
0   1  2  3  4  5  6  7  8  9  10 11
1         X <- waterloo [2 waiting]
2          \
3           \
4            \
5             \
6              \
7               \--X <- victoria [T1]
8                \/
9                /\
10              /  \
11              |   \
12              |    \
13              |     \
14              |      \
15              X <- st_pancras
16               \       \
...
23                      \-------+-X <- euston [T2]
```

`[T1]` marks a train at a station; the start and end stations show
`[n waiting]` / `[n arrived]` tallies. When stdout is a terminal the
frames animate in place (one frame every 0.5s); when piped the frames are
printed one after another. Each frame is preceded by the turn's movement
line (`Turn 2: T1-st_pancras T2-st_pancras ...`).

The flag does not change the default behaviour: without `-draw` the
output is exactly the plain schedule required by the subject. The grid is
true to the station coordinates; maps whose coordinates would not fit on
a screen fall back to a scaled rendering without the axes, and on
overcrowded rows labels are truncated rather than drawn over each other.

## How it works

1. **Parse** the map (`parse.go`): sections, comments, whitespace, and all
   the validation rules (names, coordinates, duplicates, 10k-station limit,
   connections to unknown stations, duplicate/reversed connections, ...).
2. **Find routes** (`flow.go`, `plan.go`): the network is turned into a flow
   graph where every intermediate station is split into an in/out node pair
   of capacity 1, so one unit of flow is one vertex-disjoint route.
   Successive cheapest augmentations (SPFA min-cost flow, which can re-route
   earlier paths through residual edges) yield, for every k, the k disjoint
   routes with the smallest total length.
3. **Pick the best plan** (`plan.go`): a route of length L delivers
   `T - L + 1` trains within T turns, since trains pipeline one block apart.
   For each candidate route set the minimal T moving all trains is computed,
   and the best set wins. The search stops as soon as new routes are too
   long to ever help.
4. **Simulate** (`schedule.go`): each turn, every en-route train advances one
   station and each route admits one new train, printing one line per turn.
   The `-draw` rendering (`draw.go`) replays the same simulation and plots
   it on a coordinate-scaled ASCII map instead.

This handles large cyclic maps gracefully: a 10,000-station map with 30,000
connections schedules 100 trains in well under a second.
