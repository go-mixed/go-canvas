package ti

// ImageToMask 将图像转换为遮罩（提取 alpha 通道）
func (m *AotModule) ImageToMask(input *TiImage, out *TiMask) {
	kernel := m.getCache("image_to_mask")
	kernel.Launch().ArgNdArray(input).ArgNdArray(out).Run()
}

// ComputeDistanceField 计算距离场（欧几里得距离）
func (m *AotModule) ComputeDistanceField(mask *TiMask, dist *TiGrid) {
	kernel := m.getCache("compute_distance_field")
	kernel.Launch().ArgNdArray(mask).ArgNdArray(dist).Run()
}

// ComputeFeather 应用羽化效果
func (m *AotModule) ComputeFeather(dist *TiGrid, out *TiMask, featherRadius float32, featherMode FeatherMode) {
	featherKernelName := [...]string{"feather_linear", "feather_conic", "feather_smoothstep", "feather_sigmoid"}
	kernel := m.getCache(featherKernelName[featherMode])
	kernel.Launch().ArgNdArray(dist).ArgNdArray(out).ArgFloat32(featherRadius).Run()
}
