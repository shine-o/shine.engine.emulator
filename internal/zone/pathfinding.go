package zone

import (
	"bytes"
	"fmt"
	"math"
	"sort"

	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
)

const (
	PresetNodesMargin = 8
	WallsMargin       = 2
	neighborNodes     = 8
	gridWeight        = 2
)

type ValidCoordinates map[uint16]map[uint16]uint8

type Grid map[int]map[int]*Node

type Node struct {
	Parent     *Node
	X, Y       int
	H, G, F    int
	Opened     bool
	Closed     bool
	Neighbours Nodes
}

type Nodes []*Node

func (e Nodes) Len() int {
	return len(e)
}

func (e Nodes) Less(i, j int) bool {
	return e[i].F < e[j].F
}

func (e Nodes) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (g Grid) GetNearest(x, y int) *Node {
	// for each rectangle given the epicenter X,Y
	// send a go routine to iterate over the rectangle
	// if neighbour is found, send to channel and exit function
	// close channel so all other routines are canceled
	var (
		perimeter = 250
		ch        = make(chan *Node)
	)

	go searchOneTopLeft(g, ch, x, y, perimeter)
	go searchOneTopRight(g, ch, x, y, perimeter)
	go searchOneBottomLeft(g, ch, x, y, perimeter)
	go searchOneBottomRight(g, ch, x, y, perimeter)

	return <-ch
}

func (g Grid) Get(x, y int) *Node {
	n, ok := g[x][y]
	if ok {
		return n
	}
	return nil
}

func CanWalk(wv Grid, x, y int) bool {
	_, ok := wv[x][y]
	return ok
}

func CanWalk2(wv ValidCoordinates, x, y int) bool {
	_, ok := wv[uint16(x)][uint16(y)]
	return ok
}

func GetValidCoordinates(s *data.SHBD) ValidCoordinates {
	vc := make(ValidCoordinates)
	r := bytes.NewReader(s.Data)

	for y := 0; y < s.Y; y++ {
		for x := 0; x < s.X; x++ {
			b, err := r.ReadByte()
			if err != nil {
				return vc
			}

			for i := 0; i < 8; i++ {
				if b&byte(math.Pow(2, float64(i))) == 0 {
					rX := uint16(x*8 + i)
					rY := uint16(y)
					_, ok := vc[rX]
					if !ok {
						vc[rX] = make(map[uint16]uint8)
					}
					vc[rX][rY] = 1
				}
			}
		}
	}
	return vc
}

func RawGrid(s *data.SHBD) Grid {
	wv := make(Grid)
	r := bytes.NewReader(s.Data)
	count := 0

	for y := 0; y < s.Y; y++ {
		for x := 0; x < s.X; x++ {
			b, err := r.ReadByte()
			if err != nil {
				return wv
			}

			for i := 0; i < 8; i++ {
				if b&byte(math.Pow(2, float64(i))) == 0 {
					rX := x*8 + i
					rY := y
					_, ok := wv[rX]
					if !ok {
						wv[rX] = make(map[int]*Node)
					}
					n := &Node{
						X: rX,
						Y: rY,
					}
					wv[rX][rY] = n
					count++
				}
			}
		}
	}
	fmt.Printf("count Nodes %v \n", count)
	return wv
}

func PresetNodesGrid(s *data.SHBD, margin int) Grid {
	pn := make(Grid)
	r := bytes.NewReader(s.Data)
	count := 0
	for y := 0; y < s.Y; y++ {
		for x := 0; x < s.X; x++ {
			b, err := r.ReadByte()
			if err != nil {
				return pn
			}

			for i := 0; i < 8; i++ {
				if b&byte(math.Pow(2, float64(i))) == 0 {
					rX := x*8 + i
					rY := y
					// margin between Nodes
					if rX%margin == 0 && rY%margin == 0 {
						_, ok := pn[rX]
						if !ok {
							pn[rX] = make(map[int]*Node)
						}

						pn[rX][rY] = &Node{
							X: rX,
							Y: rY,
						}
						count++
					}
				}
			}
		}
	}
	fmt.Printf("count Nodes %v \n", count)
	return pn
}

// preset Nodes + extra walls so that NPCs don't Get too close to collision areas
func PresetNodesWithMargins(s *data.SHBD, wv Grid, margin int) Grid {
	var (
		pn     = make(Grid)
		r      = bytes.NewReader(s.Data)
		ncount = 0
	)

	for y := 0; y < s.Y; y++ {
		for x := 0; x < s.X; x++ {
			b, err := r.ReadByte()
			if err != nil {
				return pn
			}

			for i := 0; i < 8; i++ {
				if b&byte(math.Pow(2, float64(i))) == 0 {
					rX := x*8 + i
					rY := y
					// add only half the Nodes needed
					count := len(adjacentNodes(wv, wv, rX, rY, WallsMargin))
					if count >= neighborNodes { // do not add Nodes if they are too close to inaccessible Nodes
						if rX%margin == 0 && rY%margin == 0 {

							_, ok := pn[rX]
							if !ok {
								pn[rX] = make(map[int]*Node)
							}

							pn[rX][rY] = &Node{
								X: rX,
								Y: rY,
							}
							ncount++
						}
					}
				}
			}
		}
	}
	fmt.Printf("count Nodes %v \n", ncount)
	return pn
}

func searchOneTopLeft(g Grid, ch chan<- *Node, x, y, perimeter int) {
	bx := x - perimeter
	by := y - perimeter
	for cy := y; cy != by; cy-- {
		for cx := x; cx != bx; cx-- {
			n := g.Get(cx, cy)
			if n != nil {
				select {
				case ch <- n:
				default:
					return
				}
			}
		}
	}
	ch <- nil
}

func searchOneTopRight(g Grid, ch chan<- *Node, x, y, perimeter int) {
	bx := x + perimeter
	by := y - perimeter
	for cy := y; cy != by; cy-- {
		for cx := x; cx != bx; cx++ {
			n := g.Get(cx, cy)
			if n != nil {
				select {
				case ch <- n:
				default:
					return
				}
			}
		}
	}
	ch <- nil
}

func searchOneBottomLeft(g Grid, ch chan<- *Node, x, y, perimeter int) {
	bx := x - perimeter
	by := y + perimeter
	for cy := y; cy != by; cy++ {
		for cx := x; cx != bx; cx-- {
			n := g.Get(cx, cy)
			if n != nil {
				ch <- n
				select {
				case ch <- n:
				default:
					return
				}
			}
		}
	}
	ch <- nil
}

func searchOneBottomRight(g Grid, ch chan<- *Node, x, y, perimeter int) {
	bx := x + perimeter
	by := y + perimeter
	for cy := y; cy != by; cy++ {
		for cx := x; cx != bx; cx++ {
			n := g.Get(cx, cy)
			if n != nil {
				select {
				case ch <- n:
				default:
					return
				}
			}
		}
	}
	ch <- nil
}

func getAllTopLeft(g Grid, ch chan<- *Node, done chan<- bool, x, y, perimeter int) {
	bx := x - perimeter
	by := y - perimeter
	for cy := y; cy != by; cy-- {
		for cx := x; cx != bx; cx-- {
			n := g.Get(cx, cy)
			if n != nil {
				ch <- n
			}
		}
	}
	done <- true
}

func getAllTopRight(g Grid, ch chan<- *Node, done chan<- bool, x, y, perimeter int) {
	bx := x + perimeter
	by := y - perimeter
	for cy := y; cy != by; cy-- {
		for cx := x; cx != bx; cx++ {
			n := g.Get(cx, cy)
			if n != nil {
				ch <- n
			}
		}
	}
	done <- true
}

func getAllBottomLeft(g Grid, ch chan<- *Node, done chan<- bool, x, y, perimeter int) {
	bx := x - perimeter
	by := y + perimeter
	for cy := y; cy != by; cy++ {
		for cx := x; cx != bx; cx-- {
			n := g.Get(cx, cy)
			if n != nil {
				ch <- n
			}
		}
	}
	done <- true
}

func getAllBottomRight(g Grid, ch chan<- *Node, done chan<- bool, x, y, perimeter int) {
	bx := x + perimeter
	by := y + perimeter
	for cy := y; cy != by; cy++ {
		for cx := x; cx != bx; cx++ {
			n := g.Get(cx, cy)
			if n != nil {
				ch <- n
			}
		}
	}
	done <- true
}

func allNearbyNodes(g Grid, x, y int, perimeter int) Nodes {
	var (
		list Nodes
		ch   = make(chan *Node, 1500)
		tld  = make(chan bool)
		trd  = make(chan bool)
		bld  = make(chan bool)
		brd  = make(chan bool)
		// sem = make(chan int, 100)
		// wg sync.WaitGroup
		topLeft, topRight, bottomLeft, bottomRight bool
	)

	go getAllTopLeft(g, ch, tld, x, y, perimeter)
	go getAllTopRight(g, ch, trd, x, y, perimeter)
	go getAllBottomLeft(g, ch, bld, x, y, perimeter)
	go getAllBottomRight(g, ch, brd, x, y, perimeter)

	for {
		select {
		case n := <-ch:
			if !in(list, n) {
				list = append(list, n)
			}
			//wg.Add(1)
			//sem <- 1
			//go func(open *nodeList, n * Node) {
			//	defer wg.Done()
			//	n.G = euclideanDistance(n.X, n.Y, b.X, b.Y)
			//	n.H = euclideanDistance(n.X, n.Y, a.X, a.Y)
			//	n.F = n.G+n.H
			//	open.RLock()
			//	open.list = append(open.list, n)
			//	open.RUnlock()
			//	<-sem
			//}(open, n)

		case <-tld:
			topLeft = true
		case <-trd:
			topRight = true
		case <-bld:
			bottomLeft = true
		case <-brd:
			bottomRight = true
		default:
			if topLeft && topRight && bottomLeft && bottomRight {
				return list
			}
		}
	}
}

// ng is a Grid with preset Nodes that have a given margin
// rg is a raw Grid, used to avoid obstacles between Nodes with margin
func astar(ng, rg Grid, fx, fy, tx, ty, margin int, heuristic string) Nodes {
	var (
		open, foundPath Nodes
		cnode           *Node
		a               = ng.Get(fx, fy)
		b               = ng.Get(tx, ty)
	)

	if a == nil {
		a = ng.GetNearest(fx, fy)
	}

	if b == nil {
		b = ng.GetNearest(tx, ty)
	}

	if a == nil || b == nil {
		logger.Fatalf("a or b is nil %v %v", a, b)
		return nil
	}

	a.G = 0
	a.F = 0

	open = append(open, a)

	a.Opened = true

	for {
		open, cnode = lowestF(open)

		if equal(cnode, b) {
			break
		}

		for _, neighbor := range adjacentNodes(ng, rg, cnode.X, cnode.Y, margin) {
			var ng int
			if neighbor.Closed {
				continue
			}

			switch heuristic {
			case "manhatan":
				ng = cnode.G + manhatanDistance(cnode.X, cnode.Y, neighbor.X, neighbor.Y)
				break
			case "euclidean":
				ng = cnode.G + euclideanDistance(cnode.X, cnode.Y, neighbor.X, neighbor.Y)
				break
			case "octile":
				ng = cnode.G + octileDistance(cnode.X, cnode.Y, neighbor.X, neighbor.Y)
				break
			default:
				ng = cnode.G + octileDistance(cnode.X, cnode.Y, neighbor.X, neighbor.Y)
			}

			if !neighbor.Opened || ng < neighbor.G {
				neighbor.G = ng
				switch heuristic {
				case "manhatan":
					neighbor.H = gridWeight * manhatanDistance(neighbor.X, neighbor.Y, b.X, b.Y)
					break
				case "octile":
					neighbor.H = gridWeight * octileDistance(neighbor.X, neighbor.Y, b.X, b.Y)
					break
				case "euclidean":
					neighbor.H = gridWeight * euclideanDistance(neighbor.X, neighbor.Y, b.X, b.Y)
					break
				default:
					neighbor.H = gridWeight * euclideanDistance(neighbor.X, neighbor.Y, b.X, b.Y)
					break
				}

				neighbor.F = neighbor.G + neighbor.H
				neighbor.Parent = cnode

				if !neighbor.Opened {
					neighbor.Opened = true
					open = append(open, neighbor)
				}
			}
		}
	}

	cnode.Closed = true

	next := cnode
	for next != nil {
		foundPath = append(foundPath, next)
		next = next.Parent
	}

	sort.Sort(sort.Reverse(foundPath))

	return foundPath
}

func manhatanDistance(fx, fy, tx, ty int) int {
	dx := math.Abs(float64(tx - fx))
	dy := math.Abs(float64(ty - fy))
	return int(dx + dy)
}

// sqrt(dx * dx + dy * dy)
func euclideanDistance(fx, fy, tx, ty int) int {
	dx := tx - fx
	dy := ty - fy
	d := math.Sqrt(float64(dx*dx + dy*dy))
	return int(d)
}

func octileDistance(fx, fy, tx, ty int) int {
	F := math.Sqrt2 - 1
	dx := math.Abs(float64(tx - fx))
	dy := math.Abs(float64(ty - fy))

	if dx < dy {
		return int(F*dx + dy)
	}

	return int(F*dy + dx)
}

// given speed, create a new Node slice with speed as euclideanDistance value between Nodes
// a => n+1 => b
// for each n
// 	if euclideanDistance between a, n+1 == speed
//	  	add a to new slice
//		a = n+1
func reduce(e Nodes, speed int) Nodes {
	var (
		rp      Nodes
		current *Node
	)

	current = e[0]
	for i := 0; i < len(e); i++ {
		d := euclideanDistance(current.X, current.Y, e[i].X, e[i].Y)
		if d >= speed {
			rp = append(rp, current)
			current = e[i]
		}
	}
	return rp
}

func equal(v1, v2 *Node) bool {
	if v1.X == v2.X && v1.Y == v2.Y {
		return true
	}
	return false
}

func lowestF(open Nodes) ([]*Node, *Node) {
	var (
		uo = make([]*Node, 0)
		v  *Node
	)

	sort.Sort(open)
	v = open[0]

	uo = append(uo, open[1:]...)

	return uo, v
}

func topObstacles(raw Grid, x, y, ty int) bool {
	for y != ty {
		if !CanWalk(raw, x, y) {
			return true
		}
		y--
	}

	return false
}

func bottomObstacles(wv Grid, x, y, ty int) bool {
	for y != ty {
		if !CanWalk(wv, x, y) {
			return true
		}
		y++
	}
	return false
}

func rightObstacles(wv Grid, x, y, tx int) bool {
	for x != tx {
		if !CanWalk(wv, x, y) {
			return true
		}
		x++
	}
	return false
}

func leftObstacles(wv Grid, x, y, tx int) bool {
	for x != tx {
		if !CanWalk(wv, x, y) {
			return true
		}
		x--
	}
	return false
}

func topRightDiagonalObstacles(wv Grid, x, y, tx, ty int) bool {
	for x != tx && y != ty {
		if !CanWalk(wv, x, y) {
			return true
		}
		x++
		y--
	}

	return false
}

func topLeftDiagonalObstacles(wv Grid, x, y, tx, ty int) bool {
	for x != tx && y != ty {
		if !CanWalk(wv, x, y) {
			return true
		}
		x--
		y--
	}

	return false
}

func bottomRightDiagonalObstacles(wv Grid, x, y, tx, ty int) bool {
	for x != tx && y != ty {
		if !CanWalk(wv, x, y) {
			return true
		}
		x++
		y++
	}

	return false
}

func bottomLeftDiagonalObstacles(wv Grid, x, y, tx, ty int) bool {
	for x != tx && y != ty {
		if !CanWalk(wv, x, y) {
			return true
		}
		x--
		y++
	}

	return false
}

// using Node Grid and raw Grid, Get neighboring Nodes
func adjacentNodes(ng, rg Grid, x, y, margin int) Nodes {
	var (
		result Nodes
		n      *Node
	)

	// ↑
	if !topObstacles(rg, x, y, y-margin) {
		n = ng.Get(x, y-margin)
		if n != nil {
			result = append(result, n)
		}
	}

	// ↓
	if !bottomObstacles(rg, x, y, y+margin) {
		n = ng.Get(x, y+margin)
		if n != nil {
			result = append(result, n)
		}
	}

	// →
	if !rightObstacles(rg, x, y, x+margin) {
		n = ng.Get(x+margin, y)
		if n != nil {
			result = append(result, n)
		}
	}

	// ←
	if !leftObstacles(rg, x, y, x-margin) {
		n = ng.Get(x-margin, y)
		if n != nil {
			result = append(result, n)
		}
	}

	// ↖
	if !topLeftDiagonalObstacles(rg, x, y, x-margin, y-margin) {
		n = ng.Get(x-margin, y-margin)
		if n != nil {
			result = append(result, n)
		}
	}

	// ↗
	if !topRightDiagonalObstacles(rg, x, y, x+margin, y-margin) {
		n = ng.Get(x+margin, y-margin)
		if n != nil {
			result = append(result, n)
		}
	}

	// ↘
	if !bottomRightDiagonalObstacles(rg, x, y, x+margin, y+margin) {
		n = ng.Get(x+margin, y+margin)
		if n != nil {
			result = append(result, n)
		}
	}

	// ↙
	if !bottomLeftDiagonalObstacles(rg, x, y, x-margin, y+margin) {
		n = ng.Get(x-margin, y+margin)
		if n != nil {
			result = append(result, n)
		}
	}

	//
	//if len(result) == 0 {
	//	fmt.Println("no Neighbours")
	//	n := ng.GetNearest(X, Y)
	//	if n != nil {
	//		result = append(result, n)
	//	} else {
	//		fmt.Println("still no Neighbours")
	//	}
	//}

	return result
}

func in(list Nodes, sv *Node) bool {
	for _, v := range list {
		if equal(v, sv) {
			return true
		}
	}
	return false
}

func solveForTime(s, d int) int {
	return d / s
}

func solveForSpeed(d, t int) int {
	return d / t
}

func copyGrid(g Grid) Grid {
	ng := make(Grid)
	for k1, v1 := range g {
		for k2, v2 := range v1 {
			n := *v2

			_, ok := ng[k1]
			if !ok {
				ng[k1] = make(map[int]*Node)
			}
			ng[k1][k2] = &n
		}
	}
	return ng
}
