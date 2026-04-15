package ctypes

import "github.com/go-mixed/go-taichi/taichi"

type TiImage = taichi.NdArray   // ndarray[w, h, [r, g, b, a]]
type BgraImage = taichi.NdArray // ndarray[w, h, u32] u32 == [0xbbggrraa]
type TiMask = taichi.NdArray    // ndarray[w, h, a]
type TiGrid = taichi.NdArray    // ndarray[w, h, f32]
type TiColor = []float32        // [4]
