package main

// Min-cost flow on the node-split network. Every intermediate station has
// capacity 1, so each unit of flow is one vertex-disjoint route and each
// augmentation yields the cheapest set of k routes by total length.

//road from a to b
type track struct {
	to       int
	capacity int
	cost     int
}

//map of the railway network
type trackMap struct {
	numStations int
	tracks      []track
	adjacency   [][]int //adjacency list of track indices (in 'tracks' slice) for each station
}

func newTrackMap(n int) *trackMap {
	return &trackMap{numStations: n, adjacency: make([][]int, n)}
}

// addTrack adds a directed track from "from" to "to" with the
// given capacity and cost, returning the index of the forward edge.
// The reverse edge is added automatically with capacity 0 and negative cost
// for return.
func (m *trackMap) addTrack(from, to, capacity, cost int) int {
	idx := len(m.tracks)

	//forward track, always idx 0,2,4,... and reverse track is idx 1,3,5,... in tracks slice
	m.tracks = append(m.tracks, track{to, capacity, cost})
	m.adjacency[from] = append(m.adjacency[from], idx)

	//phantom reserve track for return, with negative cost and 0 capacity.they are always shown in pairs
	m.tracks = append(m.tracks, track{from, 0, -cost})
	m.adjacency[to] = append(m.adjacency[to], idx+1)
	return idx
}

const infCost = int(1) << 60

//This function dispatches one train from start station to end station along the
//cheapest open route, and marks that route as used. It returns (how far the train
//traveled, did it make it?)
func (m *trackMap) findPath(start, end int) (int, bool) {
	distances := make([]int, m.numStations)
	prevEdge := make([]int, m.numStations)
	inQueue := make([]bool, m.numStations)

	//start dumb, every station is unreachable and we came from nowhere at first
	for i := range distances {
		distances[i] = infCost
		prevEdge[i] = -1
	}
	//distance from start to start is 0
	distances[start] = 0
	queue := []int{start}
	inQueue[start] = true

	//find the shortest path from start to end using a modified Bellman-Ford (SPFA) algorithm
	for len(queue) > 0 {
		station := queue[0]
		queue = queue[1:]
		inQueue[station] = false
		trackIdxs := m.adjacency[station]
		for _, trackIdx := range trackIdxs {
			track := m.tracks[trackIdx]
			//if the track is full or the new distance is not better, skip it
			if track.capacity <= 0 || distances[station]+track.cost >= distances[track.to] {
				continue
			}

			//update the distance and previous edge for the station at the end of this track,
			//which is literally means "we use this track to get to end station"
			distances[track.to] = distances[station] + track.cost
			prevEdge[track.to] = trackIdx
			if !inQueue[track.to] {
				queue = append(queue, track.to)
				inQueue[track.to] = true
			}
		}
	}

	//no path was found
	if prevEdge[end] == -1 {
		return 0, false
	}

	//walk the path from end to start, spending one unit of capacity on each
	//forward track and returning it to the paired reverse track
	for station := end; station != start; {
		trackIdx := prevEdge[station]
		m.tracks[trackIdx].capacity--     //mark the forward track as used
		m.tracks[trackIdx^1].capacity++   //XOR with 1 flips the last bit, which is how we paired forward and reverse tracks
		station = m.tracks[trackIdx^1].to //we have no from so "to" of the reverse track is the previous station
	}
	return distances[end], true
}
