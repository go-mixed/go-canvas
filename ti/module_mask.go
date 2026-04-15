package ti

import "github.com/go-mixed/go-canvas/ctypes"

// AsyncImageToMask 将图像转换为遮罩（提取 alpha 通道）
func (m *AotModule) AsyncImageToMask(input *ctypes.TiImage, out *ctypes.TiMask) {
	kernel := m.getCache("image_to_mask")
	kernel.Launch().ArgNdArray(input).ArgNdArray(out).RunAsync()
}

func (m *AotModule) ImageToMask(input *ctypes.TiImage, out *ctypes.TiMask) {
	m.AsyncImageToMask(input, out)
	m.runtime.Wait()
}

// AsyncComputeDistanceField 计算距离场（欧几里得距离）
func (m *AotModule) AsyncComputeDistanceField(mask *ctypes.TiMask, dist *ctypes.TiGrid) {
	kernel := m.getCache("compute_distance_field")
	kernel.Launch().ArgNdArray(mask).ArgNdArray(dist).RunAsync()
}

func (m *AotModule) ComputeDistanceField(mask *ctypes.TiMask, dist *ctypes.TiGrid) {
	m.AsyncComputeDistanceField(mask, dist)
	m.runtime.Wait()
}

// AsyncComputeFeather 应用羽化效果
func (m *AotModule) AsyncComputeFeather(dist *ctypes.TiGrid, out *ctypes.TiMask, featherRadius float32, featherMode ctypes.FeatherMode) {
	featherKernelName := [...]string{"feather_linear", "feather_conic", "feather_smoothstep", "feather_sigmoid"}
	kernel := m.getCache(featherKernelName[featherMode])
	kernel.Launch().ArgNdArray(dist).ArgNdArray(out).ArgFloat32(featherRadius).RunAsync()
}

func (m *AotModule) ComputeFeather(dist *ctypes.TiGrid, out *ctypes.TiMask, featherRadius float32, featherMode ctypes.FeatherMode) {
	m.AsyncComputeFeather(dist, out, featherRadius, featherMode)
	m.runtime.Wait()
}
