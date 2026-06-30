package main

import "sort"

// planMovements finds the set of vertex-disjoint routes and the number of
// trains to send down each route which together minimise the number of
// movement turns. It tries every route-set size k (the cheapest k routes by
// total length, via successive shortest augmentations), because adding a
// longer route only pays off when enough trains share the load.
func planMovements(net *Network, start, end, trains int) (paths [][]int, counts []int, ok bool) {
	n := len(net.names)
	g := newFlowGraph(2 * n)
	in := func(v int) int { return 2 * v }
	out := func(v int) int { return 2*v + 1 }
	const bigCap = 1 << 30
	for v := 0; v < n; v++ {
		c := 1
		if v == start || v == end {
			c = bigCap
		}
		g.addEdge(in(v), out(v), c, 0)
	}
	connEdges := make([]connEdge, 0, len(net.conns)*2)
	for _, c := range net.conns {
		u, v := c[0], c[1]
		connEdges = append(connEdges,
			connEdge{u, v, g.addEdge(out(u), in(v), 1, 1)},
			connEdge{v, u, g.addEdge(out(v), in(u), 1, 1)})
	}
	maxRoutes := min(trains, len(net.adj[start]), len(net.adj[end]))
	bestTurns := -1
	var bestPaths [][]int
	for k := 0; k < maxRoutes; k++ {
		cost, found := g.augmentOne(out(start), in(end))
		if !found {
			break
		}
		// Augmentation costs never decrease, so once a new route is at
		// least as long as the best turn count it can never carry a train.
		if bestTurns >= 0 && cost >= bestTurns {
			break
		}
		ps := decompose(g, connEdges, n, start, end)
		lengths := make([]int, len(ps))
		for i, p := range ps {
			lengths[i] = len(p) - 1
		}
		turns := minTurns(lengths, trains)
		if bestTurns < 0 || turns < bestTurns {
			bestTurns = turns
			bestPaths = ps
		}
	}
	if bestTurns < 0 {
		return nil, nil, false
	}
	return bestPaths, assignTrains(bestPaths, trains, bestTurns), true
}

// connEdge remembers which flow edge models each direction of a track.
type connEdge struct{ u, v, idx int }

// decompose extracts the routes carried by the current flow, shortest first.
func decompose(g *flowGraph, connEdges []connEdge, n, start, end int) [][]int {
	used := make(map[[2]int]bool)
	for _, ce := range connEdges {
		if g.edges[ce.idx].cap == 0 {
			used[[2]int{ce.u, ce.v}] = true
		}
	}
	// Opposite directions on the same track cancel out.
	for key := range used {
		rev := [2]int{key[1], key[0]}
		if used[rev] && used[key] {
			delete(used, key)
			delete(used, rev)
		}
	}
	next := make([][]int, n)
	for _, ce := range connEdges {
		if used[[2]int{ce.u, ce.v}] {
			next[ce.u] = append(next[ce.u], ce.v)
		}
	}
	var paths [][]int
	for len(next[start]) > 0 {
		cur := start
		path := []int{start}
		for cur != end && len(next[cur]) > 0 && len(path) <= n {
			nxt := next[cur][0]
			next[cur] = next[cur][1:]
			path = append(path, nxt)
			cur = nxt
		}
		if cur == end {
			paths = append(paths, path)
		}
	}
	sort.SliceStable(paths, func(i, j int) bool { return len(paths[i]) < len(paths[j]) })
	return paths
}

// minTurns is the fewest turns needed to move the given number of trains
// over routes of the given lengths: a route of length L pipelines one train
// per turn, so within T turns it delivers T-L+1 trains.
func minTurns(lengths []int, trains int) int {
	lo, hi := lengths[0], lengths[0]+trains-1
	for lo < hi {
		mid := (lo + hi) / 2
		if capacityWithin(lengths, mid) >= trains {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func capacityWithin(lengths []int, turns int) int {
	total := 0
	for _, l := range lengths {
		if turns >= l {
			total += turns - l + 1
		}
	}
	return total
}

// assignTrains decides how many trains run down each route so that every
// train still arrives within the given number of turns; surplus capacity is
// trimmed from the longest routes first.
func assignTrains(paths [][]int, trains, turns int) []int {
	counts := make([]int, len(paths))
	total := 0
	for i, p := range paths {
		c := turns - (len(p) - 1) + 1
		if c < 0 {
			c = 0
		}
		counts[i] = c
		total += c
	}
	for i := len(counts) - 1; i >= 0 && total > trains; i-- {
		d := min(counts[i], total-trains)
		counts[i] -= d
		total -= d
	}
	return counts
}
