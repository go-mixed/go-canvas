package ti

// RenderLayerOptions 渲染层选项
type RenderLayerOptions struct {
	X, Y       float32 // 相对屏幕的偏移
	Cx, Cy     float32 // 纹理中心坐标
	Scale      float32 // 缩放倍数
	Rotation   float32 // 旋转弧度
	Alpha      float32 // 透明度 0.0-1.0
	Width      float32 // 纹理宽度
	Height     float32 // 纹理高度
	MinX, MaxX int32   // 包围盒 x 范围
	MinY, MaxY int32   // 包围盒 y 范围
}

// RenderLayerNoMask 渲染层（无遮罩）
func (m *AotModule) RenderLayerNoMask(texture *TiImage, screen *TiImage, opts RenderLayerOptions) {
	kernel := m.getCache("render_layer_no_mask")
	kernel.Launch().
		ArgNdArray(texture).
		ArgFloat32(opts.X).ArgFloat32(opts.Y).
		ArgFloat32(opts.Cx).ArgFloat32(opts.Cy).
		ArgFloat32(opts.Scale).ArgFloat32(opts.Rotation).ArgFloat32(opts.Alpha).
		ArgFloat32(opts.Width).ArgFloat32(opts.Height).
		ArgInt32(opts.MinX).ArgInt32(opts.MaxX).ArgInt32(opts.MinY).ArgInt32(opts.MaxY).
		ArgNdArray(screen).
		Run()
}

// RenderLayerWithMask 渲染层（带遮罩）
func (m *AotModule) RenderLayerWithMask(texture, mask, screen *TiImage, opts RenderLayerOptions) {
	kernel := m.getCache("render_layer_with_mask")
	kernel.Launch().
		ArgNdArray(texture).
		ArgFloat32(opts.X).ArgFloat32(opts.Y).
		ArgFloat32(opts.Cx).ArgFloat32(opts.Cy).
		ArgFloat32(opts.Scale).ArgFloat32(opts.Rotation).ArgFloat32(opts.Alpha).
		ArgFloat32(opts.Width).ArgFloat32(opts.Height).
		ArgInt32(opts.MinX).ArgInt32(opts.MaxX).ArgInt32(opts.MinY).ArgInt32(opts.MaxY).
		ArgNdArray(mask).
		ArgNdArray(screen).
		Run()
}
