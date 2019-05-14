package main

import (
	"fmt"

	"github.com/faiface/pixel"
)

type Node struct {
	id        int
	neighbors []*Node
	history   []*Node
	vec       pixel.Vec
	isEaten   bool
}

func (n Node) compare(m Node) bool {
	if n.id == m.id {
		return true
	}
	return false
}

func contains(s []*Node, e *Node) bool {
	for _, a := range s {
		if a.compare(*e) {
			return true
		}
	}
	return false
}

func printNode(n Node) {
	fmt.Printf("Current ID: %d\n", n.id)
	fmt.Print("Neighbors: \n")
	for i := range n.neighbors {
		fmt.Print("ID: ")
		fmt.Println(n.neighbors[i].id)
	}

}

func Breadthwise(start, end Node) []*Node {
	start.history = start.history[:0]
	size := 10
	result := make([]*Node, 0, size)

	visited := make([]*Node, 0, size)

	work := make([]*Node, 0, size)

	visited = append(visited, &start)       //visited.Add(start)
	work = append([]*Node{&start}, work...) //work.Enqueue(start)
	for len(work) > 0 {
		current := work[len(work)-1]
		work = work[:len(work)-1] //current = work.Dequeue
		if current.compare(end) {
			//Found node
			result = current.history
			result = append(result, current)
			return result
		}
		//Didnt find node
		for i := 0; i < len(current.neighbors); i++ {
			currentNeighbor := current.neighbors[i]
			if !contains(visited, currentNeighbor) {
				currentNeighbor.history = make([]*Node, len(current.history), size)
				copy(currentNeighbor.history, current.history)
				currentNeighbor.history = append(currentNeighbor.history, current)
				visited = append(visited, currentNeighbor)
				work = append([]*Node{currentNeighbor}, work...)
			}
		}
	}
	return nil
}

/*func main() {
	nodes := make(map[int]*Node, 6)
	nodes[0] = &Node{0, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(0, 0)}
	nodes[1] = &Node{1, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(0, 0)}
	nodes[2] = &Node{2, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(0, 0)}
	nodes[3] = &Node{3, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(0, 0)}
	nodes[4] = &Node{4, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(0, 0)}
	nodes[5] = &Node{5, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(0, 0)}

	nodes[0].neighbors = append(nodes[0].neighbors, nodes[1], nodes[2])

	nodes[1].neighbors = append(nodes[1].neighbors, nodes[0], nodes[3])

	nodes[2].neighbors = append(nodes[2].neighbors, nodes[0], nodes[5])

	nodes[3].neighbors = append(nodes[3].neighbors, nodes[1], nodes[4])

	nodes[4].neighbors = append(nodes[4].neighbors, nodes[3])

	nodes[5].neighbors = append(nodes[5].neighbors, nodes[2], nodes[4])

	//delete(nodes, 2)

	res := Breadthwise(*nodes[2], *nodes[4])
	for _, n := range res {
		fmt.Printf("-> %d ", n.id)
	}
	fmt.Println()
}*/
