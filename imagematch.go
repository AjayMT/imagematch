
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

func ComputeDistance(matA [][]float64, matB [][]float64) float64 {
	var differences []float64
	integralA := make([][]float64, len(matA))
	integralB := make([][]float64, len(matB))

	for x, col := range matA {
		integralA[x] = make([]float64, len(col))
		integralB[x] = make([]float64, len(col))

		for y, _ := range col {
			prevxA, prevyA, prevxyA := 0.0, 0.0, 0.0
			prevxB, prevyB, prevxyB := 0.0, 0.0, 0.0

			if x > 0 {
				prevxA = integralA[x - 1][y]
				prevxB = integralB[x - 1][y]
			}

			if y > 0 {
				prevyA = integralA[x][y - 1]
				prevyB = integralB[x][y - 1]
			}

			if x > 0 && y > 0 {
				prevxyA = integralA[x - 1][y - 1]
				prevxyB = integralB[x - 1][y - 1]
			}

			integralA[x][y] = matA[x][y] + prevxA + prevyA - prevxyA
			integralB[x][y] = matB[x][y] + prevxB + prevyB - prevxyB

			differences = append(differences,
				math.Abs(integralA[x][y] - integralB[x][y]))
		}
	}

	return avgSlice(differences)
}

func ScaleMatrix(mat [][]float64, width int, height int) [][]float64 {
	scaled := make([][]float64, width)

	for x := 0; x < len(scaled); x++ {
		scaled[x] = make([]float64, height)

		for y := 0; y < len(scaled[x]); y++ {
			matx := int((float64(x) / float64(width)) * float64(len(mat)))
			maty := int(((float64(y)) / float64(height)) * float64(len(mat[0])))

			scaled[x][y] = mat[matx][maty]
		}
	}

	return scaled
}

func TrimMatrix(mat [][]float64) [][]float64 {
	minx := 0
	miny := 0
	maxx := len(mat)
	maxy := len(mat[0])

	for xindex, col := range mat {
		hascolor := false

		for _, pix := range col {
			if pix > 0 {
				hascolor = true
				break
			}
		}

		if hascolor {
			minx = xindex
			break
		}
	}

	for i := (len(mat) - 1); i >= 0; i-- {
		hascolor := false

		for _, pix := range mat[i] {
			if pix > 0 {
				hascolor = true
				break
			}
		}

		if hascolor {
			maxx = i
			break
		}
	}

	for y := 0; y < len(mat[0]); y++ {
		hascolor := false

		for x := 0; x < len(mat); x++ {
			if mat[x][y] > 0 {
				hascolor = true
				break
			}
		}

		if hascolor {
			miny = y
			break
		}
	}

	for y := (len(mat[0]) - 1); y >= 0; y-- {
		hascolor := false

		for x := 0; x < len(mat); x++ {
			if mat[x][y] > 0 {
				hascolor = true
				break
			}
		}

		if hascolor {
			maxy = y
			break
		}
	}

	trimmed := make([][]float64, (maxx - minx) + 1)

	for i := 0; i < len(trimmed); i++ {
		trimmed[i] = mat[minx + i][miny:(maxy + 1)]
	}

	return trimmed
}

func ToMatrix(img image.Image) [][]float64 {
	minx := img.Bounds().Min.X
	miny := img.Bounds().Min.Y
	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y

	matrix := make([][]float64, (maxx - minx) + 1)

	for x := 0; x < len(matrix); x++ {
		matrix[x] = make([]float64, (maxy - miny) + 1)

		for y := 0; y < len(matrix[x]); y++ {
			r, g, b, a := img.At(minx + x, miny + y).RGBA()
			isblack := r == 0 && g == 0 && b == 0 && a != 0

			if isblack {
				matrix[x][y] = 1.0
			} else {
				matrix[x][y] = 0.0
			}
		}
	}

	return matrix
}

func main() {
	datapath := os.Args[1]
	testpath := os.Args[2]

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

			dtfile, _ := os.Open(filepath.Join(datapath, file.Name()))
			dtimage, _ := png.Decode(dtfile)
			matrix := TrimMatrix(ToMatrix(dtimage))
			width := len(matrix)
			height := len(matrix[0])

			distance := ComputeDistance(matrix,
				ScaleMatrix(testmatrix, width, height))
			keys[index] = file.Name()
			distances[index] = distance

			fmt.Printf("Distance to %s: %f\n", file.Name(), distance)
		}(i, f)
	}

	wg.Wait()

	minindex, mindist := minSlice(distances)
	fmt.Printf("\nClosest match: %s %f\n", keys[minindex], mindist)
}
