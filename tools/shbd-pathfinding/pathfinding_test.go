package main

import (
	"bytes"
	"fmt"
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

	a := &vertex{
		x: 1596,
		y: 1320,
	}

	b := &vertex{
		x: 1628,
		y: 591,
	}

	vertices := initVertices(s)

	res := astar(vertices, a, b)

	//SHBDToImage(s)

	//vertices := Vertices(s)
	//
	logger.Info(res)
	//
	//for i, v1 := range vertices {
	//	if (i+1) > len(vertices) {
	//		break
	//	}
	//
	//	v2 := vertices[i+1]
	//	fmt.Println(distance(v1.x, v2.x, v1.y, v2.y))
	//	fmt.Println(distance(v1.x, v1.x, v1.y, v1.y))
	//}

}

func BenchmarkName(b *testing.B) {
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

func distance(a, b * vertex) int {
	d := math.Sqrt(math.Pow(float64(a.x-b.x), 2) + math.Pow(float64(a.y-b.y), 2))
	return int(d)
}


type heuristics struct {
}

// x, y,
type vertices map[int]map[int]*vertex

type vertex struct {
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
func astar(vertices vertices, a * vertex, b * vertex) []*vertex {
	var (
		open, closed []*vertex
		h, g, f int
		current * vertex
	)
	h = distance(a, b)
	f = 0 + h
	a.f = f
	current = a
	for *current != *b {

		if current == nil {
			logger.Error("nil vertex")
		}

		adjacentVertices := adjacentVertices(vertices, current)
		if len(adjacentVertices) == 0 {
			logger.Fatal("no adjacent nodes, impossible to proceed")
		}

		for _, av := range adjacentVertices {
			if !in(closed, av) && !in(open, av) {
				open = append(open, av)
				g = distance(av, a)
				h = distance(av, b)
				nf := g + h
				av.f = nf

				if nf <= f {
					f = nf
					a = av
				}
			}
		}

		closed = append(closed, current)

		open, current = lowestF(open)

		if current.f == 1 {
			current.f  = 0
		}
	}

	return closed
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

		if v.f <= open[i+1].f {
			vertex = v
			index = i
		}
	}

	//ret = append(ret, s[:index]...)
	uo = append(uo, open[:index]...)
	uo = append(uo, open[index+1:]...)

	return uo, vertex
}

func in(list []*vertex, sv *vertex) bool {
	for _ , v := range list {
		if v == sv {
			return true
		}
	}
	return false
}

func adjacentVertices(v vertices, current * vertex) []*vertex {
	var result []*vertex

	if current == nil {
		logger.Fatal("nil vertex")
	}

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

	return result
}

// A*
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