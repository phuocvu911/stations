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
func planMovements(net *Network, start, end, numTrains int) (bestPaths [][]int, trainsEachPath []int, ok bool) {
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

		//mimic 1 train at 1 station at a time
		trackMap.AddTrack(in(stationIdx), out(stationIdx), capacity, 0) //...→ in(victoria) ══[cap 1]══> out(victoria) →...
	}

	directedTracks := make([]directedTrack, 0, len(net.connections)*2) //each undirected connection is represented by two directed tracks, one in each direction
	for _, connection := range net.connections {
		from, to := connection[0], connection[1]
		directedTracks = append(directedTracks,
			directedTrack{from, to, trackMap.AddTrack(out(from), in(to), 1, 1)}, //forward track
			directedTrack{to, from, trackMap.AddTrack(out(to), in(from), 1, 1)}) //reverse track
	}

	maxRoutes := min(numTrains, len(net.adj[start]), len(net.adj[end]))
	bestTurns := -1 //thats how many lines that we prints in terminal, the min the better

	//main loop: it try using 1 route, then 2, then 3… and keeps whichever gives the fewest turns
	//i.e decide how many routes to use, and which routes to use, to minimize the number of turns needed to move all trains from start to end
	for range maxRoutes {
		cost, found := trackMap.FindPath(out(start), in(end))
		if !found { //"no path exists between start and end err"
			break
		}

		//after we have first path, if second best path is longer bestTurns, no need to consider it
		if bestTurns >= 0 && cost >= bestTurns {
			break
		}

		paths := decompose(trackMap, directedTracks, numStations, start, end)

		//number of hops each path, for example path with 4 stations has 3 hops
		//lengths is also sorted asc
		lengths := make([]int, len(paths))
		for i, path := range paths {
			lengths[i] = len(path) - 1
		}

		turns := minTurns(lengths, numTrains)
		if bestTurns < 0 || turns < bestTurns {
			bestTurns = turns
			bestPaths = paths
		}
	}
	if bestTurns < 0 {
		return nil, nil, false
	}
	return bestPaths, assignTrains(bestPaths, numTrains, bestTurns), true
}

// decompose extracts whatever is currently marked "used" in trackMap.
func decompose(m *trackMap, directedTracks []directedTrack, numStations, start, end int) (paths [][]int) {
	used := make(map[[2]int]bool)
	for _, directedTrack := range directedTracks {
		if m.tracks[directedTrack.idx].capacity == 0 {
			used[[2]int{directedTrack.from, directedTrack.to}] = true
		}
	}

	// Opposite directions on the same track cancel out, drop it
	for key := range used {
		rev := [2]int{key[1], key[0]}
		if used[rev] && used[key] {
			delete(used, key)
			delete(used, rev)
		}
	}

	//build a map of next[station] → []reachable stations
	next := make([][]int, numStations)
	for _, directedTrack := range directedTracks {
		if used[[2]int{directedTrack.from, directedTrack.to}] {
			next[directedTrack.from] = append(next[directedTrack.from], directedTrack.to)
		}
	}

	//extract "used" paths, these are the paths that we will use in final solution
	for len(next[start]) > 0 {
		cur := start
		path := []int{start}

		//follow the first available track until we reach the end or run out of options
		for cur != end && len(next[cur]) > 0 && len(path) <= numStations {
			nxt := next[cur][0]       // this is the next station to visit
			next[cur] = next[cur][1:] //remove it
			path = append(path, nxt)  //add it to the path
			cur = nxt                 //move to the next station
		}
		if cur == end {
			paths = append(paths, path)
		}
	}

	//sort the paths by length, shortest first
	sort.Slice(paths, func(i, j int) bool { return len(paths[i]) < len(paths[j]) })
	return paths
}

// minTurns is the fewest turns needed to move the given number of trains
// over routes of the given lengths: a route of length L pipelines one train
// per turn, so within T turns it delivers T-L+1 trains. -> T= N+L-1
func minTurns(lengths []int, numTrains int) (turns int) {
	shortestLength := lengths[0]

	//binary search for the minimum number of turns needed to deliver numTrains, x is ranged [0,numTrains)
	i := sort.Search(numTrains, func(x int) bool {
		return totalTrains(lengths, shortestLength+x) >= numTrains
	})
	turns = shortestLength + i
	return turns
}

// totalTrains returns the total number of trains that can be delivered within the given turns
func totalTrains(lengths []int, turns int) int {
	total := 0
	for _, length := range lengths {
		if turns >= length {
			total += turns - length + 1
		}
	}
	return total
}

// assignTrains decides how many trains sent through each route so that every
// train still arrives within the given number of turns; surplus capacity is
// trimmed from the longest routes first.
func assignTrains(paths [][]int, numTrains, turns int) []int {
	trainsEachPath := make([]int, len(paths))

	//calculate how many trains with these setups can move, some times it is bigger than numTrains
	total := 0
	for i, p := range paths {
		trainsCount := max(turns-(len(p)-1)+1, 0)
		trainsEachPath[i] = trainsCount
		total += trainsCount
	}

	//trim the longest routes first until we have exactly numTrains trains
	for i := len(trainsEachPath) - 1; i >= 0 && total > numTrains; i-- {
		d := min(trainsEachPath[i], total-numTrains)
		trainsEachPath[i] -= d
		total -= d
	}
	return trainsEachPath
}
