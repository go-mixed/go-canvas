import taichi as ti

@ti.kernel
def fill_color(
    texture: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    color: ti.types.vector(4, ti.f32),
):
    """
    填充纹理
    """
    for i, j in texture:
        texture[i, j] = color


# ============================================================================
# Blur 模糊效果 kernels
# ============================================================================

@ti.kernel
def blur_box(
    input: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    output: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    kernel_size: ti.i32,
):
    """
    普通模糊（Box Blur）
    每个像素取周围 kernel_size x kernel_size 区域的平均值

    Args:
        input: 输入纹理
        output: 输出纹理
        kernel_size: 模糊半径（3 = 3x3 区域，5 = 5x5 区域，以此类推）
    """
    width = input.shape[0]
    height = input.shape[1]
    k = kernel_size // 2

    for x, y in output:
        r, g, b, a = 0.0, 0.0, 0.0, 0.0
        count = 0.0

        for dy in range(-k, k + 1):
            for dx in range(-k, k + 1):
                px = x + dx
                py = y + dy
                if 0 <= px < width and 0 <= py < height:
                    c = input[px, py]
                    r += c[0]
                    g += c[1]
                    b += c[2]
                    a += c[3]
                    count += 1.0

        if count > 0:
            output[x, y] = ti.math.vec4(r / count, g / count, b / count, a / count)


@ti.kernel
def blur_gaussian(
    input: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    output: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    kernel_size: ti.i32,
    # sigma: ti.f32,
):
    """
    高斯模糊（Gaussian Blur）
    使用高斯分布权重进行模糊，效果更自然柔和

    Args:
        input: 输入纹理
        output: 输出纹理
        kernel_size: 模糊半径（3 = 3x3 区域，5 = 5x5 区域，以此类推）
        sigma: 高斯分布的标准差（通常等于 kernel_size / 3）
    """
    width = input.shape[0]
    height = input.shape[1]
    k = kernel_size // 2
    sigma = kernel_size / 3

    for x, y in output:
        r, g, b, a = 0.0, 0.0, 0.0, 0.0
        weight_sum = 0.0

        for dy in range(-k, k + 1):
            for dx in range(-k, k + 1):
                px = x + dx
                py = y + dy
                if 0 <= px < width and 0 <= py < height:
                    # 高斯权重
                    dist_sq = ti.cast(dx * dx + dy * dy, ti.f32)
                    weight = ti.exp(-dist_sq / (2.0 * sigma * sigma))
                    c = input[px, py]
                    r += c[0] * weight
                    g += c[1] * weight
                    b += c[2] * weight
                    a += c[3] * weight
                    weight_sum += weight

        if weight_sum > 0.0:
            output[x, y] = ti.math.vec4(r / weight_sum, g / weight_sum, b / weight_sum, a / weight_sum)


@ti.kernel
def blur_mosaic(
    input: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    output: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    block_size: ti.i32,
):
    """
    马赛克效果（Mosaic / Pixelation）
    将图像分割成 block_size x block_size 的块，每个块使用中心像素的颜色

    Args:
        input: 输入纹理
        output: 输出纹理
        block_size: 马赛克块大小（3 = 3x3 块，5 = 5x5 块，以此类推）
    """
    width = input.shape[0]
    height = input.shape[1]
    k = block_size // 2

    for x, y in output:
        # 找到当前块中心点
        center_x = ((x // block_size) * block_size) + k
        center_y = ((y // block_size) * block_size) + k

        # 边界检查
        if center_x >= width:
            center_x = width - 1
        if center_y >= height:
            center_y = height - 1

        output[x, y] = input[center_x, center_y]

