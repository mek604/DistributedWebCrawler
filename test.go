package main

import (
	WG "csc462/DistributedWebCrawler/webgraph"
)
var g WG.Webgraph

func fillGraph() {
	nA := WG.Node{Url:"A"}
	nB := WG.Node{Url:"B"}
	nC := WG.Node{Url:"C"}
	nD := WG.Node{Url:"D"}
	nE := WG.Node{Url:"E"}
	nF := WG.Node{Url:"F"}
	g.AddNode(&nA)
	g.AddNode(&nB)
	g.AddNode(&nC)
	g.AddNode(&nD)
	g.AddNode(&nE)
	g.AddNode(&nF)

	g.AddEdge(&nA, &nB)
	g.AddEdge(&nA, &nC)
	g.AddEdge(&nB, &nE)
	g.AddEdge(&nC, &nE)
	g.AddEdge(&nE, &nF)
	g.AddEdge(&nD, &nA)
}


func main() {
	fillGraph()
	g.String()
}
