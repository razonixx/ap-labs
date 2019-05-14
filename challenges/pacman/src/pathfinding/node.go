package pathfinding

import (
	"fmt"

	"github.com/faiface/pixel"
)

type Node struct {
	ID        int
	Neighbors []*Node
	History   []*Node
	Vec       pixel.Vec
	IsEaten   bool
}

func (n Node) compare(m Node) bool {
	if n.ID == m.ID {
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
	fmt.Printf("Current ID: %d\n", n.ID)
	fmt.Print("Neighbors: \n")
	for i := range n.Neighbors {
		fmt.Print("ID: ")
		fmt.Println(n.Neighbors[i].ID)
	}

}

func Breadthwise(start, end Node) []*Node {
	start.History = start.History[:0]
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
			result = current.History
			result = append(result, current)
			return result
		}
		//DIDnt find node
		for i := 0; i < len(current.Neighbors); i++ {
			currentNeighbor := current.Neighbors[i]
			if !contains(visited, currentNeighbor) {
				currentNeighbor.History = make([]*Node, len(current.History), size)
				copy(currentNeighbor.History, current.History)
				currentNeighbor.History = append(currentNeighbor.History, current)
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

	nodes[0].Neighbors = append(nodes[0].Neighbors, nodes[1], nodes[2])

	nodes[1].Neighbors = append(nodes[1].Neighbors, nodes[0], nodes[3])

	nodes[2].Neighbors = append(nodes[2].Neighbors, nodes[0], nodes[5])

	nodes[3].Neighbors = append(nodes[3].Neighbors, nodes[1], nodes[4])

	nodes[4].Neighbors = append(nodes[4].Neighbors, nodes[3])

	nodes[5].Neighbors = append(nodes[5].Neighbors, nodes[2], nodes[4])

	//delete(nodes, 2)

	res := Breadthwise(*nodes[2], *nodes[4])
	for _, n := range res {
		fmt.Printf("-> %d ", n.ID)
	}
	fmt.Println()
}*/
