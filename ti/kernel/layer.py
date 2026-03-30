import taichi as ti

from ti.kernel.sample import *


@ti.func
def build_inverse_affine_matrix(
        x: ti.f32, y: ti.f32,
        cx: ti.f32, cy: ti.f32,
        scale_x: ti.f32, scale_y: ti.f32,
        rotation: ti.f32
) -> ti.types.matrix(3, 3, ti.f32):
    """
    构建逆仿射变换矩阵（屏幕坐标 → 纹理坐标）

    正变换顺序：平移 → 旋转 → 缩放
    逆变换顺序：逆缩放 → 逆旋转 → 逆平移

    参数说明：
        x, y: 纹理左边界在屏幕上的绝对坐标
        cx, cy: 纹理中心相对于纹理左边界(x=0)的偏移（比如 texture_width/2, texture_height/2）
                最终纹理中心在屏幕上的绝对位置 = (x + cx, y + cy)
    """
    cos_r = ti.cos(rotation)
    sin_r = ti.sin(rotation)
    inv_scale_x = 1.0 / scale_x
    inv_scale_y = 1.0 / scale_y

    # 逆仿射矩阵
    # [ cos/sx   sin/sx  -(x+cx)*cos/sx - (y+cy)*sin/sx + x + cx ]
    # [ -sin/sy  cos/sy   (x+cx)*sin/sy - (y+cy)*cos/sy + y + cy ]
    # [   0       0                    1                             ]

    tx = -(x + cx) * cos_r * inv_scale_x - (y + cy) * sin_r * inv_scale_x + x + cx
    ty = (x + cx) * sin_r * inv_scale_y - (y + cy) * cos_r * inv_scale_y + y + cy


    return ti.math.mat3(
        cos_r * inv_scale_x,  sin_r * inv_scale_x, tx,
       -sin_r * inv_scale_y, cos_r * inv_scale_y, ty,
        0.0, 0.0, 1.0
    )


@ti.func
def apply_affine_transform(matrix: ti.types.matrix(3, 3, ti.f32), x: ti.f32, y: ti.f32) -> ti.types.vector(2, ti.f32):
    """应用仿射变换（齐次坐标）"""
    v = ti.math.vec3(x, y, 1.0)
    result = matrix @ v
    return ti.math.vec2(result.x, result.y)


@ti.func
def sample_and_blend(
        texture: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        tex_x: ti.f32, tex_y: ti.f32,
        width: ti.f32, height: ti.f32,
        alpha: ti.f32,
        use_scale: ti.i32,
        screen_color: ti.types.vector(4, ti.f32),
        mask_value: ti.f32,
) -> ti.types.vector(4, ti.f32):
    """
    采样纹理并进行 Alpha 混合

    Args:
        mask_value: 遮罩值（无遮罩时传入 1.0）
    """
    # 采样纹理
    tex_color = screen_color
    if use_scale == 0:
        tex_color = texture[ti.cast(tex_x, ti.i32), ti.cast(tex_y, ti.i32)]
    else:
        # 双三次插值
        tex_color = bilinear_sample(texture, tex_x, tex_y, width, height)

    # 计算最终透明度
    final_alpha = ti.min(tex_color.w * alpha * mask_value, 1.0)

    # Alpha 混合
    new_color = tex_color
    if final_alpha >= 0.999:
        new_color.w = 1.0
    elif final_alpha > 1e-6:
        new_color = ti.math.mix(screen_color, tex_color, final_alpha)
        new_color.w = 1.0
    else:
        new_color = screen_color  # 完全透明，保持原色

    return new_color


@ti.kernel
def render_layer_no_mask(
        texture: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        x: ti.f32, y: ti.f32,
        cx: ti.f32, cy: ti.f32,
        scale_x: ti.f32, scale_y: ti.f32,
        rotation: ti.f32,
        alpha: ti.f32,
        width: ti.f32, height: ti.f32,
        min_x: ti.i32, max_x: ti.i32,
        min_y: ti.i32, max_y: ti.i32,
        screen: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
):
    """
    渲染层（无遮罩）

    :param texture: 纹理层，结构为[w, h] = [r, g, b, a]
    :param x: 前景相对screen的x偏移值，如果x为0，表示从左上角开始绘制
    :param y: 前景相对screen的y偏移值，如果y为0，表示从左上角开始绘制
    :param cx: 纹理中心相对于纹理左边界(x=0)的偏移（比如 texture_width/2）
    :param cy: 纹理中心相对于纹理上边界(y=0)的偏移（比如 texture_height/2）
    :param scale_x: x轴缩放，（单位：倍数）, 默认为1.0, 表示不缩放
    :param scale_y: y轴缩放，（单位：倍数）, 默认为1.0, 表示不缩放
    :param rotation: 旋转，（单位：弧度）, 默认为0.0, 表示不旋转
    :param alpha: 透明度，（单位：0.0-1.0）, 默认为1.0, 表示不透明
    :param width: 纹理宽度
    :param height: 纹理高度
    :param min_x: 纹理和屏幕的包围盒的左上角x坐标
    :param max_x: 纹理和屏幕的包围盒的右下角x坐标
    :param min_y: 纹理和屏幕的包围盒的左上角y坐标
    :param max_y: 纹理和屏幕的包围盒的右下角y坐标
    """
    use_scale = 1 if ti.abs(scale_x - 1.0) > 1e-6 or ti.abs(scale_y - 1.0) > 1e-6 else 0
    inv_matrix = build_inverse_affine_matrix(x, y, cx, cy, scale_x, scale_y, rotation)

    for x_screen, y_screen in ti.ndrange((min_x, max_x + 1), (min_y, max_y + 1)):
        # 屏幕坐标 → 纹理坐标
        tex_pos = apply_affine_transform(inv_matrix, ti.f32(x_screen), ti.f32(y_screen))

        # 边界检查
        if 0 <= tex_pos.x < width and 0 <= tex_pos.y < height:
            screen_color = screen[x_screen, y_screen]
            new_color = sample_and_blend(
                texture, tex_pos.x, tex_pos.y, width, height,
                alpha, use_scale, screen_color, 1.0  # mask_value = 1.0
            )
            if new_color.w > 0:  # 只有非透明时才写入
                screen[x_screen, y_screen] = new_color


@ti.kernel
def render_layer_with_mask(
        texture: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        x: ti.f32, y: ti.f32,
        cx: ti.f32, cy: ti.f32,
        scale_x: ti.f32, scale_y: ti.f32,
        rotation: ti.f32,
        alpha: ti.f32,
        width: ti.f32, height: ti.f32,
        min_x: ti.i32, max_x: ti.i32,
        min_y: ti.i32, max_y: ti.i32,
        mask: ti.types.ndarray(dtype=ti.f32, ndim=2),
        screen: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
):
    """
    渲染层（带遮罩）

    :param texture: 纹理层，结构为[w, h] = [r, g, b, a]
    :param x: 前景相对screen的x偏移值，如果x为0，表示从左上角开始绘制
    :param y: 前景相对screen的y偏移值，如果y为0，表示从左上角开始绘制
    :param cx: 纹理中心相对于纹理左边界(x=0)的偏移（比如 texture_width/2）
    :param cy: 纹理中心相对于纹理上边界(y=0)的偏移（比如 texture_height/2）
    :param scale_x: x轴缩放，（单位：倍数）, 默认为1.0, 表示不缩放
    :param scale_y: y轴缩放，（单位：倍数）, 默认为1.0, 表示不缩放
    :param rotation: 旋转，（单位：弧度）, 默认为0.0, 表示不旋转
    :param alpha: 透明度，（单位：0.0-1.0）, 默认为1.0, 表示不透明
    :param width: 纹理宽度
    :param height: 纹理高度
    :param min_x: 纹理和屏幕的包围盒的左上角x坐标
    :param max_x: 纹理和屏幕的包围盒的右下角x坐标
    :param min_y: 纹理和屏幕的包围盒的左上角y坐标
    :param max_y: 纹理和屏幕的包围盒的右下角y坐标
    :param mask: 遮罩，结构为[w, h] = alpha f32
    """
    use_scale = 1 if ti.abs(scale_x - 1.0) > 1e-6 or ti.abs(scale_y - 1.0) > 1e-6 else 0
    inv_matrix = build_inverse_affine_matrix(x, y, cx, cy, scale_x, scale_y, rotation)

    for x_screen, y_screen in ti.ndrange((min_x, max_x + 1), (min_y, max_y + 1)):
        # 屏幕坐标 → 纹理坐标
        tex_pos = apply_affine_transform(inv_matrix, ti.f32(x_screen), ti.f32(y_screen))

        # 边界检查
        if 0 <= tex_pos.x < width and 0 <= tex_pos.y < height:
            screen_color = screen[x_screen, y_screen]
            tex_x_i = ti.cast(tex_pos.x, ti.i32)
            tex_y_i = ti.cast(tex_pos.y, ti.i32)
            mask_value = mask[tex_x_i, tex_y_i]

            new_color = sample_and_blend(
                texture, tex_pos.x, tex_pos.y, width, height,
                alpha, use_scale, screen_color, mask_value
            )
            if new_color.w > 0:
                screen[x_screen, y_screen] = new_color

