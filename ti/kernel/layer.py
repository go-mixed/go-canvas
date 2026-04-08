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
    # [ cos/sx   sin/sx  -(x+cx)*cos/sx - (y+cy)*sin/sx + cx ]
    # [ -sin/sy  cos/sy   (x+cx)*sin/sy - (y+cy)*cos/sy + cy ]
    # [   0       0                    1                             ]

    tx = -(x + cx) * cos_r * inv_scale_x - (y + cy) * sin_r * inv_scale_x + cx
    ty = (x + cx) * sin_r * inv_scale_y - (y + cy) * cos_r * inv_scale_y + cy


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
    tex_color = screen_color
    if use_scale == 0:
        tex_color = texture[ti.cast(tex_x, ti.i32), ti.cast(tex_y, ti.i32)]
    else:
        # 双三次插值
        tex_color = bilinear_sample(texture, tex_x, tex_y, width, height)

    # 计算最终源透明度（纹理alpha * 图层alpha * mask）
    src_a = ti.min(tex_color.w * alpha * mask_value, 1.0)
    out = screen_color

    if src_a > 1e-6:
        # 标准 Over（直通道颜色）
        dst_a = screen_color.w
        one_minus_src = 1.0 - src_a
        out_a = src_a + dst_a * one_minus_src

        if out_a > 1e-6:
            out.x = (tex_color.x * src_a + screen_color.x * dst_a * one_minus_src) / out_a
            out.y = (tex_color.y * src_a + screen_color.y * dst_a * one_minus_src) / out_a
            out.z = (tex_color.z * src_a + screen_color.z * dst_a * one_minus_src) / out_a
            out.w = out_a
        else:
            out = ti.math.vec4(0.0, 0.0, 0.0, 0.0)

    return out


@ti.kernel
def render_layer_no_mask1(
        image: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),  # 纹理图像
        x: ti.f32, y: ti.f32,  # 精灵位置
        cx: ti.f32, cy: ti.f32,  # 纹理中心
        scale_x: ti.f32, scale_y: ti.f32,  # 缩放
        rotation: ti.f32,  # 旋转
        alpha: ti.f32,  # 透明度
        width: ti.f32, height: ti.f32,  # 纹理尺寸
        min_x: ti.i32, max_x: ti.i32,  # 包围盒
        min_y: ti.i32, max_y: ti.i32,
        screen: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
):
    """
    精灵渲染 kernel
    """
    # 是否使用插值采样（缩放比例不是1.0）
    use_scale = (ti.abs(scale_x - 1.0) > 1e-6 or ti.abs(scale_y - 1.0) > 1e-6)

    # 预计算逆旋转矩阵（用于从屏幕坐标反推纹理坐标）
    cos_rot = ti.cos(rotation)
    sin_rot = ti.sin(rotation)
    rot_matrix = ti.math.mat2(cos_rot, sin_rot, -sin_rot, cos_rot)

    for x_screen, y_screen in ti.ndrange((min_x, max_x + 1), (min_y, max_y + 1)):
        screen_color = screen[x_screen, y_screen]

        # 计算屏幕像素相对于精灵中心的偏移
        dx = x_screen - (x + cx)
        dy = y_screen - (y + cy)
        screen_offset = ti.math.vec2(dx, dy)

        # 逆变换：屏幕偏移 → 纹理本地坐标（先缩放再旋转）
        scaled_offset = ti.math.vec2(screen_offset.x / scale_x, screen_offset.y / scale_y)
        tex_offset = rot_matrix @ scaled_offset
        tex_x_f = cx + tex_offset.x
        tex_y_f = cy + tex_offset.y

        # 边界检查：纹理坐标越界则跳过
        if not (0 <= tex_x_f < width and 0 <= tex_y_f < height):
            continue

        tex_color = screen_color
        if not use_scale:
            tex_color = image[ti.cast(tex_x_f, ti.i32), ti.cast(tex_y_f, ti.i32)]
        else:
            # 双三次插值
            # tex_color = bicubic_sample(image, tex_x_f, tex_y_f, width, height)
            # 线性插值
            tex_color = bilinear_sample(image, tex_x_f, tex_y_f, width, height)  # 双线性
            # lanczos4
            #tex_color = lanczos4_sample(image, tex_x_f, tex_y_f, width, height)  # Lanczos4

        # Alpha混合
        final_alpha = ti.min(tex_color.w * alpha, 1.0)

        # 透明度过低则跳过
        if final_alpha <= 1e-6:
            continue

        new_color = tex_color
        if final_alpha >= 0.999:
            new_color.w = 1.0
        else:
            new_color = ti.math.mix(screen_color, tex_color, final_alpha)
            new_color.w = 1.0
        screen[x_screen, y_screen] = new_color


@ti.kernel
def render_layer_with_mask1(
        image: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),  # 纹理图像
        x: ti.f32, y: ti.f32,  # 精灵位置
        cx: ti.f32, cy: ti.f32,  # 纹理中心
        scale_x: ti.f32, scale_y: ti.f32,  # 缩放
        rotation: ti.f32,  # 旋转
        alpha: ti.f32,  # 透明度
        width: ti.f32, height: ti.f32,  # 纹理尺寸
        min_x: ti.i32, max_x: ti.i32,  # 包围盒
        min_y: ti.i32, max_y: ti.i32,
        mask: ti.types.ndarray(dtype=ti.f32, ndim=2),  # 遮罩（mask2d.field，标量 f32）
        screen: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),  # 输出屏幕
):
    # 是否使用插值采样（缩放比例不是1.0）
    use_scale = (ti.abs(scale_x - 1.0) > 1e-6 or ti.abs(scale_y - 1.0) > 1e-6)

    # 预计算逆旋转矩阵（用于从屏幕坐标反推纹理坐标）
    cos_rot = ti.cos(rotation)
    sin_rot = ti.sin(rotation)
    rot_matrix = ti.math.mat2(cos_rot, sin_rot, -sin_rot, cos_rot)

    # 只遍历包围盒区域
    for x_screen, y_screen in ti.ndrange((min_x, max_x + 1), (min_y, max_y + 1)):
        screen_color = screen[x_screen, y_screen]

        # 计算屏幕像素相对于精灵中心的偏移
        dx = x_screen - (x + cx)
        dy = y_screen - (y + cy)
        screen_offset = ti.math.vec2(dx, dy)

        # 逆变换：屏幕偏移 → 纹理本地坐标（先缩放再旋转）
        scaled_offset = ti.math.vec2(screen_offset.x / scale_x, screen_offset.y / scale_y)
        tex_offset = rot_matrix @ scaled_offset
        tex_x_f = cx + tex_offset.x
        tex_y_f = cy + tex_offset.y

        # 边界检查：纹理坐标越界则跳过
        if not (0 <= tex_x_f < width and 0 <= tex_y_f < height):
            continue

        tex_x_i = ti.cast(tex_x_f, ti.i32)
        tex_y_i = ti.cast(tex_y_f, ti.i32)

        tex_color = screen_color
        if not use_scale:
            tex_color = image[ti.cast(tex_x_f, ti.i32), ti.cast(tex_y_f, ti.i32)]
        else:
            # 双三次插值
            # tex_color = bicubic_sample(image, tex_x_f, tex_y_f, width, height)
            # 线性插值
            tex_color = bilinear_sample(image, tex_x_f, tex_y_f, width, height)  # 双线性
            # lanczos4
            # tex_color = lanczos4_sample(image, tex_x_f, tex_y_f, width, height)  # Lanczos4

        # Alpha混合
        final_alpha = ti.min(tex_color.w * alpha, 1.0)

        # 应用 mask（mask2d 是标量 field，直接访问）
        mask_value = mask[tex_x_i, tex_y_i]
        final_alpha *= mask_value

        # 透明度过低则跳过
        if final_alpha <= 1e-6:
            continue

        new_color = tex_color
        # Alpha 混合优化：当 alpha 接近 1.0 时直接覆盖，避免混合计算
        if final_alpha >= 0.999:
            new_color.w = 1.0
        else:
            new_color = ti.math.mix(screen_color, tex_color, final_alpha)
            new_color.w = 1.0
        screen[x_screen, y_screen] = new_color

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

    screen_w, screen_h = screen.shape
    max_x1 = ti.math.min(max_x, screen_w)
    max_y1 = ti.math.min(max_y, screen_h)

    for x_screen, y_screen in ti.ndrange((min_x, max_x1), (min_y, max_y1)):
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

    screen_w, screen_h = screen.shape
    max_x1 = ti.math.min(max_x, screen_w)
    max_y1 = ti.math.min(max_y, screen_h)
    for x_screen, y_screen in ti.ndrange((min_x, max_x1), (min_y, max_y1)):
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

