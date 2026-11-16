package graph

type Edge struct {
	u, v   int
	active bool
}

type Graph struct {
	n        int
	edges    []Edge
	adjEdges [][]int
}

func New(n int) *Graph {
	return &Graph{
		n:        n,
		edges:    []Edge{},
		adjEdges: make([][]int, n),
	}
}

func (g *Graph) AddEdge(u, v int) {
	id := len(g.edges)
	g.edges = append(g.edges, Edge{u, v, true})

	g.adjEdges[u] = append(g.adjEdges[u], id)
	g.adjEdges[v] = append(g.adjEdges[v], id)
}

func (g *Graph) Degree(v int) int {
	count := 0

	for _, eID := range g.adjEdges[v] {
		if g.edges[eID].active {
			count++
		}
	}

	return count
}

func (g *Graph) hasActiveEdges() bool {
	for _, e := range g.edges {
		if e.active {
			return true
		}
	}

	return false
}

func (g *Graph) isBridge(eID int) bool {
	e := g.edges[eID]
	if !e.active {
		return false
	}

	g.edges[eID].active = false
	defer func() { g.edges[eID].active = true }()
	return !g.reachable(e.u, e.v)
}

func (g *Graph) reachable(from, to int) bool {
	visited := make([]bool, g.n)
	stack := []int{from}

	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[v] {
			continue
		}

		visited[v] = true

		if v == to {
			return true
		}

		for _, eID := range g.adjEdges[v] {
			if !g.edges[eID].active {
				continue
			}
			e := g.edges[eID]
			next := e.u
			if next == v {
				next = e.v
			}
			if !visited[next] {
				stack = append(stack, next)
			}
		}
	}

	return false
}

func (g *Graph) IsEulerian() bool {
	for i := 0; i < g.n; i++ {
		if g.Degree(i)%2 != 0 {
			return false
		}
	}

	start := -1

	for i := 0; i < g.n; i++ {
		if g.Degree(i) > 0 {
			start = i
			break
		}
	}

	if start == -1 {
		return false
	}

	visited := make([]bool, g.n)
	stack := []int{start}

	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[v] {
			continue
		}

		visited[v] = true

		for _, eID := range g.adjEdges[v] {
			if !g.edges[eID].active {
				continue
			}

			e := g.edges[eID]

			next := e.u
			if next == v {
				next = e.v
			}

			if !visited[next] {
				stack = append(stack, next)
			}
		}
	}

	for i := 0; i < g.n; i++ {
		if g.Degree(i) > 0 && !visited[i] {
			return false
		}
	}

	return true
}

func (g *Graph) Fleury(start int) []int {
	cycle := []int{start}
	current := start

	for g.hasActiveEdges() {
		chosen := -1
		for _, eID := range g.adjEdges[current] {
			if !g.edges[eID].active {
				continue
			}
			if !g.isBridge(eID) {
				chosen = eID
				break
			}
		}

		if chosen == -1 {
			for _, eID := range g.adjEdges[current] {
				if g.edges[eID].active {
					chosen = eID
					break
				}
			}
		}

		if chosen == -1 {
			break
		}

		e := g.edges[chosen]
		next := e.u
		if next == current {
			next = e.v
		}

		g.edges[chosen].active = false
		cycle = append(cycle, next)
		current = next
	}

	return cycle
}

func (g *Graph) GetEdges() []Edge {
	return g.edges
}

func (g *Graph) GetN() int {
	return g.n
}

func (e *Edge) IsActive() bool {
	return e.active
}
