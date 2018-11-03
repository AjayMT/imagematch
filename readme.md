
# imagematch
A very simple k-nearest-neighbours algorithm that classifies images of characters. Still in development.

## Build
This is a Go program so you will need Go to build it.

```sh
git clone http://github.com/AjayMT/imagematch.git
cd imagematch
go build
```

## Usage
```
imagematch <training-data> <test-image> <k>
```

For example, to find the 3 nearest neighbours of `testdata/test_g.png` from the dataset in `traindata`:
```
imagematch traindata testdata/test_g.png 3
```

TODO figure out how to calculate optimal `k` value

## Design rationale
This program does not process colors or binarize input images -- it ignores all pixels that are not completely black. It is also not very accurate without a large dataset of PNGs, which I haven't found or compiled yet. My primary goals when writing this program were:
- learn Go
- write a simple classification algorithm

which is why it isn't very sophisticated.

This program is also entirely stateless, i.e it can classify any test image based on any dataset. As such, it cannot be 'trained' in the traditional sense -- 'reference' dataset is probably a more accurate term than 'training' dataset. This is because it is not based on a neural network or any form of gradient descent. I wanted to find a simpler, mathematically deterministic way to compare two arbitrary matrices, which is what this program does.

## Distance function
PNG images are first translated into matrices -- as of now, black pixels are `1.0` and all other colors are `0.0`, but it is possible to use this same distance function with a range of different values. Before distances are computed, the matrices are cropped to fit the black pixels and scaled to the same dimensions.

The distance function first computes the [integral images](https://en.wikipedia.org/wiki/Summed-area_table) of the two matrices, and then converts each value in the integral image to a fraction of the total integral. This effectively expresses the image as a probability density function `P(x, y)`, where

`P(x, y)` = probability that rectangle from `(0, 0)` to `(x, y)` contains a black pixel

The distance between the two matrices is the mean of the differences between the values of `P(x, y)` at every point in the two matrices.
