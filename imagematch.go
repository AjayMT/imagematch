
package main

import (
	"fmt"
	"image/png"
	"image"
	"os"
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

func computeDistance(matA [][]bool, matB [][]bool, tolerance float64) float64 {
	refmatrix := matA
	testmatrix := matB
	sizeA := 0
	sizeB := 0

	for _, col := range matA {
		for _, pix := range col {
			if pix { sizeA++ }
		}
	}

	for _, col := range matB {
		for _, pix := range col {
			if pix { sizeB++ }
		}
	}

	if sizeA > sizeB {
		refmatrix = matB
		testmatrix = matA
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
			} else {
				_, mindist := minSlice(distances)
				pixdistances = append(pixdistances, mindist)
			}
		}
	}

	return avgSlice(pixdistances)
}

func scaleMatrix(mat [][]bool, width int, height int) [][]bool {
	squashed := make([][]bool, width)

	for x := 0; x < len(squashed); x++ {
		squashed[x] = make([]bool, height)

		for y := 0; y < len(squashed[x]); y++ {
			matx := int((float64(x) / float64(width)) * float64(len(mat)))
			maty := int(((float64(y)) / float64(height)) * float64(len(mat[0])))

			squashed[x][y] = mat[matx][maty]
		}
	}

	return squashed
}

func trimMatrix(mat [][]bool) [][]bool {
	minx := 0
	miny := 0
	maxx := len(mat)
	maxy := len(mat[0])

	for xindex, col := range mat {
		hastrue := false

		for _, pix := range col {
			if pix {
				hastrue = true
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

func toMatrix(img image.Image) [][]bool {
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
	testmatrix := trimMatrix(toMatrix(testimage))

	fmt.Printf("done.\n\n")

	var keys []string
	var distances []float64

	for _, file := range files {
		fmt.Printf("Computing distance to %s... ", file.Name())
		dtfile, _ := os.Open(filepath.Join(datapath, file.Name()))
		dtimage, _ := png.Decode(dtfile)
		matrix := trimMatrix(toMatrix(dtimage))

		width := len(matrix)
		height := len(matrix[0])
		distance := computeDistance(matrix, scaleMatrix(testmatrix, width, height), tolerance)
		keys = append(keys, file.Name())
		distances = append(distances, distance)
		fmt.Printf("%f\n", distance)
	}

	minindex, _ := minSlice(distances)
	fmt.Printf("\nClosest match: %s\n", keys[minindex])
}
