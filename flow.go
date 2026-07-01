package main

// Min-cost flow on the node-split network. Every intermediate station has
// capacity 1, so each unit of flow is one vertex-disjoint route and each
// augmentation yields the cheapest set of k routes by total length.

//road from a to b
type track struct {
	to, cap, cost int
}

//map of the railway network
type trackMap struct {
	numOfStations int
	tracks        []track
	adjacentsList [][]int
}

func newTrackMap(n int) *trackMap {
	return &trackMap{numOfStations: n, adjacentsList: make([][]int, n)}
}

// addTrack adds a directed track from "from" to "to" with the
// given capacity and cost, returning the index of the forward edge.
// The reverse edge is added automatically with capacity 0 and negative cost
// for return.
func (f *trackMap) addTrack(from, to, capacity, cost int) int {
	idx := len(f.tracks)

	//forward track
	f.adjacentsList[from] = append(f.adjacentsList[from], idx)
	f.tracks = append(f.tracks, track{to, capacity, cost})

	//phantom track for return, with negative cost and 0 capacity
	f.adjacentsList[to] = append(f.adjacentsList[to], idx+1)
	f.tracks = append(f.tracks, track{from, 0, -cost})
	return idx
}

const infCost = int(1) << 60

//This function dispatches one train from start station to end station along the
//cheapest open route, and marks that route as used. It returns (how far the train
//traveled, did it make it?)
func (f *trackMap) findPath(start, end int) (int, bool) {
	distances := make([]int, f.numOfStations)
	prevs := make([]int, f.numOfStations)
	inQueue := make([]bool, f.numOfStations)

	//start dumb, every stations is unreachable and we came from nowhere at first
	for i := range distances {
		distances[i] = infCost
		prevs[i] = -1
	}
	//distance from start to start is 0
	distances[start] = 0
	queue := []int{start}
	inQueue[start] = true

	//find the shortest path from start to end using a modified Bellman-Ford algorithm
	for len(queue) > 0 {
		station := queue[0]
		queue = queue[1:]
		inQueue[station] = false
		for _, adjacents := range f.adjacentsList[station] {
			e := f.tracks[adjacents]
			if e.cap <= 0 || distances[station]+e.cost >= distances[e.to] {
				continue
			}
			distances[e.to] = distances[station] + e.cost
			prevs[e.to] = adjacents
			if !inQueue[e.to] {
				queue = append(queue, e.to)
				inQueue[e.to] = true
			}
		}
	}

	//log if no path was found
	if prevs[end] == -1 {
		return 0, false
	}

	//update the capacities of the tracks along the path found, use the capacity of the forward track
	// and increase the capacity of the reverse edge, go from end to start, following the previous stations
	for station := end; station != start; {
		from := prevs[station]
		f.tracks[from].cap--
		f.tracks[from^1].cap++
		station = f.tracks[from^1].to
	}
	return distances[end], true
}
