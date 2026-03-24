import taichi as ti

@ti.kernel
def image_to_mask(
    image: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    mask: ti.types.ndarray(dtype=ti.f32, ndim=2),
):
    """
    将 4 通道图像转换为单通道遮罩（使用 alpha 通道）

    :param image: 输入图像 [w, h, [r, g, b, a]]
    :param mask: 输出遮罩 [w, h]
    """
    w, h = mask.shape
    for i, j in ti.ndrange(w, h):
        # 取 alpha 通道作为遮罩值
        mask[i, j] = image[i, j][3]



@ti.kernel
def feather_linear(dist: ti.template(), output: ti.template(), feather_radius: ti.f32):
    """
    应用线性羽化效果

    :param dist: 距离场
    :param output: 输出遮罩（0-1）
    :param feather_radius: 羽化半径
    """
    w, h = dist.shape
    for i, j in ti.ndrange(w, h):

        alpha = ti.min(dist[i, j] / feather_radius, 1.0)
        output[i, j] = ti.clamp(alpha, 0.0, 1.0)


@ti.kernel
def feather_conic(dist: ti.template(), output: ti.template(), feather_radius: ti.f32):
    """
    应用圆锥羽化效果

    :param dist: 距离场
    :param output: 输出遮罩（0-1）
    :param feather_radius: 羽化半径
    """
    w, h = dist.shape
    for i, j in ti.ndrange(w, h):
        norm_dist = ti.min(dist[i, j] / feather_radius, 1.0)

        alpha = ti.pow(norm_dist, 1.6)
        output[i, j] = ti.clamp(alpha, 0.0, 1.0)


@ti.kernel
def feather_smoothstep(dist: ti.template(), output: ti.template(), feather_radius: ti.f32):
    """
    应用平滑步函数羽化效果

    :param dist: 距离场
    :param output: 输出遮罩（0-1）
    :param feather_radius: 羽化半径
    """
    w, h = dist.shape
    for i, j in ti.ndrange(w, h):
        norm_dist = ti.min(dist[i, j] / feather_radius, 1.0)

        alpha = norm_dist * norm_dist * (3.0 - 2.0 * norm_dist)
        output[i, j] = ti.clamp(alpha, 0.0, 1.0)


@ti.kernel
def feather_sigmoid(dist: ti.template(), output: ti.template(), feather_radius: ti.f32):
    """
    应用 sigmoid 函数羽化效果

    :param dist: 距离场
    :param output: 输出遮罩（0-1）
    :param feather_radius: 羽化半径
    """
    w, h = dist.shape
    for i, j in ti.ndrange(w, h):
        norm_dist = ti.min(dist[i, j] / feather_radius, 1.0)

        k = 6.0
        alpha = 1.0 / (1.0 + ti.exp(-k * (norm_dist - 0.5)))
        output[i, j] = ti.clamp(alpha, 0.0, 1.0)