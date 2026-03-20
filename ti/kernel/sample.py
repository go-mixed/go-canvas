import taichi as ti


@ti.func
def lanczos_weight(x: ti.f32, a: ti.i32) -> ti.f32:
    """
    Lanczos 窗口函数
    a: 窗口大小（通常为 2, 3, 或 4）
    """
    result = 0.0
    x_abs = ti.abs(x)
    if x_abs < 1e-6:
        result = 1.0
    elif x_abs < a:
        pi_x = 3.14159265359 * x_abs
        result = a * ti.sin(pi_x) * ti.sin(pi_x / a) / (pi_x * pi_x)
    return result


@ti.func
def lanczos4_sample(
    image: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    x: ti.f32,
    y: ti.f32,
    width: ti.f32,
    height: ti.f32,
) -> ti.types.vector(4, ti.f32):
    """
    Lanczos4 插值采样（质量最高，适合高质量放大）

    :param image: 图像字段
    :param x: 浮点x坐标
    :param y: 浮点y坐标
    :param width: 图像宽度
    :param height: 图像高度
    :return: 插值后的颜色
    """
    result = ti.math.vec4(0.0, 0.0, 0.0, 0.0)

    # 中心像素坐标
    x_center = ti.floor(x)
    y_center = ti.floor(y)

    # 采样 8x8 邻域（Lanczos4 需要 a=4）
    for dy in ti.static(range(-3, 5)):
        for dx in ti.static(range(-3, 5)):
            # 邻近像素坐标
            px = ti.cast(x_center + dx, ti.i32)
            py = ti.cast(y_center + dy, ti.i32)

            # 边界检查
            if 0 <= px < width and 0 <= py < height:
                # 计算权重
                wx = lanczos_weight(x - (x_center + dx), 4)
                wy = lanczos_weight(y - (y_center + dy), 4)
                weight = wx * wy

                # 累加加权颜色
                result += image[px, py] * weight

    return result


@ti.func
def cubic_weight(t: ti.f32) -> ti.f32:
    """
    三次插值权重函数（Catmull-Rom）
    """
    t_abs = ti.abs(t)
    result = 0.0
    if t_abs <= 1.0:
        result = 1.5 * t_abs * t_abs * t_abs - 2.5 * t_abs * t_abs + 1.0
    elif t_abs <= 2.0:
        result = -0.5 * t_abs * t_abs * t_abs + 2.5 * t_abs * t_abs - 4.0 * t_abs + 2.0
    return result


@ti.func
def bicubic_sample(
    image: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    x: ti.f32,
    y: ti.f32,
    width: ti.f32,
    height: ti.f32,
) -> ti.types.vector(4, ti.f32):
    """
    双三次插值采样（Bicubic Interpolation）
    质量远超双线性插值，适合图像放大

    :param image: 图像字段
    :param x: 浮点x坐标
    :param y: 浮点y坐标
    :param width: 图像宽度
    :param height: 图像高度
    :return: 插值后的颜色
    """
    result = ti.math.vec4(0.0, 0.0, 0.0, 0.0)

    # 中心像素坐标
    x_center = ti.floor(x)
    y_center = ti.floor(y)

    # 采样 4x4 邻域
    for dy in ti.static(range(-1, 3)):
        for dx in ti.static(range(-1, 3)):
            # 邻近像素坐标
            px = ti.cast(x_center + dx, ti.i32)
            py = ti.cast(y_center + dy, ti.i32)

            # 边界检查
            if 0 <= px < width and 0 <= py < height:
                # 计算权重
                wx = cubic_weight(x - (x_center + dx))
                wy = cubic_weight(y - (y_center + dy))
                weight = wx * wy

                # 累加加权颜色
                result += image[px, py] * weight

    return result


@ti.func
def bilinear_sample(
    image: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    x: ti.f32,
    y: ti.f32,
    width: ti.f32, 
    height: ti.f32,
) -> ti.types.vector(4, ti.f32):
    """
    双线性插值采样

    :param image: 图像字段
    :param x: 浮点x坐标
    :param y: 浮点y坐标
    :param width: 图像宽度
    :param height: 图像高度
    :return: 插值后的颜色
    """
    result = ti.math.vec4(0.0, 0.0, 0.0, 0.0)

    # 获取四个邻近像素的坐标
    # Bottom-left corner
    x1 = ti.cast(ti.floor(x), ti.i32)
    y1 = ti.cast(ti.floor(y), ti.i32)
    # Top-right corner
    x2 = ti.min(x1 + 1, ti.cast(width - 1, ti.i32))
    y2 = ti.min(y1 + 1, ti.cast(height - 1, ti.i32))

    # 边界检查
    if 0 <= x1 < width and 0 <= y1 < height:
        # 获取四个角点的颜色
        Q11 = image[x1, y1]
        Q21 = image[x2, y1]
        Q12 = image[x1, y2]
        Q22 = image[x2, y2]

        # 计算插值权重
        fx = x - x1
        fy = y - y1

        # 双线性插值
        R1 = Q11 * (1.0 - fx) + Q21 * fx
        R2 = Q12 * (1.0 - fx) + Q22 * fx
        result = R1 * (1.0 - fy) + R2 * fy

    return result