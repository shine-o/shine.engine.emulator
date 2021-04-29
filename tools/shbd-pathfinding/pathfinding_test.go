package main

import (
	"bytes"
	"fmt"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"image/color"
	"math"
	"sort"
	"testing"
)

func Test_Path_A_B2(t *testing.T) {
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
		x: 734,
		y: 612,
	}

	b := &vertex{
		x: 1500,
		y: 1440,
	}

	v := initVertices(s)

	if err != nil {
		logger.Fatal(err)
	}

	pathVertices := astar(v, a, b)

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

func Benchmark_InitVertices(b *testing.B) {
	m := "Rou"
	var s *data.SHBD
	s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", m))

	if err != nil {
		logger.Error(err)
	}

	for i := 0; i < b.N; i++ {
		initVertices(s)
	}
}

func Benchmark_aStar_Algorithm(b *testing.B) {
	m := "Rou"
	var s *data.SHBD
	s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", m))

	if err != nil {
		logger.Error(err)
	}

	v := initVertices(s)

	pa := &vertex{
		x: 106,
		y: 510,
	}

	pb := &vertex{
		x: 1662,
		y: 1111,
	}


	for i := 0; i < b.N; i++ {
		astar(v, pa, pb)
	}
}

func distance(a, b * vertex) int {
	d := math.Sqrt(math.Pow(float64(a.x - b.x), 2) + math.Pow(float64(a.y - b.y), 2))
	return int(d)
}

// x, y,
type walkableVectors map[int]map[int]*vertex

type vertex struct {
	parent  *vertex
	x, y    int
	h, g, f int
	opened  bool
	closed  bool
}

type vertexSegment []*vertex

func (e vertexSegment) Len() int {
	return len(e)
}

func (e vertexSegment) Less(i, j int) bool {
	return e[i].h < e[j].h
}

func (e vertexSegment) Swap(i, j int) {
	e[i], e[j] = e[j],e[i]
}

func astar(wv walkableVectors, a * vertex, b * vertex) vertexSegment {
	var (
		open, shortestPath vertexSegment
		//node, neighbor             * vertex
		node				          * vertex
		//neighbors                  vertexSegment
		ng                   float64
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

		for _, neighbor := range adjacentVertices(wv, node) { // 744 609

			if neighbor.closed {
				continue
			}

			ng = float64(node.g)
			if neighbor.x - node.x == 0 || neighbor.y - node.y == 0 {
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

func equal(v1, v2 *vertex) bool  {
	if v1.x == v2.x && v1.y == v2.y	{
		return true
	}
	return false
}

func lowestF(open vertexSegment) ([]*vertex, *vertex) {
	var (
		uo = make([]*vertex, 0)
		v *vertex
	)

	sort.Sort(open)
	v = open[0]

	uo = append(uo, open[1:]...)

	return uo, v
}

func in(list []*vertex, sv *vertex) bool {
	for _ , v := range list {
		if equal(v, sv) {
			return true
		}
	}
	return false
}

func canWalk(wv walkableVectors, x, y int) bool {
	_, ok := wv[x][y]
	return ok
}

func adjacentVertices(wv walkableVectors, node * vertex) vertexSegment {
	var result vertexSegment

	if canWalk(wv, node.x, node.y-1) {
		result = append(result, wv[node.x][node.y-1])
	}

	if canWalk(wv, node.x, node.y+1) {
		result = append(result, wv[node.x][node.y+1])
	}

	if canWalk(wv, node.x+1, node.y) {
		result = append(result, wv[node.x+1][node.y])
	}

	if canWalk(wv, node.x-1, node.y) {
		result = append(result, wv[node.x-1][node.y])
	}

	return result
}

func initVertices(s *data.SHBD) walkableVectors {

	var vertices = make(walkableVectors)

	r := bytes.NewReader(s.Data)

	for y := 0; y < s.Y; y++ {
		for x := 0; x < s.X; x++ {
			b, err := r.ReadByte()
			if err != nil {
				return vertices
			}

			for i := 0; i < 8; i++ {
				if b&byte(math.Pow(2, float64(i))) == 0 {
					rX := x*8 + i
					rY := y

					_, ok := vertices[rX]
					if !ok {
						vertices[rX] = make(map[int]*vertex)
					}

					vertices[rX][rY] = &vertex{
						x: rX,
						y: rY,
					}
				}
			}

		}
	}

	return vertices
}
