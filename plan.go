package main

import "sort"

const bigCap = 1 << 60

// directedTrack represents a directed track from one station to another, along with the index of the forward edge in the trackMap.tracks slice.
type directedTrack struct {
	from int
	to   int
	idx  int //index of the forward edge in the trackMap.tracks slice
}

// planMovements finds the set of vertex-disjoint routes (non-overlapping) and the number of
// trains to send down each route which together minimise the number of
// movement turns. It tries every route-set size k (the cheapest k routes by
// total length, via successive shortest augmentations), because adding a
// longer route only pays off when enough trains share the load.
func planMovements(net *Network, start, end, numTrains int) (paths [][]int, trainsEachPath []int, ok bool) {
	numStations := len(net.stationNames)
	trackMap := newTrackMap(2 * numStations) //each station is split into an "in" and "out" node, with a track of capacity 1 between them. So it needs double the number of stations to represent the split network.

	in := func(stationIdx int) int { return 2 * stationIdx }    //entry node
	out := func(stationIdx int) int { return 2*stationIdx + 1 } //exit node

	//add a track from in to out for each station, with capacity 1 (except for start and end stations, which have infinite capacity). This ensures that each intermediate station can only be used by one train at a time.
	for stationIdx := range numStations {
		capacity := 1
		if stationIdx == start || stationIdx == end {
			capacity = bigCap
		}

		//mimic 1 train at 1 stattion at a time
		trackMap.AddTrack(in(stationIdx), out(stationIdx), capacity, 0) //...→ in(victoria) ══[cap 1]══> out(victoria) →...
	}
	directedTracks := make([]directedTrack, 0, len(net.connections)*2) //each undirected connection is represented by two directed tracks, one in each direction
	for _, connection := range net.connections {
		from, to := connection[0], connection[1]
		directedTracks = append(directedTracks,
			directedTrack{from, to, trackMap.AddTrack(out(from), in(to), 1, 1)}, //foward track
			directedTrack{to, from, trackMap.AddTrack(out(to), in(from), 1, 1)}) //reverse track
	}

	maxRoutes := min(numTrains, len(net.adj[start]), len(net.adj[end]))
	bestTurns := -1
	var bestPaths [][]int
	for range maxRoutes {
		cost, found := trackMap.FindPath(out(start), in(end))
		if !found { //"no path exists between start and end`err"
			break
		}
		// Augmentation costs never decrease, so once a new route is at
		// least as long as the best turn count it can never carry a train.
		if bestTurns >= 0 && cost >= bestTurns {
			break
		}
		ps := decompose(trackMap, directedTracks, numStations, start, end)
		lengths := make([]int, len(ps))
		for i, p := range ps {
			lengths[i] = len(p) - 1
		}
		turns := minTurns(lengths, numTrains)
		if bestTurns < 0 || turns < bestTurns {
			bestTurns = turns
			bestPaths = ps
		}
	}
	if bestTurns < 0 {
		return nil, nil, false
	}
	return bestPaths, assignTrains(bestPaths, numTrains, bestTurns), true
}

// decompose extracts the routes carried by the current flow, shortest first.
func decompose(g *trackMap, connEdges []directedTrack, n, start, end int) [][]int {
	used := make(map[[2]int]bool)
	for _, ce := range connEdges {
		if g.tracks[ce.idx].capacity == 0 {
			used[[2]int{ce.from, ce.to}] = true
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
		if used[[2]int{ce.from, ce.to}] {
			next[ce.from] = append(next[ce.from], ce.to)
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
func minTurns(lengths []int, numTrains int) int {
	lo, hi := lengths[0], lengths[0]+numTrains-1
	for lo < hi {
		mid := (lo + hi) / 2
		if capacityWithin(lengths, mid) >= numTrains {
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
func assignTrains(paths [][]int, numTrains, turns int) []int {
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
	for i := len(counts) - 1; i >= 0 && total > numTrains; i-- {
		d := min(counts[i], total-numTrains)
		counts[i] -= d
		total -= d
	}
	return counts
}
