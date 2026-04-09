package ti

// AsyncImageToMask 将图像转换为遮罩（提取 alpha 通道）
func (m *AotModule) AsyncImageToMask(input *TiImage, out *TiMask) {
	kernel := m.getCache("image_to_mask")
	kernel.Launch().ArgNdArray(input).ArgNdArray(out).RunAsync()
}

func (m *AotModule) ImageToMask(input *TiImage, out *TiMask) {
	m.AsyncImageToMask(input, out)
	m.runtime.Wait()
}

// AsyncComputeDistanceField 计算距离场（欧几里得距离）
func (m *AotModule) AsyncComputeDistanceField(mask *TiMask, dist *TiGrid) {
	kernel := m.getCache("compute_distance_field")
	kernel.Launch().ArgNdArray(mask).ArgNdArray(dist).RunAsync()
}

func (m *AotModule) ComputeDistanceField(mask *TiMask, dist *TiGrid) {
	m.AsyncComputeDistanceField(mask, dist)
	m.runtime.Wait()
}

// AsyncComputeFeather 应用羽化效果
func (m *AotModule) AsyncComputeFeather(dist *TiGrid, out *TiMask, featherRadius float32, featherMode FeatherMode) {
	featherKernelName := [...]string{"feather_linear", "feather_conic", "feather_smoothstep", "feather_sigmoid"}
	kernel := m.getCache(featherKernelName[featherMode])
	kernel.Launch().ArgNdArray(dist).ArgNdArray(out).ArgFloat32(featherRadius).RunAsync()
}

func (m *AotModule) ComputeFeather(dist *TiGrid, out *TiMask, featherRadius float32, featherMode FeatherMode) {
	m.AsyncComputeFeather(dist, out, featherRadius, featherMode)
	m.runtime.Wait()
}
