import taichi as ti
import taichi.math as tm

@ti.kernel
def copy_region(
    src: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    dst: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    src_x: ti.i32, src_y: ti.i32, src_w: ti.i32, src_h: ti.i32,
    dst_x: ti.i32, dst_y: ti.i32, dst_w: ti.i32, dst_h: ti.i32,
):
    """将 src 的 [src_x, src_y, src_w, src_h] 区块直接复制到 dst 的 [dst_x, dst_y, dst_w, dst_h] 区块，无插值"""
    copy_w = ti.min(src_w, dst_w)
    copy_h = ti.min(src_h, dst_h)
    for i, j in ti.ndrange(copy_w, copy_h):
        sx = src_x + i
        sy = src_y + j
        dx = dst_x + i
        dy = dst_y + j
        if 0 <= sx < src.shape[0] and 0 <= sy < src.shape[1] \
                and 0 <= dx < dst.shape[0] and 0 <= dy < dst.shape[1]:
            dst[dx, dy] = src[sx, sy]


@ti.kernel
def ti_image_to_bgra(
    texture: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    output: ti.types.ndarray(dtype=ti.u32, ndim=2)
):
    for x, y in texture:
        rgba = texture[x, y]
        rgba = ti.math.clamp(rgba, 0.0, 1.0)

        # Export premultiplied BGRA for better downstream compositing consistency
        # in raw/video pipelines (e.g. ffmpeg).
        a_f = rgba[3]
        r_f = 0.0
        g_f = 0.0
        b_f = 0.0
        if a_f <= 1e-6:
            a_f = 0.0
        else:
            r_f = rgba[0] * a_f
            g_f = rgba[1] * a_f
            b_f = rgba[2] * a_f

        r = ti.cast(r_f * 255.0 + 0.5, ti.u32)
        g = ti.cast(g_f * 255.0 + 0.5, ti.u32)
        b = ti.cast(b_f * 255.0 + 0.5, ti.u32)
        a = ti.cast(a_f * 255.0 + 0.5, ti.u32)

        # BGRA 打包（ bgra 小端内存布局）
        output[y, x] = (a << 24) | (r << 16) | (g << 8) | b

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

@ti.kernel
def draw_line(
    dst: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    x1: ti.i32, y1: ti.i32, x2: ti.i32, y2: ti.i32,
    color: ti.types.vector(4, ti.f32),
):
    """
    绘制线段
    """
    # 1. 计算线段在 X 和 Y 方向的跨度
    dx = x2 - x1
    dy = y2 - y1

    # 2. 确定步数（取长边的像素数，保证线条连续无断点）
    steps = ti.max(ti.abs(dx), ti.abs(dy))

    # 3. 并行化处理：每一个 i 都是线段上的一个像素点
    for i in range(steps + 1):
        # 计算当前点的坐标 (线性插值)
        # 使用 float 计算位置以保证精度，最后转回 i32 索引
        curr_x = ti.cast(x1 + i * dx / steps, ti.i32)
        curr_y = ti.cast(y1 + i * dy / steps, ti.i32)

        # 4. 边界检查（防止 ndarray 越界崩溃）
        img_w, img_h = dst.shape[0], dst.shape[1]
        if 0 <= curr_x < img_w and 0 <= curr_y < img_h:
            dst[curr_x, curr_y] = color


@ti.kernel
def draw_rect(
    dst: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    x: ti.i32, y: ti.i32, w: ti.i32, h: ti.i32,
    color: ti.types.vector(4, ti.f32),
):
    # 1. 确定图像边界
    img_w, img_h = dst.shape[0], dst.shape[1]

    # 2. 计算有效的起始点和终止点（防止 w, h 为负数或超出屏幕）
    x_min = ti.max(0, x)
    y_min = ti.max(0, y)
    x_max = ti.min(img_w, x + w)
    y_max = ti.min(img_h, y + h)

    # 3. 并行填充：ti.ndrange 会自动在 GPU 上展开为二维并行循环
    for i, j in ti.ndrange((x_min, x_max), (y_min, y_max)):
        dst[i, j] = color


@ti.kernel
def render_border(
    dst: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    border_widths: ti.types.vector(4, ti.f32), # [top, right, bottom, left]
    border_top_color: ti.types.vector(4, ti.f32),
    border_right_color: ti.types.vector(4, ti.f32),
    border_bottom_color: ti.types.vector(4, ti.f32),
    border_left_color: ti.types.vector(4, ti.f32),
    radii: ti.types.vector(4, ti.f32),         # [tl, tr, br, bl]
):
    res = tm.vec2(dst.shape[0], dst.shape[1])
    half_res = res * 0.5

    for i, j in dst:
        p = tm.vec2(i + 0.5, j + 0.5)

        # --- 1. 计算圆角矩形 SDF (dist > 0 在外, dist < 0 在内) ---
        p_centered = p - half_res
        # 选择象限对应的圆角半径
        r = 0.0
        if p_centered.x < 0 and p_centered.y < 0: r = radii[0] # TL
        elif p_centered.x >= 0 and p_centered.y < 0: r = radii[1] # TR
        elif p_centered.x >= 0 and p_centered.y >= 0: r = radii[2] # BR
        else: r = radii[3] # BL

        # 圆角矩形核心公式
        q = ti.abs(p_centered) - half_res + r
        dist_outer = tm.length(tm.max(q, 0.0)) + tm.min(tm.max(q.x, q.y), 0.0) - r

        # --- 2. 完美的裁剪逻辑 (Content Clipping) ---
        # 即使这里有浮点数精度误差，我们通过 Alpha 混合解决它
        alpha_out = 1.0 - tm.smoothstep(-0.5, 0.5, dist_outer)

        # --- 3. 完美的边框绘制逻辑 (Border Rendering) ---
        # 核心：必须处理四条不同颜色边的对角线切分
        uv = p / res
        active_color = tm.vec4(0.0)
        bw = 0.0

        if uv.x < uv.y: # 左下三角
            if uv.x < (1.0 - uv.y): # 左
                active_color = border_left_color
                bw = border_widths[3]
            else: # 下
                active_color = border_bottom_color
                bw = border_widths[2]
        else: # 右上三角
            if uv.x < (1.0 - uv.y): # 上
                active_color = border_top_color
                bw = border_widths[0]
            else: # 右
                active_color = border_right_color
                bw = border_widths[1]

        # 计算边框强度的 Mask
        # 此 Mask 确保边框是一个内沿圆角的环形区域
        alpha_in = 1.0 - tm.smoothstep(-bw - 0.5, -bw + 0.5, dist_outer)

        # --- 4. 无缝混合的核心逻辑 ---
        if alpha_out > 0.0:
            # 第一步：不管有没有边框，先把 content 裁剪成完美的圆角形状
            # (注意：这是在内存中修改 Alpha 通道)
            current_pixel = dst[i, j]
            current_pixel.w *= alpha_out

            # 第二步：将边框 Alpha 混合叠加在已裁剪的内容之上
            # 当 dist_outer 恰好等于边界时，content 的 alpha 会微小减小，
            # 但边框的 alpha 会精确填补这个空缺，两者共享一个 Alpha 通道边界。

            # 边框在 dist_outer 处的有效强度
            b_mask = alpha_out - alpha_in

            if b_mask > 0.0:
                # 使用标准的 Alpha 混合公式进行覆盖
                src = active_color * b_mask # 边框色带其 Mask

                # 最终颜色写入
                dst[i, j].rgb = src.rgb * src.w + current_pixel.rgb * current_pixel.w * (1.0 - src.w)
                dst[i, j].w = src.w + current_pixel.w * (1.0 - src.w)
            else:
                # 只有 content 区域，更新裁剪后的 alpha
                dst[i, j].w = current_pixel.w
        else:
            # 外部
            dst[i, j] = tm.vec4(0.0)

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

