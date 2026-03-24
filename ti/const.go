package ti

import "github.com/go-mixed/go-taichi/taichi"

type CvImage = taichi.NdArray // ndarray[h, w, [b, g, r]]
type TiImage = taichi.NdArray // ndarray[w, h, [r, g, b, a]]
type TiMask = taichi.NdArray  // ndarray[w, h, a]
type TiGrid = taichi.NdArray  // ndarray[w, h, f32]
type TiColor = []float32      // [4]
