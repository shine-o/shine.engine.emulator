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

	a := &vertex{
		x: 835,
		y: 700,
	}

	b := &vertex{
		x: 1070,
		y: 1540,
	}

	v := initVertices(s)

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
	// print an image with the vertex data
}

// +10
func Test_Paint_Path_Nodes(*testing.T) {
	m := "Rou"
	var s *data.SHBD
	s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", m))

	if err != nil {
		logger.Error(err)
	}

	v := initVertices(s)

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
	// print an image with the vertex data

	// paint nodes that will be used for paths.

}

func solveForTime(s, d int) int {
	return d/s
}

func solveForSpeed(d, t int) int {
	return d / t
}

func Test_Map_Intermitent_By_Speed_Path_A_B_AStar(t *testing.T)  {
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

	// I will also need to calculate the distance per second


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

	a := &vertex{
		x: 515,
		y: 1177,
	}

	b := &vertex{
		x: 1250,
		y: 1500,
	}

	v := initVertices(s)

	pn := initPathNodes(s, v)

	if err != nil {
		logger.Fatal(err)
	}

	pathVertices := astar(pn, a, b)

	pathVertices = reduce(pathVertices,15)

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

func Benchmark_ReduceVertices(b *testing.B)  {

	m := "Rou"
	var s *data.SHBD
	s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", m))

	if err != nil {
		logger.Error(err)
	}

	pa := &vertex{
		x: 515,
		y: 1177,
	}

	pb := &vertex{
		x: 1250,
		y: 1500,
	}

	v := initVertices(s)

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

	var v walkableVertices
	b.Run("initVertices", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v = initVertices(s)
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

	v := initVertices(s)
	pn := initPathNodes(s, v)

	if err != nil {
		logger.Fatal(err)
	}

	pa := &vertex{
		x: 835,
		y: 700,
	}

	pb := &vertex{
		x: 1070,
		y: 1540,
	}


	b.Run("astar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			astar(pn, pa, pb)
		}
	})
}

func Benchmark_jjs_algorithm(b *testing.B) {
	// base astar algo
	// find successors
	// foreach neighbour, jump
}

func SHBDToImage(s *data.SHBD, wv walkableVertices) (*image.RGBA, error) {
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
					count := len(adjacentVertices(wv, &vertex{
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

func initPathNodes(s *data.SHBD, wv walkableVertices) walkableVertices {

	var (
		pn = make(walkableVertices)
		r = bytes.NewReader(s.Data)
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
					count := len(adjacentVertices(wv, &vertex{
						x: rX,
						y: rY,
					}, nodesMargin))

					if count == neighborNodes {	// do not add nodes if they are too close to inaccessible nodes
						_, ok := pn[rX]
						if !ok {
							pn[rX] = make(map[int]*vertex)
						}

						pn[rX][rY] = &vertex{
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
	nodesMargin = 7
)

func distance(a, b *vertex) int {
	d := math.Sqrt(math.Pow(float64(a.x-b.x), 2) + math.Pow(float64(a.y-b.y), 2))
	return int(d)
}

type walkableVertices map[int]map[int]*vertex

type vertex struct {
	parent  *vertex
	x, y    int
	h, g, f int
	opened  bool
	closed  bool
}

type vertices []*vertex

func (e vertices) Len() int {
	return len(e)
}

func (e vertices) Less(i, j int) bool {
	return e[i].h < e[j].h
}

func (e vertices) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// given speed, create a new vertex slice with speed as distance value between nodes
// a => n+1 => b
// for each n
// 	if distance between a, n+1 == speed
//	  	add a to new slice
//		a = n+1
func reduce(e vertices, speed int) vertices {
	var (
		rp vertices
		current * vertex
	)

	current = e[0]
	for i := 0; i < len(e); i++ {
		d := distance(current, e[i])
		if d >= speed {
			rp = append(rp, current)
			current = e[i]
		}
	}
	return rp
}

func astar(wv walkableVertices, a *vertex, b *vertex) vertices {
	var (
		open, shortestPath vertices
		node               *vertex
		ng                 float64
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

		for _, neighbor := range adjacentVertices(wv, node, 1) { // 744 609

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
				neighbor.h = 2 * distance(neighbor, b)
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
		shortestPath = append(shortestPath, next)
		next = next.parent
	}

	sort.Sort(sort.Reverse(shortestPath))

	return shortestPath
}

func equal(v1, v2 *vertex) bool {
	if v1.x == v2.x && v1.y == v2.y {
		return true
	}
	return false
}

func lowestF(open vertices) ([]*vertex, *vertex) {
	var (
		uo = make([]*vertex, 0)
		v  *vertex
	)

	sort.Sort(open)
	v = open[0]

	uo = append(uo, open[1:]...)

	return uo, v
}

func in(list []*vertex, sv *vertex) bool {
	for _, v := range list {
		if equal(v, sv) {
			return true
		}
	}
	return false
}

func canWalk(wv walkableVertices, x, y int) bool {
	_, ok := wv[x][y]
	return ok
}

func adjacentVertices(wv walkableVertices, node *vertex, margin int) vertices {
	var (
		result vertices
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

func initVertices(s *data.SHBD) walkableVertices {
	var (
		wv = make(walkableVertices)
		r = bytes.NewReader(s.Data)
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
						wv[rX] = make(map[int]*vertex)
					}

					wv[rX][rY] = &vertex{
						x: rX,
						y: rY,
					}
				}
			}
		}
	}
	return wv
}
