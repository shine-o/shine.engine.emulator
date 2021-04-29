package main

import (
	"bytes"
	"fmt"
	"image/color"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"math"
	"testing"
)

func Test_Path_A_B(t *testing.T) {
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
		x: 1036,
		y: 766,
	}

	b := &vertex{
		x: 1045,
		y: 766,
	}

	vertices := initVertices(s)

	pathVertices := astar(vertices, a, b)
	//pathVertices := astar2(vertices, a, b)

	for _, v := range pathVertices {
		img.Set(v.x, v.y, color.RGBA{
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

	vertices := initVertices(s)

	pa := &vertex{
		x: 106,
		y: 510,
	}

	//pb := &vertex{
	//	x: 1597,
	//	y: 1321,
	//}

	pb := &vertex{
		x: 1662,
		y: 1111,
	}


	for i := 0; i < b.N; i++ {
		astar(vertices, pa, pb)
	}
}

func distance(a, b * vertex) int {
	d := math.Sqrt(math.Pow(float64(a.x - b.x), 2) + math.Pow(float64(a.y - b.y), 2))
	//d := math.Abs(float64(a.x - b.x)) + math.Abs(float64(a.y - b.y))
	return int(d)
}

// x, y,
type vertices map[int]map[int]*vertex

type vertex struct {
	prev *vertex
	x, y int
	h, g, f int
}

//*** A* pseudocode***
//
//Initialise open and closed lists
//Make the start vertex current
//Calculate heuristic distance of start vertex to destination (h)
//Calculate f value for start vertex (f = g + h, where g = 0)
//WHILE current vertex is not the destination
//    FOR each vertex adjacent to current
//        IF vertex not in closed list and not in open list THEN
//            Add vertex to open list
//        END IF
//        Calculate distance from start (g)
//        Calculate heuristic distance to destination (h)
//        Calculate f value (f = g + h)
//        IF new f value < existing f value or there is no existing f value THEN
//            Update f value
//            Set parent to be the current vertex
//        END IF
//    NEXT adjacent vertex
//    Add current vertex to closed list
//    Remove vertex with lowest f value from open list and make it current
//END WHILE
// https://www.youtube.com/watch?v=eSOJ3ARN5FM
func astar(vertices vertices, a * vertex, b * vertex) []*vertex {
	var (
		//open, closed, shortestPath []*vertex
		open, closed []*vertex
		current * vertex
		f int
	)

	a.h = distance(a, b)
	a.f = 0 + a.h

	current = a
	for !equal(current, b) {

		adjacentVertices := adjacentVertices(vertices, current)

		if len(adjacentVertices) == 0 {
			logger.Fatal("no adjacent nodes, impossible to proceed")
		}

		for _, av := range adjacentVertices {
			if !in(closed, av) && !in(open, av) {
				open = append(open, av)
			}

			av.g = distance(av, a)
			av.h = distance(av, b)
			av.f = av.g + av.h

			if av.f <= f || f == 0 {
				f = av.f
				av.prev = current
			}
		}

		closed = append(closed, current)

		open, current = lowestF(open)

		if current == nil {
			logger.Error("nil vertex")
			break
		}
	}
	//
	//if current == nil {
	//	logger.Error("nil vertex")
	//	return shortestPath
	//}


	//next := current
	//for next != nil {
	//	shortestPath = append(shortestPath, next.prev)
	//	next = next.prev
	//}

	return closed
}

func equal(v1, v2 *vertex) bool  {
	if v1.x == v2.x && v1.y == v2.y	{
		return true
	}
	return false
}

func lowestF(open []*vertex) ([]*vertex, *vertex) {
	var (
		uo = make([]*vertex, 0)
		vertex *vertex
		index int
	)

	for i, v := range open {

		if i+1 >= len(open) {
			break
		}

		if v.f < open[i+1].f && v.h < open[i+1].h {
			vertex = v
			index = i
		}
	}

	uo = append(uo, open[:index]...)
	uo = append(uo, open[index+1:]...)

	return uo, vertex
}

func in(list []*vertex, sv *vertex) bool {
	for _ , v := range list {
		if equal(v, sv) {
			return true
		}
	}
	return false
}

func adjacentVertices(v vertices, current * vertex) []*vertex {
	var result []*vertex

	if current != nil {
		upv, ok := v[current.x][current.y+1]
		if ok {
			result = append(result, upv)
		}

		dpv, ok := v[current.x][current.y-1]
		if ok {
			result = append(result, dpv)
		}

		rpv, ok := v[current.x+1][current.y]

		if ok {
			result = append(result, rpv)
		}

		lpv, ok := v[current.x-1][current.y]

		if ok {
			result = append(result, lpv)
		}
	} else {
		logger.Fatal("nil vertex")
	}

	return result
}

func initVertices(s *data.SHBD) vertices {

	var vertices = make(vertices)

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