package main

// Min-cost flow on the node-split network. Every intermediate station has
// capacity 1, so each unit of flow is one vertex-disjoint route and each
// augmentation yields the cheapest set of k routes by total length.

type flowEdge struct {
	to, cap, cost int
}

type flowGraph struct {
	n     int
	edges []flowEdge
	adj   [][]int
}

func newFlowGraph(n int) *flowGraph {
	return &flowGraph{n: n, adj: make([][]int, n)}
}

func (g *flowGraph) addEdge(u, v, capacity, cost int) int {
	idx := len(g.edges)
	g.adj[u] = append(g.adj[u], idx)
	g.edges = append(g.edges, flowEdge{v, capacity, cost})
	g.adj[v] = append(g.adj[v], idx+1)
	g.edges = append(g.edges, flowEdge{u, 0, -cost})
	return idx
}

const infCost = int(1) << 60

// augmentOne pushes one unit of flow along the cheapest residual path from s
// to t and returns that path's cost. SPFA is used because residual edges of
// earlier augmentations carry negative costs.
func (g *flowGraph) augmentOne(s, t int) (int, bool) {
	dist := make([]int, g.n)
	parent := make([]int, g.n)
	inQueue := make([]bool, g.n)
	for i := range dist {
		dist[i] = infCost
		parent[i] = -1
	}
	dist[s] = 0
	queue := []int{s}
	inQueue[s] = true
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		inQueue[u] = false
		for _, ei := range g.adj[u] {
			e := g.edges[ei]
			if e.cap <= 0 || dist[u]+e.cost >= dist[e.to] {
				continue
			}
			dist[e.to] = dist[u] + e.cost
			parent[e.to] = ei
			if !inQueue[e.to] {
				queue = append(queue, e.to)
				inQueue[e.to] = true
			}
		}
	}
	if parent[t] == -1 {
		return 0, false
	}
	for v := t; v != s; {
		ei := parent[v]
		g.edges[ei].cap--
		g.edges[ei^1].cap++
		v = g.edges[ei^1].to
	}
	return dist[t], true
}
