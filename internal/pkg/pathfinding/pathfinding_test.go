package path

import (
	"bytes"
	"fmt"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"image"
	"image/color"
	"math"
	"os"
	"sync"
	"testing"
)

var (
	mn                    string
	s                     *data.SHBD
	v                     Grid
	vNodesNoWallsMargin   = PresetNodesMargin
	vNodesWithWallsMargin = PresetNodesMargin
	vNodesWithWalls       Grid
	vNodesNoWalls         Grid
)

func TestMain(m *testing.M) {
	mn = "Rou"
	s, _ = data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", mn))
	//s, _ = data.LoadSHBDFile(fmt.Sprintf("/home/marius/projects/shine/shine.engine.emulator/files/blocks/%v.shbd", mn))
	v = RawGrid(s)
	vNodesNoWalls = PresetNodesGrid(s, vNodesNoWallsMargin)
	vNodesWithWalls = PresetNodesWithMargins(s, v, vNodesWithWallsMargin)
	os.Exit(m.Run())
}

func Test_Path_A_astar_paint_large_distance(t *testing.T) {
	v1 := copyGrid(vNodesNoWalls)

	img1, err := data.SHBDToImage(s)

	if err != nil {
		logger.Error(err)
	}

	pathVertices1 := astar(v1, v, 1566, 500, 466, 1566, vNodesNoWallsMargin, "octile")
	fmt.Printf("vNodesNoWallsMargin path Nodes %v\n", len(pathVertices1))

	for _, pv := range pathVertices1 {
		img1.Set(pv.X, pv.Y, color.RGBA{
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
	// print an image with the Node data

	v2 := copyGrid(vNodesWithWalls)

	img2, err := data.SHBDToImage(s)

	if err != nil {
		logger.Error(err)
	}

	pathVertices2 := astar(v2, v, 1566, 500, 466, 1566, vNodesWithWallsMargin, "octile")
	fmt.Printf("vNodesWithWalls path Nodes %v\n", len(pathVertices2))

	for _, pv := range pathVertices2 {
		img2.Set(pv.X, pv.Y, color.RGBA{
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

	pathVertices3 := astar(v3, v, 1566, 500, 466, 1566, 1, "octile")

	fmt.Printf("raw path Nodes %v\n", len(pathVertices3))

	for _, pv := range pathVertices3 {
		img3.Set(pv.X, pv.Y, color.RGBA{
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

func Test_Path_A_astar_paint_short_distance(t *testing.T) {
	v1 := copyGrid(vNodesNoWalls)

	img1, err := data.SHBDToImage(s)

	if err != nil {
		logger.Error(err)
	}

	pathVertices1 := astar(v1, v, 900, 862, 908, 821, vNodesNoWallsMargin, "octile")
	fmt.Printf("vNodesNoWallsMargin path Nodes %v\n", len(pathVertices1))

	for _, pv := range pathVertices1 {
		img1.Set(pv.X, pv.Y, color.RGBA{
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
	// print an image with the Node data

	v2 := copyGrid(vNodesWithWalls)

	img2, err := data.SHBDToImage(s)

	if err != nil {
		logger.Error(err)
	}

	pathVertices2 := astar(v2, v, 900, 862, 908, 821, vNodesWithWallsMargin, "octile")
	fmt.Printf("vNodesWithWalls path Nodes %v\n", len(pathVertices2))

	for _, pv := range pathVertices2 {
		img2.Set(pv.X, pv.Y, color.RGBA{
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

	pathVertices3 := astar(v3, v, 900, 862, 908, 821, 1, "octile")

	fmt.Printf("raw path Nodes %v\n", len(pathVertices3))

	for _, pv := range pathVertices3 {
		img3.Set(pv.X, pv.Y, color.RGBA{
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

	v := RawGrid(s)

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

	// print an image with the Node data

	// paint Nodes that will be used for paths.
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

			pathVertices := astar(cgrid, v, 835, 700, 1070, 1540, 1, "octile")

			for _, pv := range pathVertices {
				img.Set(pv.X, pv.Y, color.RGBA{
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

	// print an image with the Node data

	// paint Nodes that will be used for paths.
}

func Test_Map_Intermitent_By_Speed_Path_A_B_AStar(t *testing.T) {
	// raw path has too many Nodes
	// given a speed, reduce the raw path using the speed
	// have a function that calculates how many Nodes should be sent per second

	// speed is a constant, the entity may be walking / running, e.G for entities 120, 60

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

	pn := PresetNodesGrid(s, 4)

	if err != nil {
		logger.Fatal(err)
	}

	pathVertices := astar(pn, v, 835, 700, 1070, 1540, 1, "octile")

	pathVertices = reduce(pathVertices, 15)

	for _, pv := range pathVertices {
		img.Set(pv.X, pv.Y, color.RGBA{
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
	// astar
	b.Run("astar_large_distance_raw", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ng := copyGrid(v)
			astar(ng, v, 1566, 500, 466, 1566, 1, "octile")
		}
	})

	b.Run("astar_large_distance_preset_nodes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ng := copyGrid(vNodesNoWalls)
			astar(ng, v, 1566, 500, 466, 1566, vNodesNoWallsMargin, "octile")
		}
	})

	b.Run("astar_large_distancepreset_nodes_with_wall_margin", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ng := copyGrid(vNodesWithWalls)
			astar(ng, v, 1566, 500, 466, 1566, vNodesWithWallsMargin, "octile")
		}
	})

	b.Run("astar_short_distance_raw", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ng := copyGrid(v)
			astar(ng, v, 900, 862, 908, 821, 1, "octile")
		}
	})

	b.Run("astar_short_distance_preset_nodes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ng := copyGrid(vNodesNoWalls)
			astar(ng, v, 900, 862, 908, 821, vNodesNoWallsMargin, "octile")
		}
	})

	b.Run("astar_short_distance_preset_nodes_with_wall_margin", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ng := copyGrid(vNodesWithWalls)
			astar(ng, v, 900, 862, 908, 821, vNodesWithWallsMargin, "octile")
		}
	})

	// jps
	// trace
	// bfs
	// dijkstra

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

	b.Run("RawGrid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v = RawGrid(s)
		}
	})

	b.Run("preset_nodes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			PresetNodesGrid(s, 4)
		}
	})

	b.Run("preset_nodes_wall_margin", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			PresetNodesGrid(s, 4)
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

	pn := PresetNodesGrid(s, 4)

	if err != nil {
		logger.Fatal(err)
	}

	pathVertices := astar(pn, v, 835, 700, 1070, 1540, 1, "octile")
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

	v := RawGrid(s)

	pn := PresetNodesGrid(s, 4)

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

func PaintNodesAndWallMargins(s *data.SHBD, wv Grid) (*image.RGBA, error) {
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
					countn := len(adjacentNodes(wv, v, rX, rY, 2))
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
	fmt.Printf("count Nodes %v\n", count)
	return img, nil
}

func PaintNodesWithoutWallMargins(s *data.SHBD, wv Grid) (*image.RGBA, error) {
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
	fmt.Printf("count Nodes %v\n", count)
	return img, nil
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