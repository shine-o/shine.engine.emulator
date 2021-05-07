package main

import (
	"bytes"
	"fmt"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"image"
	"image/color"
	"math"
	"os"
	"sort"
	"sync"
	"testing"
)

var (
	mn string
	s  *data.SHBD
	v  grid
	vNodesNoWallsMargin = 8
	vNodesWithWallsMargin = 8
	vNodesWithWalls  grid
	vNodesNoWalls    grid
)

func TestMain(m *testing.M) {
	mn = "Rou"
	//s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", m))
	s, _ = data.LoadSHBDFile(fmt.Sprintf("/home/marius/projects/shine/shine.engine.emulator/files/blocks/%v.shbd", mn))
	v = rawGrid(s)
	vNodesNoWalls = gridWithReducedNodes(s, vNodesNoWallsMargin)
	vNodesWithWalls = gridWithWallsMargin(s, v, vNodesWithWallsMargin)
	os.Exit(m.Run())
}

func Test_Path_A_astar_paint(t *testing.T) {
	v1 := copyGrid(vNodesNoWalls)

	img1, err := data.SHBDToImage(s)

	if err != nil {
		logger.Error(err)
	}

	pathVertices1 := astar(v1, v, 835, 700, 1070, 1540, vNodesNoWallsMargin)
	fmt.Printf("vNodesNoWallsMargin path nodes %v\n", len(pathVertices1))

	for _, pv := range pathVertices1 {
		img1.Set(pv.x, pv.y, color.RGBA{
			R: 207,
			G: 0,
			B: 15,
			A: 1,
		})
	}

	err = data.SaveBmpFile(img1, "./", "path_nodes_no_walls")

	if err != nil {
		logger.Error(err)
	}
	// print an image with the node data

	v2 := copyGrid(vNodesWithWalls)

	img2, err := data.SHBDToImage(s)

	if err != nil {
		logger.Error(err)
	}

	pathVertices2 := astar(v2, v,835, 700, 1070, 1540, vNodesWithWallsMargin)
	fmt.Printf("vNodesWithWalls path nodes %v\n", len(pathVertices2))

	for _, pv := range pathVertices2 {
		img2.Set(pv.x, pv.y, color.RGBA{
			R: 207,
			G: 0,
			B: 15,
			A: 1,
		})
	}

	err = data.SaveBmpFile(img2, "./", "path_nodes_with_walls")

	if err != nil {
		logger.Error(err)
	}
	//

	v3 := copyGrid(v)

	img3, err := data.SHBDToImage(s)

	if err != nil {
		logger.Error(err)
	}

	pathVertices3 := astar(v3,v, 835, 700, 1070, 1540, 1)

	fmt.Printf("raw path nodes %v\n", len(pathVertices3))

	for _, pv := range pathVertices3 {
		img3.Set(pv.x, pv.y, color.RGBA{
			R: 207,
			G: 0,
			B: 15,
			A: 1,
		})
	}

	err = data.SaveBmpFile(img3, "./", "path_raw")

	if err != nil {
		logger.Error(err)
	}
}

func Test_Paint_Path_Nodes(*testing.T) {
	var (
		m = "Rou"
		//s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", m))
		s, err = data.LoadSHBDFile(fmt.Sprintf("/home/marius/projects/shine/shine.engine.emulator/files/blocks/%v.shbd", m))
	)

	if err != nil {
		logger.Error(err)
	}

	v := rawGrid(s)

	if err != nil {
		logger.Fatal(err)
	}

	img, err := PaintNodesAndWallMargins(s, v)

	if err != nil {
		logger.Error(err)
	}

	err = data.SaveBmpFile(img, "./", "painted_nodes_and_wall_margins")

	if err != nil {
		logger.Error(err)
	}

	img, err = PaintNodesWithoutWallMargins(s, v)

	if err != nil {
		logger.Error(err)
	}

	err = data.SaveBmpFile(img, "./", "painted_nodes_without_wall_margins")

	// print an image with the node data

	// paint nodes that will be used for paths.
}

func Test_Paint_Path_Nodes_Multiple(*testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			img, err := data.SHBDToImage(s)

			if err != nil {
				logger.Error(err)
			}

			cgrid := copyGrid(vNodesNoWalls)

			if err != nil {
				logger.Fatal(err)
			}

			pathVertices := astar(cgrid,v, 835, 700, 1070, 1540, 1)

			for _, pv := range pathVertices {
				img.Set(pv.x, pv.y, color.RGBA{
					R: 207,
					G: 0,
					B: 15,
					A: 1,
				})
			}

			err = data.SaveBmpFile(img, "./", name)

			if err != nil {
				logger.Fatal(err)
			}

		}(fmt.Sprintf("paintedpath%v", i))
	}
	wg.Wait()

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

	pn := gridWithReducedNodes(s, 4)

	if err != nil {
		logger.Fatal(err)
	}

	pathVertices := astar(pn,v, 835, 700, 1070, 1540, 1)

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

func Benchmark_algorithms(b *testing.B) {
	b.Run("astar_raw", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ng := copyGrid(v)
			astar(ng, v, 835, 700, 1070, 1540, 1)
		}
	})

	b.Run("astar_reduced_nodes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ng := copyGrid(vNodesNoWalls)
			astar(ng, v, 835, 700, 1070, 1540, vNodesNoWallsMargin)
		}
	})

	b.Run("astar_reduced_nodes_with_wall_margin", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ng := copyGrid(vNodesWithWalls)
			astar(ng,v, 835, 700, 1070, 1540, vNodesWithWallsMargin)
		}
	})
}

func Benchmark_InitVertices(b *testing.B) {
	var (
		m = "Rou"
		//s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", m))
		s, err = data.LoadSHBDFile(fmt.Sprintf("/home/marius/projects/shine/shine.engine.emulator/files/blocks/%v.shbd", m))
	)
	if err != nil {
		logger.Error(err)
	}

	b.Run("rawGrid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v = rawGrid(s)
		}
	})

	b.Run("pathNodes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gridWithReducedNodes(s, 4)
		}
	})
}

func Benchmark_ReduceVertices(b *testing.B) {
	m := "Rou"
	var s *data.SHBD
	s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", m))

	if err != nil {
		logger.Error(err)
	}

	pn := gridWithReducedNodes(s, 4)

	if err != nil {
		logger.Fatal(err)
	}

	pathVertices := astar(pn, v, 835, 700, 1070, 1540, 1)
	b.Run("reduce", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = reduce(pathVertices, 15)
		}
	})
}

func Benchmark_Copy_Grid(b *testing.B) {
	var (
		m = "Rou"
		//s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", m))
		s, err = data.LoadSHBDFile(fmt.Sprintf("/home/marius/projects/shine/shine.engine.emulator/files/blocks/%v.shbd", m))
	)

	if err != nil {
		logger.Error(err)
	}

	v := rawGrid(s)

	pn := gridWithReducedNodes(s, 4)

	if err != nil {
		logger.Fatal(err)
	}

	b.Run("copyGrid_raw", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			copyGrid(v)
		}
	})

	b.Run("copyGrid_nodes_path", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			copyGrid(pn)
		}
	})
}

func Benchmark_Close_Distance_Algorithms(b *testing.B) {

}

func PaintNodesAndWallMargins(s *data.SHBD, wv grid) (*image.RGBA, error) {
	var (
		r = bytes.NewReader(s.Data)

		img = image.NewRGBA(image.Rectangle{
			Min: image.Point{
				X: 0,
				Y: 0,
			},
			Max: image.Point{
				X: s.X * 8,
				Y: s.Y,
			},
		})
		count = 0
	)

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
					countn := len(adjacentNodes(wv, v, rX, rY, 8))
					if rX%vNodesWithWallsMargin == 0 && rY%vNodesWithWallsMargin == 0 {
						if countn >= neighborNodes {
								c = color.RGBA{
									R: 255,
									G: 0,
									B: 0,
									A: 1,
								}
								count++
							} else {
								c = color.White
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
	fmt.Printf("count nodes %v\n", count)
	return img, nil
}

func PaintNodesWithoutWallMargins(s *data.SHBD, wv grid) (*image.RGBA, error) {
	var (
		r = bytes.NewReader(s.Data)

		img = image.NewRGBA(image.Rectangle{
			Min: image.Point{
				X: 0,
				Y: 0,
			},
			Max: image.Point{
				X: s.X * 8,
				Y: s.Y,
			},
		})
		count = 0
	)

	for y := 0; y < s.Y; y++ {
		for x := 0; x < s.X; x++ {
			b, err := r.ReadByte()
			if err != nil {
				return img, err
			}

			// every ten pixels
			for i := 0; i < 8; i++ {
				var (
					rX, rY    int
					c         color.Color
				)

				rX = x*8 + i
				rY = y

				if b&byte(math.Pow(2, float64(i))) == 0 {
					if rX%vNodesNoWallsMargin == 0 && rY%vNodesNoWallsMargin == 0 {
						c = color.RGBA{
							R: 255,
							G: 0,
							B: 0,
							A: 1,
						}
						count++
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
	fmt.Printf("count nodes %v\n", count)
	return img, nil
}

func rawGrid(s *data.SHBD) grid {
	var (
		wv    = make(grid)
		r     = bytes.NewReader(s.Data)
		count = 0
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
					n := &node{
						x: rX,
						y: rY,
					}
					wv[rX][rY] = n
					count++
				}
			}
		}
	}
	fmt.Printf("count nodes %v \n", count)
	return wv
}

func gridWithReducedNodes(s *data.SHBD, margin int) grid {
	var (
		pn    = make(grid)
		r     = bytes.NewReader(s.Data)
		count = 0
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
					// margin between nodes
					if rX%margin == 0 && rY%margin == 0 {
						_, ok := pn[rX]
						if !ok {
							pn[rX] = make(map[int]*node)
						}

						pn[rX][rY] = &node{
							x: rX,
							y: rY,
						}
						count++
					}
				}
			}
		}
	}
	fmt.Printf("count nodes %v \n", count)
	return pn
}

// preset nodes + extra walls so that NPCs don't get too close to collision areas
func gridWithWallsMargin(s *data.SHBD, wv grid, margin int) grid {
	var (
		pn    = make(grid)
		r     = bytes.NewReader(s.Data)
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
					// add only half the nodes needed
					count := len(adjacentNodes(wv, v, rX, rY, margin))
					if count >= neighborNodes { // do not add nodes if they are too close to inaccessible nodes
						if rX%margin == 0 && rY%margin == 0 {

							_, ok := pn[rX]
							if !ok {
								pn[rX] = make(map[int]*node)
							}

							pn[rX][rY] = &node{
								x: rX,
								y: rY,
							}
							ncount++
						}
					}
				}
			}
		}
	}
	fmt.Printf("count nodes %v \n", ncount)
	return pn
}

const (
	neighborNodes = 8
	nodesMargin   = 2
)

// sqrt(dx * dx + dy * dy)
func euclideanDistance(fx, fy, tx, ty int) int {
	dx := tx - fx
	dy := ty - fy
	d := math.Sqrt(float64(dx*dx + dy*dy))
	return int(d)
}

// sqrt(dx * dx + dy * dy)
func octileDistance(a, b *node) int {
	var F = math.Sqrt2 - 1
	dx := math.Abs(float64(a.x - b.x))
	dy := math.Abs(float64(a.y - b.y))

	if dx < dy {
		return int(F*dx + dy)
	}

	return int(F*dy + dx)
}

type grid map[int]map[int]*node

type node struct {
	parent     *node
	x, y       int
	h, g, f    int
	opened     bool
	closed     bool
	neighbours nodes
}

type nodes []*node

func (e nodes) Len() int {
	return len(e)
}

func (e nodes) Less(i, j int) bool {
	return e[i].f < e[j].f
}

func (e nodes) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// use K nearest neighbor algorithm
func (g grid) getNearest(x, y int) *node {
	// for each rectangle given the epicenter x,y
	// send a go routine to iterate over the rectangle
	// if neighbour is found, send to channel and exit function
	// close channel so all other routines are canceled
	var (
		perimeter = 150
		ch = make(chan *node)
	)

	go searchOneTopLeft(g, ch, x, y, perimeter)
	go searchOneTopRight(g, ch, x, y, perimeter)
	go searchOneBottomLeft(g, ch, x, y, perimeter)
	go searchOneBottomRight(g, ch, x, y, perimeter)

	return <- ch
}

//func (g grid) getNearest(x, y int) *node {
//	// given this point, move forward in all directions until a node is found
//	// sounds good, doesn't work
//	//return nearestNodes(g, x, y , 1,256)
//
//	// I should try to build a square around the point
//	// then iterate each node in the square, if ts a walkable position
//	// add to the list ( or return if its needed to reduce resources, add this as a function param lazyMode)
//	// iterate each item in the list, calculate euclidean distance of node, return the one with the shortest distance
//
//
//	// for each rectangle given the epicenter x,y
//	// send a go routine to iterate over the rectangle
//	// if neighbour is found, send to channel and exit function
//	// close channel so all other routines are canceled
//	var (
//		// boundary x,y
//		bx, by int
//		perimeter = 150
//	)
//
//	// upper left side
//	// x-50, y-50
//	// for y++
//	//		for x++
//	//		if x,y is a node
//	bx = x-perimeter
//	by = y-perimeter
//	for cy := y; cy != by; cy-- {
//		for cx := x; cx != bx; cx-- {
//			n := g.get(cx, cy)
//			if n != nil{
//				return n
//			}
//		}
//	}
//	// upper right side
//	// x+50, y-50
//	// for y++
//	//		for x++
//	//		if x,y is a node
//	bx = x+perimeter
//	by = y-perimeter
//	for cy := y; cy != by; cy-- {
//		for cx := x; cx != bx; cx++ {
//			n := g.get(cx, cy)
//			if n != nil{
//				return n
//			}
//		}
//	}
//
//	// bottom left side
//	// x-50, y+50
//	bx = x-perimeter
//	by = y+perimeter
//	for cy := y; cy != by; cy++ {
//		for cx := x; cx != bx; cx-- {
//			n := g.get(cx, cy)
//			if n != nil{
//				return n
//			}
//		}
//	}
//
//	// bottom left side
//	// x+50, y+50
//	bx = x+perimeter
//	by = y+perimeter
//	for cy := y; cy != by; cy++ {
//		for cx := x; cx != bx; cx++ {
//			n := g.get(cx, cy)
//			if n != nil{
//				return n
//			}
//		}
//	}
//	return nil
//}

func searchOneTopLeft(g grid, ch chan<- *node, x, y, perimeter int) {
	bx := x-perimeter
	by := y-perimeter
	for cy := y; cy != by; cy-- {
		for cx := x; cx != bx; cx-- {
			n := g.get(cx, cy)
			if n != nil{
				select {
				case ch <- n:
				default:
					return
				}
			}
		}
	}
}

func searchOneTopRight(g grid, ch chan<- *node, x, y, perimeter int) {
	bx := x+perimeter
	by := y-perimeter
	for cy := y; cy != by; cy-- {
		for cx := x; cx != bx; cx++ {
			n := g.get(cx, cy)
			if n != nil {
				select {
				case ch <- n:
				default:
					return
				}
			}
		}
	}
}

func searchOneBottomLeft(g grid, ch chan<- *node, x, y, perimeter int) {
	bx := x-perimeter
	by := y+perimeter
	for cy := y; cy != by; cy++ {
		for cx := x; cx != bx; cx-- {
			n := g.get(cx, cy)
			if n != nil{
				ch <- n
				select {
				case ch <- n:
				default:
					return
				}
			}
		}
	}
}

func searchOneBottomRight(g grid, ch chan<- *node, x, y, perimeter int) {
	bx := x+perimeter
	by := y+perimeter
	for cy := y; cy != by; cy++ {
		for cx := x; cx != bx; cx++ {
			n := g.get(cx, cy)
			if n != nil {
				select {
				case ch <- n:
				default:
					return
				}
			}
		}
	}
}


//
func getAllTopLeft(g grid, ch chan<- *node, done chan <- bool, x, y, perimeter int) {
	bx := x-perimeter
	by := y-perimeter
	for cy := y; cy != by; cy-- {
		for cx := x; cx != bx; cx-- {
			n := g.get(cx, cy)
			if n != nil{
				ch <- n
			}
		}
	}
	done <- true
}

func getAllTopRight(g grid, ch chan<- *node, done chan <- bool, x, y, perimeter int) {
	bx := x+perimeter
	by := y-perimeter
	for cy := y; cy != by; cy-- {
		for cx := x; cx != bx; cx++ {
			n := g.get(cx, cy)
			if n != nil{
				ch <- n
			}
		}
	}
	done <- true
}

func getAllBottomLeft(g grid, ch chan<- *node, done chan <- bool, x, y, perimeter int) {
	bx := x-perimeter
	by := y+perimeter
	for cy := y; cy != by; cy++ {
		for cx := x; cx != bx; cx-- {
			n := g.get(cx, cy)
			if n != nil{
				ch <- n
			}
		}
	}
	done <- true
}

func getAllBottomRight(g grid, ch chan<- *node, done chan <- bool, x, y, perimeter int) {
	bx := x+perimeter
	by := y+perimeter
	for cy := y; cy != by; cy++ {
		for cx := x; cx != bx; cx++ {
			n := g.get(cx, cy)
			if n != nil{
				ch <- n
			}
		}
	}
	done <- true
}

func (g grid) get(x, y int) *node {
	n, ok := g[x][y]
	if ok {
		return n
	}
	return nil
}

// A*
//func astar(wv grid, fx, fy, tx, ty, margin int) nodes {
//	var (
//		open, foundPath nodes
//		node            *node
//		ng              float64
//		a               = wv.get(fx, fy)
//		b               = wv.get(tx, ty)
//		//start = time.Now().UnixNano()
//	)
//
//	if a == nil {
//		a = wv.getNearest(fx, fy)
//	}
//
//	if b == nil {
//		b = wv.getNearest(tx, ty)
//	}
//
//	if a == nil || b == nil {
//		logger.Fatalf("a or b is nil %v %v", a, b)
//		return nil
//	}
//
//	a.g = 0
//	a.f = 0
//
//	open = append(open, a)
//
//	a.opened = true
//
//	for len(open) != 0 {
//		open, node = lowestF(open)
//		node.closed = true
//
//		if equal(node, b) {
//			break
//		}
//
//		for _, neighbor := range adjacentNodes(wv, node.x, node.y, margin) {
//			if neighbor.closed {
//				continue
//			}
//
//			//ng = float64(node.g)
//			//if neighbor.x-node.x == 0 || neighbor.y-node.y == 0 {
//			//	ng += 1
//			//} else {
//			//	ng += math.Sqrt2
//			//}
//
//			ng = float64(node.g)
//			if neighbor.x-node.x == 0 || neighbor.y-node.y == 0 {
//				ng += float64(margin/2)
//			} else {
//				ng += math.Sqrt2
//			}
//
//			if !neighbor.opened || ng < float64(neighbor.g) {
//				neighbor.g = int(ng)
//				//neighbor.h = octileDistance(neighbor, b)
//				//neighbor.h = manhatanDistance(neighbor.x, neighbor.y, b.x, b.y)
//				neighbor.h = euclideanDistance(neighbor.x, neighbor.y, b.x, b.y)
//				neighbor.f = neighbor.g + neighbor.h
//				neighbor.parent = node
//
//				if !neighbor.opened {
//					open = append(open, neighbor)
//					neighbor.opened = true
//				}
//			}
//		}
//	}
//
//	//fmt.Println(time.Now().UnixNano()-start)
//	next := node
//	for next != nil {
//		foundPath = append(foundPath, next)
//		next = next.parent
//	}
//
//	sort.Sort(sort.Reverse(foundPath))
//
//	return foundPath
//}

type nodeList struct {
	list nodes
	sync.RWMutex
}

func TestAllNearbyNodes(t *testing.T) {
	//pathVertices3 := astar(v3, 835, 700, 1070, 1540, 1)
	nodes := allNearbyNodes(vNodesWithWalls, 835, 700, 100)
	fmt.Println(nodes)
}

func Benchmark_AllNearbyNodes(b *testing.B) {
	b.Run("nodes_with_walls_10_perimeter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			allNearbyNodes(vNodesWithWalls, 835, 700, 10)
		}
	})
	b.Run("nodes_with_walls_100_perimeter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			allNearbyNodes(vNodesWithWalls, 835, 700, 100)
		}
	})
	b.Run("raw_10_perimeter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			allNearbyNodes(v, 835, 700, 10)
		}
	})
	b.Run("raw_100_perimeter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			allNearbyNodes(v, 835, 700, 100)
		}
	})
}

func allNearbyNodes(g grid, x, y int, perimeter int) nodes {
	var (
		list nodes
		ch = make(chan *node, 1500)
		tld = make(chan bool)
		trd = make(chan bool)
		bld = make(chan bool)
		brd = make(chan bool)
		//sem = make(chan int, 100)
		//wg sync.WaitGroup
		topLeft, topRight, bottomLeft, bottomRight bool
	)

	go getAllTopLeft(g, ch, tld, x, y, perimeter)
	go getAllTopRight(g, ch, trd, x, y, perimeter)
	go getAllBottomLeft(g, ch, bld, x, y, perimeter)
	go getAllBottomRight(g, ch, brd, x, y, perimeter)

	for {
		select {
		case n := <- ch:
			if !in(list, n) {
				list = append(list, n)
			}
			//wg.Add(1)
			//sem <- 1
			//go func(open *nodeList, n * node) {
			//	defer wg.Done()
			//	n.g = euclideanDistance(n.x, n.y, b.x, b.y)
			//	n.h = euclideanDistance(n.x, n.y, a.x, a.y)
			//	n.f = n.g+n.h
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

func astar(ng, rg grid, fx, fy, tx, ty, margin int) nodes {
	var (
		//open, foundPath nodes
		open, foundPath nodes
		cnode     *node
		a         = ng.get(fx, fy)
		b         = ng.get(tx, ty)
		//start = time.Now().UnixNano()
	)

	if a == nil {
		a = ng.getNearest(fx, fy)
	}

	if b == nil {
		b = ng.getNearest(tx, ty)
	}

	if a == nil || b == nil {
		logger.Fatalf("a or b is nil %v %v", a, b)
		return nil
	}

	a.g = 0
	a.f = 0

	open = append(open, a)

	a.opened = true

	for {
		open, cnode = lowestF(open)

		if equal(cnode, b) {
			break
		}

		for _, neighbor := range adjacentNodes(ng, rg, cnode.x, cnode.y, margin) {
			if neighbor.closed {
				continue
			}

			ng := cnode.g + euclideanDistance(cnode.x, cnode.y, neighbor.x, neighbor.y)

			if !neighbor.opened || ng < neighbor.g {
				neighbor.g = ng
				//neighbor.h = octileDistance(neighbor, b)
				//neighbor.h = manhatanDistance(neighbor.x, neighbor.y, b.x, b.y)
				neighbor.h = euclideanDistance(neighbor.x, neighbor.y, b.x, b.y)
				neighbor.f = neighbor.g + neighbor.h
				neighbor.parent = cnode

				if !neighbor.opened {
					// if diagonally no obstacle
					neighbor.opened = true
					open = append(open, neighbor)
				}
			}
		}
	}

	cnode.closed = true

	//fmt.Println(time.Now().UnixNano()-start)
	next := cnode
	for next != nil {
		foundPath = append(foundPath, next)
		next = next.parent
	}

	sort.Sort(sort.Reverse(foundPath))

	return foundPath
}

func manhatanDistance(fx, fy, tx, ty int) int {
	dx := math.Abs(float64(tx - fx))
	dy := math.Abs(float64(ty - fy))
	return int(dx + dy)
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
		d := euclideanDistance(current.x, current.y, e[i].x, e[i].y)
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

func topObstacles(raw grid, x, y, ty int) bool {
	for y != ty {
		if !canWalk(raw, x, y) {
			return true
		}
		y--
	}

	return false
}

func bottomObstacles(wv grid, x, y, ty int) bool {
	for y != ty {
		if !canWalk(wv, x, y) {
			return true
		}
		y++
	}
	return false
}

func rightObstacles(wv grid, x, y, tx int) bool {
	for x != tx {
		if !canWalk(wv, x, y) {
			return true
		}
		x++
	}
	return false
}

func leftObstacles(wv grid, x, y, tx int) bool {
	for x != tx {
		if !canWalk(wv, x, y) {
			return true
		}
		x--
	}
	return false
}

func topRightDiagonalObstacles(wv grid, x, y, tx, ty int) bool {
	for x != tx && y != ty {
		if !canWalk(wv, x, y) {
			return true
		}
		x++
		y--
	}

	return false
}


func topLeftDiagonalObstacles(wv grid, x, y, tx, ty int) bool {
	for x != tx && y != ty {
		if !canWalk(wv, x, y) {
			return true
		}
		x--
		y--
	}

	return false
}


func bottomRightDiagonalObstacles(wv grid, x, y, tx, ty int) bool {
	for x != tx && y != ty {
		if !canWalk(wv, x, y) {
			return true
		}
		x++
		y++
	}

	return false
}

func bottomLeftDiagonalObstacles(wv grid, x, y, tx, ty int) bool {
	for x != tx && y != ty {
		if !canWalk(wv, x, y) {
			return true
		}
		x--
		y++
	}

	return false
}

// using node grid and raw grid, get neighboring nodes
func adjacentNodes(ng, rg grid, x, y , margin int) nodes {
	var (
		result nodes
		n * node
	)

	// between x,y and adjacent node
	// assert that between x,y and adjacent node there are only walkable nodes

	// ↑
	if !topObstacles(rg, x,y, y-margin) {
		n = ng.get(x, y-margin)
		if n != nil {
			result = append(result, n)
		}
	}

	// ↓
	if !bottomObstacles(rg, x, y, y+margin) {
		n = ng.get(x, y+margin)
		if n != nil {
			result = append(result, n)
		}
	}

	// →
	if !rightObstacles(rg, x, y, x+margin) {
		n = ng.get(x+margin, y)
		if n != nil {
			result = append(result, n)
		}
	}

	// left
	if !leftObstacles(rg, x,y, x-margin) {
		n = ng.get(x-margin, y)
		if n != nil {
			result = append(result, n)
		}
	}

	// ↖
	if !topLeftDiagonalObstacles(rg, x, y, x-margin, y-margin) {
		n = ng.get(x-margin, y-margin)
		if n != nil {
			result = append(result, n)
		}
	}

	// ↗
	if !topRightDiagonalObstacles(rg, x,y, x+margin, y-margin) {
		n = ng.get(x+margin, y-margin)
		if n != nil {
			result = append(result, n)
		}
	}

	// ↘
	if !bottomRightDiagonalObstacles(rg, x, y, x+margin, y+margin) {
		n = ng.get(x+margin, y+margin)
		if n != nil {
			result = append(result, n)
		}
	}

	// ↙
	if !bottomLeftDiagonalObstacles(rg, x,y, x-margin, y+margin) {
		n = ng.get(x-margin, y+margin)
		if n != nil {
			result = append(result, n)
		}
	}

	return result
}

func in(list nodes, sv *node) bool {
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

// really inefficient
// I should check how to avoid making use of the grid pointers
func copyGrid(g grid) grid {
	var ng = make(grid)
	for k1, v1 := range g {
		for k2, v2 := range v1 {
			n := *v2

			_, ok := ng[k1]
			if !ok {
				ng[k1] = make(map[int]*node)
			}
			ng[k1][k2] = &n
		}
	}
	return ng
}
