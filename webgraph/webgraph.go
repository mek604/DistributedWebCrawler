package webgraph

// acknowledgement
// https://flaviocopes.com/golang-data-structure-graph/

import (
	"fmt"
	"sync"
)

type Webgraph struct {
	mu sync.Mutex
	nodes []*Node
	edges map[Node][]*Node
}

// web page info
type Node struct {
	Url string
	// content ?
}

func (g *Webgraph) AddNode(n *Node) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.nodes = append(g.nodes, n)
}

func (g *Webgraph) AddEdge(n1, n2 *Node) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.edges == nil {
	    g.edges = make(map[Node][]*Node)
	}
	g.edges[*n1] = append(g.edges[*n1], n2)
	g.edges[*n2] = append(g.edges[*n2], n1)
}

// to remove
func (g *Webgraph) String() {
    g.mu.Lock()
    g.mu.Unlock()
    s := ""
    for i := 0; i < len(g.nodes); i++ {
        s += g.nodes[i].Url + " -> "
        near := g.edges[*g.nodes[i]]
        for j := 0; j < len(near); j++ {
            s += near[j].Url + " "
        }
        s += "\n"
    }
    fmt.Println(s)
}
