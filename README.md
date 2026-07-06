# stations-pathfinder

A command-line pathfinder that moves trains across a rail network in the
minimum number of movement turns. Only one train may occupy an intermediate
station at a time, and each track may be used once per turn.

## Usage

### Clone the project 
```
git clone https://gitea.kood.tech/hoangphuocvu/pathfinder
cd pathfinder
```

### Run
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

## How it works

1. **Parse** the map (`parse.go`): sections, comments, whitespace, and all
   the validation rules (names, coordinates, duplicates, 10k-station limit,
   connections to unknown stations, duplicate/reversed connections, ...). Station name got indexed for even better performance.
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

## Extras

### Super Advanced Error Handling
In any cases of error, the program prints out `<Time> Error: <explaining what happend and what clause cause err>`

Example:
```
go run . maps/london.map waterloo st_pancras -2

2026/07/05 19:29:47 Error: number of trains (-2) is not a valid positive integer
```
```
go run . maps/london.map waterloo hakaniemi 4

2026/07/05 19:34:18 Error: end station "hakaniemi" does not exist
```
### Suite of Tests
A suite of tests has been created in advance in a make file, covering the cases described in the school's testing tab.How to use: the test case is mark 1-30 based on the order of the testing cases.

For example: you want to test this
```
It finds more than one valid route for 3 trains between waterloo and st_pancras in the London Network Map.
```
And it is task number 3 in the testing tab.
In the terminal:
```
make 3
```
So on and so forth for the remaining test cases.
### Super Fast Performance
Try:
```
go run . maps/big.map station0 station9999 100
```
This is the map with 10,000 stations and 30,000 connections. The program will finish in less than a second.
### Bonus
