
package main

import (
	"fmt"
	"image/png"
	"image"
	"os"
	"sync"
	"io/ioutil"
	"path/filepath"
	"math"
	"strconv"
)

func minSlice(slice []float64) (int, float64) {
	if len(slice) == 0 { panic("cannot find minimum of empty slice") }

	index := 0
	min := slice[index]

	for i, n := range slice {
		if n < min {
			min = n
			index = i
		}
	}

	return index, min
}

func avgSlice(slice []float64) float64 {
	if len(slice) == 0 { panic("cannot find average of empty slice") }

	sum := 0.0
	for _, n := range slice { sum += n }

	return sum / float64(len(slice))
}

func ComputeDistance(refmatrix [][]bool, testmatrix [][]bool, tolerance float64) float64 {
	refsize := 0
	testsize := 0

	for _, col := range refmatrix {
		for _, pix := range col {
			if pix { refsize++ }
		}
	}

	for _, col := range testmatrix {
		for _, pix := range col {
			if pix { testsize++ }
		}
	}

	if refsize > testsize {
		refsize, testsize = testsize, refsize
		refmatrix, testmatrix = testmatrix, refmatrix
	}

	var pixdistances []float64

	for xindex, col := range testmatrix {
		for yindex, pix := range col {
			if ! pix { continue }

			minx := xindex - int(float64(len(refmatrix)) * tolerance)
			maxx := xindex + int(float64(len(refmatrix)) * tolerance)
			miny := yindex - int(float64(len(refmatrix[0])) * tolerance)
			maxy := yindex + int(float64(len(refmatrix[0])) * tolerance)

			if minx < 0 { minx = 0 }
			if maxx > len(refmatrix) { maxx = len(refmatrix) }
			if miny < 0 { miny = 0 }
			if maxy > len(refmatrix[0]) { maxy = len(refmatrix[0]) }

			var distances []float64

			for refx := minx; refx < maxx; refx++ {
				for refy := miny; refy < maxy; refy++ {
					if ! refmatrix[refx][refy] { continue }

					xdiff := float64(refx - xindex)
					ydiff := float64(refy - yindex)
					dist := math.Sqrt(math.Pow(xdiff, 2) + math.Pow(ydiff, 2))

					distances = append(distances, dist)
				}
			}

			if len(distances) == 0 {
				// no pixel within tolerance range
				// TODO real solution
				faraway := math.Sqrt(math.Pow(float64(len(refmatrix)), 2) +
					math.Pow(float64(len(refmatrix[0])), 2))
				pixdistances = append(pixdistances, faraway)
			} else {
				_, mindist := minSlice(distances)
				pixdistances = append(pixdistances, mindist)
			}
		}
	}

	return avgSlice(pixdistances)// + (float64(testsize) / float64(refsize))
}

func ScaleMatrix(mat [][]bool, width int, height int) [][]bool {
	scaled := make([][]bool, width)

	for x := 0; x < len(scaled); x++ {
		scaled[x] = make([]bool, height)

		for y := 0; y < len(scaled[x]); y++ {
			matx := int((float64(x) / float64(width)) * float64(len(mat)))
			maty := int(((float64(y)) / float64(height)) * float64(len(mat[0])))

			scaled[x][y] = mat[matx][maty]
		}
	}

	return scaled
}

func TrimMatrix(mat [][]bool) [][]bool {
	minx := 0
	miny := 0
	maxx := len(mat)
	maxy := len(mat[0])

	for xindex, col := range mat {
		hastrue := false

		for _, pix := range col {
			if pix {
				hastrue = true
				break
			}
		}

		if hastrue {
			minx = xindex
			break
		}
	}

	for i := (len(mat) - 1); i >= 0; i-- {
		hastrue := false

		for _, pix := range mat[i] {
			if pix {
				hastrue = true
				break
			}
		}

		if hastrue {
			maxx = i
			break
		}
	}

	for y := 0; y < len(mat[0]); y++ {
		hastrue := false

		for x := 0; x < len(mat); x++ {
			if mat[x][y] {
				hastrue = true
				break
			}
		}

		if hastrue {
			miny = y
			break
		}
	}

	for y := (len(mat[0]) - 1); y >= 0; y-- {
		hastrue := false

		for x := 0; x < len(mat); x++ {
			if mat[x][y] {
				hastrue = true
				break
			}
		}

		if hastrue {
			maxy = y
			break
		}
	}

	trimmed := make([][]bool, (maxx - minx) + 1)

	for i := 0; i < len(trimmed); i++ {
		trimmed[i] = mat[minx + i][miny:(maxy + 1)]
	}

	return trimmed
}

func ToMatrix(img image.Image) [][]bool {
	minx := img.Bounds().Min.X
	miny := img.Bounds().Min.Y
	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y

	matrix := make([][]bool, (maxx - minx) + 1)

	for x := 0; x < len(matrix); x++ {
		matrix[x] = make([]bool, (maxy - miny) + 1)

		for y := 0; y < len(matrix[x]); y++ {
			r, g, b, a := img.At(minx + x, miny + y).RGBA()
			isblack := r == 0 && g == 0 && b == 0 && a != 0
			matrix[x][y] = isblack
		}
	}

	return matrix
}

func main() {
	datapath := os.Args[1]
	testpath := os.Args[2]
	tolerance, _ := strconv.ParseFloat(os.Args[3], 64)

	files, _ := ioutil.ReadDir(datapath)

	fmt.Printf("Processing %s... ", testpath)

	testfile, _ := os.Open(testpath)
	testimage, _ := png.Decode(testfile)
	testmatrix := TrimMatrix(ToMatrix(testimage))

	fmt.Printf("done.\n\n")

	keys := make([]string, len(files))
	distances := make([]float64, len(files))
	var wg sync.WaitGroup

	for i, f := range files {
		wg.Add(1)

		go func(index int, file os.FileInfo) {
			defer wg.Done()

			fmt.Printf("Computing distance to %s...\n", file.Name())

			dtfile, _ := os.Open(filepath.Join(datapath, file.Name()))
			dtimage, _ := png.Decode(dtfile)
			matrix := TrimMatrix(ToMatrix(dtimage))
			width := len(matrix)
			height := len(matrix[0])

			distance := ComputeDistance(matrix, ScaleMatrix(testmatrix, width, height), tolerance)
			keys[index] = file.Name()
			distances[index] = distance

			fmt.Printf("Distance to %s: %f\n", file.Name(), distance)
		}(i, f)
	}

	wg.Wait()

	minindex, _ := minSlice(distances)
	fmt.Printf("\nClosest match: %s\n", keys[minindex])
}
