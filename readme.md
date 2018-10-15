
# imagematch
A very simple nearest-neighbour image matching algorithm.

Still incomplete and broken, do not use this yet.

# Build
This is a Go program so you will need Go to build it.

```sh
git clone http://github.com/AjayMT/imagematch.git
cd imagematch
go build
```

# Usage
```
imagematch <training-data> <test-image> <tolerance>
```

For example, to match `testdata/test_a.png` with an image from `traindata/fixed-width/` (PNGs of a fixed width font) with a tolerance of 0.5:
```
imagematch traindata/fixed-width testdata/test_a.png 0.5
```

I will eventually get around to documenting this better.
