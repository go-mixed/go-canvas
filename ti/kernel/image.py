import taichi as ti

@ti.func
def cv_color_to_ti(x: ti.types.vector(3, ti.u32)) -> ti.types.vector(4, ti.f32):
    """
    将 cv color (b, g, r) 转换为 (r, g, b, a)
    输入：vec3<u32> (0-255 范围的整数)
    输出：vec4<f32> (0.0-1.0 范围的浮点数)
    """
    return ti.math.vec4(
        ti.cast(x[2], ti.f32) / 255.0,
        ti.cast(x[1], ti.f32) / 255.0,
        ti.cast(x[0], ti.f32) / 255.0,
        1.0
    )

@ti.kernel
def cv_image_to_ti(
    input: ti.types.ndarray(element_shape=(3,), dtype=ti.u32, ndim=2),
    output: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
):
    """
    将 OpenCV 的 [y, x] = uint32([b, g, r]) 转换为 Taichi 的 [x, y] = float32([r, g, b, a])
    同时完成坐标转置和颜色通道重排

    注意：使用 u32 而非 u8 以兼容 Vulkan/SPIR-V 后端
    Go 端需要将 uint8 数据转换为 uint32 后传入
    """
    for i, j in output:
        output[i, j] = cv_color_to_ti(input[j, i])


@ti.kernel
def fill_texture(
    texture: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    r: ti.f32, g: ti.f32, b: ti.f32, a: ti.f32
):
    """
    填充纹理
    """
    v = ti.math.vec4(r, g, b, a)
    for i, j in texture:
        texture[i, j] = v

