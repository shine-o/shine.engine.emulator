package main

import (
	"bytes"
	"fmt"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"image"
	"image/color"
	"math"
	"sort"
	"testing"
)

func Test_Path_A_B_astar(t *testing.T) {
	m := "Rou"
	var s *data.SHBD
	s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", m))

	if err != nil {
		logger.Error(err)
	}

	img, err := data.SHBDToImage(s)

	if err != nil {
		logger.Error(err)
	}

	a := &node{
		x: 835,
		y: 700,
	}

	b := &node{
		x: 1070,
		y: 1540,
	}

	v := initNodes(s)

	pn := initPathNodes(s, v)

	if err != nil {
		logger.Fatal(err)
	}

	pathVertices := astar(pn, a, b)

	for _, pv := range pathVertices {
		img.Set(pv.x, pv.y, color.RGBA{
			R: 207,
			G: 0,
			B: 15,
			A: 1,
		})
	}

	err = data.SaveBmpFile(img, "./", "painted path")

	if err != nil {
		logger.Error(err)
	}
	// print an image with the node data
}

func Test_Path_A_B_jps(t *testing.T) {
	m := "Rou"
	var s *data.SHBD
	s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", m))

	if err != nil {
		logger.Error(err)
	}

	img, err := data.SHBDToImage(s)

	if err != nil {
		logger.Error(err)
	}

	//a := &node{
	//	x: 1111,
	//	y: 1577,
	//}
	//
	//b := &node{
	//	x: 1111,
	//	y: 1582,
	//}

	a := &node{
		x: 1441,
		y: 172,
	}

	b := &node{
		x: 500,
		y: 1570,
	}

	v := initNodes(s)
	//
	//pn := initPathNodes(s, v)
	//
	//if err != nil {
	//	logger.Fatal(err)
	//}

	pathVertices := jps(v, a, b)

	for _, pv := range pathVertices {
		img.Set(pv.x, pv.y, color.RGBA{
			R: 207,
			G: 0,
			B: 15,
			A: 1,
		})
	}

	err = data.SaveBmpFile(img, "./", "painted path")

	if err != nil {
		logger.Error(err)
	}
	// print an image with the node data
}

// +10
func Test_Paint_Path_Nodes(*testing.T) {
	m := "Rou"
	var s *data.SHBD
	s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", m))

	if err != nil {
		logger.Error(err)
	}

	v := initNodes(s)

	if err != nil {
		logger.Fatal(err)
	}

	img, err := SHBDToImage(s, v)

	if err != nil {
		logger.Error(err)
	}

	err = data.SaveBmpFile(img, "./", "painted path")

	if err != nil {
		logger.Error(err)
	}
	// print an image with the node data

	// paint nodes that will be used for paths.
}

func Test_Map_Intermitent_By_Speed_Path_A_B_AStar(t *testing.T) {
	// raw path has too many nodes
	// given a speed, reduce the raw path using the speed
	// have a function that calculates how many nodes should be sent per second

	// speed is a constant, the entity may be walking / running, e.g for entities 120, 60

	// s = d / t
	// 120 = d / t

	// 120 = 180 / t

	// t = 180 / 120
	// t = 1.5

	// t will be the time in seconds
	// create a ticker that will send those packets

	// I will also need to calculate the euclideanDistance per second

	m := "Rou"
	var s *data.SHBD
	s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", m))

	if err != nil {
		logger.Error(err)
	}

	img, err := data.SHBDToImage(s)

	if err != nil {
		logger.Error(err)
	}

	a := &node{
		x: 515,
		y: 1177,
	}

	b := &node{
		x: 1250,
		y: 1500,
	}

	v := initNodes(s)

	pn := initPathNodes(s, v)

	if err != nil {
		logger.Fatal(err)
	}

	pathVertices := astar(pn, a, b)

	pathVertices = reduce(pathVertices, 15)

	for _, pv := range pathVertices {
		img.Set(pv.x, pv.y, color.RGBA{
			R: 207,
			G: 0,
			B: 15,
			A: 1,
		})
	}

	err = data.SaveBmpFile(img, "./", "painted path")

	if err != nil {
		logger.Error(err)
	}

}

func Benchmark_ReduceVertices(b *testing.B) {

	m := "Rou"
	var s *data.SHBD
	s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", m))

	if err != nil {
		logger.Error(err)
	}

	pa := &node{
		x: 515,
		y: 1177,
	}

	pb := &node{
		x: 1250,
		y: 1500,
	}

	v := initNodes(s)

	pn := initPathNodes(s, v)

	if err != nil {
		logger.Fatal(err)
	}

	pathVertices := astar(pn, pa, pb)
	b.Run("reduce", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = reduce(pathVertices, 15)
		}
	})
}

func Benchmark_InitVertices(b *testing.B) {
	m := "Rou"
	var s *data.SHBD
	s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", m))

	if err != nil {
		logger.Error(err)
	}

	var v grid
	b.Run("initNodes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v = initNodes(s)
		}
	})

	b.Run("pathNodes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			initPathNodes(s, v)
		}
	})
}

func Benchmark_astar_algorithm(b *testing.B) {
	m := "Rou"
	var s *data.SHBD
	s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", m))

	if err != nil {
		logger.Error(err)
	}

	v := initNodes(s)
	pn := initPathNodes(s, v)

	if err != nil {
		logger.Fatal(err)
	}

	pa := &node{
		x: 835,
		y: 700,
	}

	pb := &node{
		x: 1070,
		y: 1540,
	}

	b.Run("astar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			astar(pn, pa, pb)
		}
	})
}

func Benchmark_jps_algorithm(b *testing.B) {
	// base astar algo
	// find successors
	// foreach neighbour, jump
}

func SHBDToImage(s *data.SHBD, wv grid) (*image.RGBA, error) {
	r := bytes.NewReader(s.Data)

	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: s.X * 8,
			Y: s.Y,
		},
	})

	for y := 0; y < s.Y; y++ {
		for x := 0; x < s.X; x++ {
			b, err := r.ReadByte()
			if err != nil {
				return img, err
			}

			// every ten pixels
			for i := 0; i < 8; i++ {
				var (
					rX, rY int
					c      color.Color
				)

				rX = x*8 + i
				rY = y

				if b&byte(math.Pow(2, float64(i))) == 0 {
					count := len(adjacentNodes(wv, &node{
						x: rX,
						y: rY,
					}, nodesMargin))

					if count == neighborNodes {
						c = color.RGBA{
							R: 255,
							G: 0,
							B: 0,
							A: 1,
						}
					} else {
						c = color.White
					}
				} else {
					c = color.Black
				}

				img.Set(rX, rY, c)
			}
		}
	}
	return img, nil
}

func initPathNodes(s *data.SHBD, wv grid) grid {

	var (
		pn = make(grid)
		r  = bytes.NewReader(s.Data)
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
					// add only half the nodes needed
					count := len(adjacentNodes(wv, &node{
						x: rX,
						y: rY,
					}, nodesMargin))

					if count == neighborNodes { // do not add nodes if they are too close to inaccessible nodes
						_, ok := pn[rX]
						if !ok {
							pn[rX] = make(map[int]*node)
						}

						pn[rX][rY] = &node{
							x: rX,
							y: rY,
						}
					}
				}
			}
		}
	}
	return pn
}

const (
	neighborNodes = 8
	nodesMargin   = 4
)

// sqrt(dx * dx + dy * dy)
func euclideanDistance(a, b *node) int {
	d := math.Sqrt(math.Pow(math.Abs(float64(a.x-b.x)), 2) + math.Pow(math.Abs(float64(a.y-b.y)), 2))
	return int(d)
}

// sqrt(dx * dx + dy * dy)
func octileDistance(a, b * node) int {
	var F = math.Sqrt2 - 1
	dx := math.Abs(float64(a.x - b.x))
	dy := math.Abs(float64(a.y - b.y))

	if dx < dy {
		return int(F * dx + dy)
	}

	return int(F * dy + dx)
}

type grid map[int]map[int]*node

type node struct {
	parent  *node
	x, y    int
	h, g, f int
	opened  bool
	closed  bool
}

type nodes []*node

func (e nodes) Len() int {
	return len(e)
}

func (e nodes) Less(i, j int) bool {
	return e[i].h < e[j].h
}

func (e nodes) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// jump point search
// good for open spaces
func jps(grid grid, a *node, b *node) nodes {
	var (
		open, foundPath nodes
		node            *node
	)

	a.g = 0
	a.f = 0

	open = append(open, a)
	a.opened = true

	for len(open) != 0 {
		open, node = lowestF(open)
		node.closed = true
		if equal(node, b) {
			break
		}
		identifySuccessors(grid, &open, node, b)
	}

	next := node
	for next != nil {
		foundPath = append(foundPath, next)
		next = next.parent
	}

	sort.Sort(sort.Reverse(foundPath))

	return foundPath
}

// this one fails to add nodes to openList
func identifySuccessors(grid grid, openList *nodes, node, b *node) {
	var (
		neighbors = adjacentNodes(grid, node, 1)
		nx, ny    int
		x         = node.x
		y         = node.y
		ng              int

	)

	for _, neighbor := range neighbors {
		if neighbor == nil {
			continue
		}
		nx = neighbor.x
		ny = neighbor.y
		jumpNode := jump(grid, nx, ny, x, y, b)
		if jumpNode != nil {

			if jumpNode.closed {
				continue
			}

			d := octileDistance(jumpNode, node)
			//d := euclideanDistance(jumpNode, node)

			ng = node.g + d

			if !jumpNode.opened || ng < jumpNode.g {
				jumpNode.g = ng
				jumpNode.h = euclideanDistance(jumpNode, b)
				jumpNode.f = jumpNode.g + jumpNode.h
				jumpNode.parent = node
				if !jumpNode.opened {
					*openList = append(*openList, jumpNode)
					jumpNode.opened = true
				}
			}
		}
	}
}

func findNeighbors(grid grid, node *node) nodes {
	var (
		parent         = node.parent
		x              = node.x
		y              = node.y
		px, py, dx, dy int
		neighbors      nodes
	)

	// directed pruning: can ignore most neighbors, unless forced.
	if parent != nil {
		px = parent.x
		py = parent.y
		// get the normalized direction of travel
		dx = int(float64(x-px) / math.Max(math.Abs(float64(x-px)), 1))
		dy = int(float64(y-py) / math.Max(math.Abs(float64(y-py)), 1))

		if dx != 0 && dy != 0 {
			if canWalk(grid, x, y+dy) {
				neighbors = append(neighbors, grid[x][y+dy])
			}

			if canWalk(grid, x+dx, y) {
				neighbors = append(neighbors, grid[x+dx][y])
			}

			if canWalk(grid, x, y+dy) && canWalk(grid, x+dx, y) {
				neighbors = append(neighbors, grid[x+dx][y+dy])
			}
		} else { // search horizontally/vertically
			var isNextWalkable bool
			if dx != 0 {
				isNextWalkable = canWalk(grid, x+dx, y)
				var (
					isTopWalkable = canWalk(grid, x, y+1)
					isBottomWalkable = canWalk(grid,x, y-1)
				)

				if isNextWalkable {
					neighbors = append(neighbors, grid[x+dx][y])
					if isTopWalkable {
						neighbors = append(neighbors, grid[x+dx][y+1])
					}
					if isBottomWalkable {
						neighbors = append(neighbors, grid[x+dx][y-1])
					}
				}

			} else if dy != 0 {
				isNextWalkable = canWalk(grid, x, y+dy)

				var (
					isRightWalkable = canWalk(grid, x+1, y)
					isLeftWalkable = canWalk(grid, x-1, y)
				)

				if isNextWalkable {
					neighbors = append(neighbors, grid[x][y+dy])
					if isRightWalkable {
						neighbors = append(neighbors, grid[x+1][y+dy])
					}
					if isLeftWalkable {
						neighbors = append(neighbors, grid[x-1][y+dy])
					}
				}
			}
		}
	} else { // return all neighbors
		for _, neighborNode := range adjacentNodes(grid, node, 1) {
			neighbors = append(neighbors, neighborNode)
		}
	}
	return neighbors
}

func jump(grid grid, x, y, px, py int, b *node) *node {
	var (
		dx = x - px
		dy = y - py
	)

	if !canWalk(grid, x, y) {
		return nil
	}

	node, ok := grid[x][y]
	if ok && equal(node, b) {
		return node
	}

	if dx != 0 && dy != 0 {
		if (canWalk(grid, x-dx, y+dy) && !canWalk(grid, x-dx, y)) ||
			canWalk(grid, x+dx, y-dy) && !canWalk(grid, x, y-dy) {
			return grid[x][y]
		}
		// when moving diagonally, must check for vertical/horizontal jump points
		if  jump(grid, x+dx, y, x, y, b) != nil ||
			jump(grid, x, y+dy, x, y, b) != nil {
			return grid[x][y]
		}
	} else {    // horizontally/vertically
		if dx != 0 { // moving along x
			if (canWalk(grid, x, y-1) && !canWalk(grid, x-dx, y-1)) ||
				canWalk(grid, x, y+1) && !canWalk(grid, x-dx, y+1) {
				return grid[x][y]
			}
		} else if dy != 0 {
			if (canWalk(grid, x-1, y) && !canWalk(grid, x-1, y-dy)) ||
				canWalk(grid, x+1, y) && !canWalk(grid, x+1, y-dy) {
				return grid[x][y]
			}
		}
	}

	// moving diagonally, must make sure one of the vertical/horizontal
	// neighbors is open to allow the path
	if canWalk(grid, x+dx, y) && canWalk(grid, x, y+dy) {
		return jump(grid, x+dx, y+dy, x, y, b)
	}
	return nil
}

// A*
func astar(wv grid, a *node, b *node) nodes {
	var (
		open, foundPath nodes
		node            *node
		ng              float64
	)

	a.g = 0
	a.f = 0

	open = append(open, a)

	a.opened = true

	for len(open) != 0 {
		open, node = lowestF(open)
		node.closed = true

		if equal(node, b) {
			break
		}

		for _, neighbor := range adjacentNodes(wv, node, 1) { // 744 609

			if neighbor.closed {
				continue
			}

			ng = float64(node.g)
			if neighbor.x-node.x == 0 || neighbor.y-node.y == 0 {
				ng += 1
			} else {
				ng += math.Sqrt2
			}

			if !neighbor.opened || ng < float64(neighbor.g) {
				neighbor.g = int(ng)
				neighbor.h = 2 * euclideanDistance(neighbor, b)
				//neighbor.h = octileDistance(neighbor, b)

				neighbor.f = neighbor.g + neighbor.h
				neighbor.parent = node

				if !neighbor.opened {
					open = append(open, neighbor)
					neighbor.opened = true
				}
			}
		}
	}

	next := node
	for next != nil {
		foundPath = append(foundPath, next)
		next = next.parent
	}

	sort.Sort(sort.Reverse(foundPath))

	return foundPath
}

// given speed, create a new node slice with speed as euclideanDistance value between nodes
// a => n+1 => b
// for each n
// 	if euclideanDistance between a, n+1 == speed
//	  	add a to new slice
//		a = n+1
func reduce(e nodes, speed int) nodes {
	var (
		rp      nodes
		current *node
	)

	current = e[0]
	for i := 0; i < len(e); i++ {
		d := euclideanDistance(current, e[i])
		if d >= speed {
			rp = append(rp, current)
			current = e[i]
		}
	}
	return rp
}

func equal(v1, v2 *node) bool {
	if v1.x == v2.x && v1.y == v2.y {
		return true
	}
	return false
}

func lowestF(open nodes) ([]*node, *node) {
	var (
		uo = make([]*node, 0)
		v  *node
	)

	sort.Sort(open)
	v = open[0]

	uo = append(uo, open[1:]...)

	return uo, v
}

func canWalk(wv grid, x, y int) bool {
	_, ok := wv[x][y]
	return ok
}

func adjacentNodes(wv grid, node *node, margin int) nodes {
	var (
		result nodes
	)

	// ↑
	if canWalk(wv, node.x, node.y-margin) {
		result = append(result, wv[node.x][node.y-margin])
	}

	// ↓
	if canWalk(wv, node.x, node.y+margin) {
		result = append(result, wv[node.x][node.y+margin])
	}

	// →
	if canWalk(wv, node.x+margin, node.y) {
		result = append(result, wv[node.x+1][node.y])
	}

	// ↓
	if canWalk(wv, node.x-margin, node.y) {
		result = append(result, wv[node.x-1][node.y])
	}

	// ↖
	if canWalk(wv, node.x-margin, node.y-margin) {
		result = append(result, wv[node.x-1][node.y-1])
	}

	// ↗
	if canWalk(wv, node.x+margin, node.y-margin) {
		result = append(result, wv[node.x+1][node.y-1])
	}

	// ↘
	if canWalk(wv, node.x+margin, node.y+margin) {
		result = append(result, wv[node.x+margin][node.y+margin])
	}

	// ↙
	if canWalk(wv, node.x-margin, node.y+margin) {
		result = append(result, wv[node.x-margin][node.y+margin])
	}

	return result
}

func initNodes(s *data.SHBD) grid {
	var (
		wv = make(grid)
		r  = bytes.NewReader(s.Data)
	)

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
						wv[rX] = make(map[int]*node)
					}

					wv[rX][rY] = &node{
						x: rX,
						y: rY,
					}
				}
			}
		}
	}
	return wv
}

func in(list []*node, sv *node) bool {
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
