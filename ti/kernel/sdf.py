import taichi as ti

# signed distance field


@ti.kernel
def compute_normalized_coords(
    dx: ti.types.ndarray(dtype=ti.f32, ndim=2),
    dy: ti.types.ndarray(dtype=ti.f32, ndim=2),
    cx: ti.f32,
    cy: ti.f32
):
    # 计算从中心点到四个角的距离，取最大值作为归一化基准
    # 这样无论中心点在哪里，t=1 时都能覆盖整个屏幕

    w, h = dx.shape
    fw = ti.cast(w, ti.f32)
    fh = ti.cast(h, ti.f32)
    dist_to_top_left = ti.sqrt(cx * cx + cy * cy)
    dist_to_top_right = ti.sqrt((fw - cx) * (fw - cx) + cy * cy)
    dist_to_bottom_left = ti.sqrt(cx * cx + (fh - cy) * (fh - cy))
    dist_to_bottom_right = ti.sqrt((fw - cx) * (fw - cx) + (fh - cy) * (fh - cy))

    max_dist = ti.max(
        ti.max(dist_to_top_left, dist_to_top_right),
        ti.max(dist_to_bottom_left, dist_to_bottom_right)
    )
    scale = ti.max(max_dist, 1.0)

    # 归一化坐标网格（从中心到最远角的距离归一化为 1.0）
    for i, j in ti.ndrange(w, h):
        dx[i, j] = (i - cx) / scale
        dy[i, j] = (j - cy) / scale


@ti.kernel
def compute_circle(
    data: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    dx: ti.types.ndarray(dtype=ti.f32, ndim=2),
    dy: ti.types.ndarray(dtype=ti.f32, ndim=2),
    t_val: ti.f32,
    color: ti.types.vector(4, ti.f32)
):
    # 归一化半径（从中心到角的距离已归一化为 1.0）
    # t=1 时，radius=1.0，刚好覆盖整个屏幕（屏幕内切于圆）
    radius = t_val
    radius_sq = radius * radius

    for i, j in ti.ndrange(data.shape[0], data.shape[1]):
        dist_sq = dx[i, j] * dx[i, j] + dy[i, j] * dy[i, j]
        if dist_sq <= radius_sq:
            data[i, j] = color


@ti.kernel
def compute_diamond(
    data: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    dx: ti.types.ndarray(dtype=ti.f32, ndim=2),
    dy: ti.types.ndarray(dtype=ti.f32, ndim=2),
    t_val: ti.f32,
    color: ti.types.vector(4, ti.f32)
):
    for i, j in ti.ndrange(data.shape[0], data.shape[1]):
        manhattan_dist = ti.abs(dx[i, j]) + ti.abs(dy[i, j])
        if manhattan_dist <= t_val:
            data[i, j] = color


@ti.kernel
def compute_rect(
    data: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    dx: ti.types.ndarray(dtype=ti.f32, ndim=2),
    dy: ti.types.ndarray(dtype=ti.f32, ndim=2),
    t_val: ti.f32,
    dir_val: ti.i32,
    color: ti.types.vector(4, ti.f32)
):
    for i, j in ti.ndrange(data.shape[0], data.shape[1]):
        show = False
        if dir_val == 0:  # TOP
            show = dy[i, j] + 0.5 <= t_val
        elif dir_val == 1:  # BOTTOM
            show = dy[i, j] + 0.5 >= (1.0 - t_val)
        elif dir_val == 2:  # LEFT
            show = dx[i, j] + 0.5 <= t_val
        elif dir_val == 3:  # RIGHT
            show = dx[i, j] + 0.5 >= (1.0 - t_val)
        elif dir_val == 4:  # TOP_LEFT
            show = dx[i, j] <= t_val and dy[i, j] <= t_val
        elif dir_val == 5:  # TOP_RIGHT
            show = dx[i, j] >= (1.0 - t_val) and dy[i, j] <= t_val
        elif dir_val == 6:  # BOTTOM_LEFT
            show = dx[i, j] <= t_val and dy[i, j] >= (1.0 - t_val)
        elif dir_val == 7:  # BOTTOM_RIGHT
            show = dx[i, j] >= (1.0 - t_val) and dy[i, j] >= (1.0 - t_val)
        elif dir_val == 8:  # CENTER (radial rect)
            show = ti.abs(dx[i, j]) <= t_val and ti.abs(dy[i, j]) <= t_val

        if show:
            data[i, j] = color


@ti.kernel
def compute_directional(
    data: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    dx: ti.types.ndarray(dtype=ti.f32, ndim=2),
    dy: ti.types.ndarray(dtype=ti.f32, ndim=2),
    t_val: ti.f32,
    dir_x: ti.f32,
    dir_y: ti.f32,
    use_radial: ti.f32,
    manhattan_weight: ti.f32,
    chebyshev_weight: ti.f32,
    color: ti.types.vector(4, ti.f32)
):
    """
    通用方向性遮罩计算

    参数说明：
    - dir_x, dir_y: 方向向量（use_radial=0.0 时使用）
    - use_radial: 0.0=线性投影, 1.0=径向距离
    - manhattan_weight: 0.0=欧几里得, 1.0=曼哈顿
    - chebyshev_weight: 0.0=不使用, 1.0=切比雪夫
    """
    for i, j in ti.ndrange(data.shape[0], data.shape[1]):
        x, y = dx[i, j], dy[i, j]

        # 线性投影距离
        projection = x * dir_x + y * dir_y
        linear_dist = (projection + 1.0) * 0.5

        # 三种径向距离度量
        euclidean_dist = ti.sqrt(x * x + y * y)
        manhattan_dist = ti.abs(x) + ti.abs(y)
        chebyshev_dist = ti.max(ti.abs(x), ti.abs(y))

        # 混合径向距离
        mixed_dist = euclidean_dist * (1.0 - manhattan_weight) + manhattan_dist * manhattan_weight
        radial_dist = mixed_dist * (1.0 - chebyshev_weight) + chebyshev_dist * chebyshev_weight

        # 最终距离：线性 vs 径向
        final_dist = linear_dist * (1.0 - use_radial) + radial_dist * use_radial

        if final_dist <= t_val:
            data[i, j] = color


@ti.kernel
def compute_triangle(
    data: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    dx: ti.types.ndarray(dtype=ti.f32, ndim=2),
    dy: ti.types.ndarray(dtype=ti.f32, ndim=2),
    t_val: ti.f32,
    color: ti.types.vector(4, ti.f32)
):
    # 等边三角形：顶点向上，底边在下
    scaled_t = t_val * 2.0
    half_width = scaled_t * 0.5

    # 顶点位置（归一化坐标，y 向上为正）
    top_y = scaled_t * 0.5
    bottom_y = -scaled_t * 0.5

    for i, j in ti.ndrange(data.shape[0], data.shape[1]):
        x = dx[i, j]
        y = dy[i, j]

        # 使用重心坐标判断点是否在三角形内
        v0_x = 0.0
        v0_y = top_y
        v1_x = -half_width
        v1_y = bottom_y
        v2_x = half_width
        v2_y = bottom_y

        denom = (v1_y - v2_y) * (v0_x - v2_x) + (v2_x - v1_x) * (v0_y - v2_y)
        a = ((v1_y - v2_y) * (x - v2_x) + (v2_x - v1_x) * (y - v2_y)) / denom
        b = ((v2_y - v0_y) * (x - v2_x) + (v0_x - v2_x) * (y - v2_y)) / denom
        c = 1.0 - a - b

        if a >= 0.0 and b >= 0.0 and c >= 0.0:
            data[i, j] = color


@ti.kernel
def compute_star(
    data: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    dx: ti.types.ndarray(dtype=ti.f32, ndim=2),
    dy: ti.types.ndarray(dtype=ti.f32, ndim=2),
    t_val: ti.f32,
    color: ti.types.vector(4, ti.f32)
):
    # 五角星参数：外半径和内半径比例
    inner_ratio = 0.382  # 内半径 / 外半径 ≈ 0.382 (黄金比例)
    outer_radius = t_val * 2.7  # 放大2.7倍确保t=1时屏幕完全被覆盖
    inner_radius = outer_radius * inner_ratio

    for i, j in ti.ndrange(data.shape[0], data.shape[1]):
        x = dx[i, j]
        y = dy[i, j]

        angle_per_point = 2.0 * ti.math.pi / 10.0
        base_angle = -ti.math.pi * 0.5  # 从顶部开始

        # 射线法判断点是否在多边形内
        intersections = 0

        for k in range(10):
            angle1 = base_angle + k * angle_per_point
            angle2 = base_angle + (k + 1) * angle_per_point

            r1 = outer_radius if (k % 2) == 0 else inner_radius
            r2 = outer_radius if ((k + 1) % 2) == 0 else inner_radius

            x1 = r1 * ti.cos(angle1)
            y1 = r1 * ti.sin(angle1)
            x2 = r2 * ti.cos(angle2)
            y2 = r2 * ti.sin(angle2)

            # 检查边是否跨越射线的y坐标
            if (y1 <= y and y2 > y) or (y2 <= y and y1 > y):
                t = (y - y1) / (y2 - y1)
                x_intersect = x1 + t * (x2 - x1)

                if x_intersect > x:
                    intersections += 1

        if (intersections % 2) == 1:
            data[i, j] = color


@ti.kernel
def compute_heart(
    data: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    dx: ti.types.ndarray(dtype=ti.f32, ndim=2),
    dy: ti.types.ndarray(dtype=ti.f32, ndim=2),
    t_val: ti.f32,
    color: ti.types.vector(4, ti.f32)
):
    # 心形的最窄部分在凹陷处，需要更大的缩放以确保完全覆盖屏幕
    scaled_t = t_val * 2.5

    for i, j in ti.ndrange(data.shape[0], data.shape[1]):
        x = dx[i, j] / scaled_t * 2.0
        y = -dy[i, j] / scaled_t * 1.5 + 0.3
        term1 = ti.pow(x * x + y * y - 1.0, 3.0)
        term2 = x * x * y * y * y
        heart_shape = term1 - term2 <= 0.0

        if heart_shape:
            data[i, j] = color


@ti.kernel
def compute_cross(
    data: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
    dx: ti.types.ndarray(dtype=ti.f32, ndim=2),
    dy: ti.types.ndarray(dtype=ti.f32, ndim=2),
    t_val: ti.f32,
    color: ti.types.vector(4, ti.f32)
):
    # 十字架：两条相交的矩形条，从中心向外扩展
    arm_width_ratio = 0.3  # 臂宽与长度的比例
    scaled_t = t_val * 1.5
    arm_width = arm_width_ratio * scaled_t

    for i, j in ti.ndrange(data.shape[0], data.shape[1]):
        x = dx[i, j]
        y = dy[i, j]

        horizontal = ti.abs(y) <= arm_width and ti.abs(x) <= scaled_t
        vertical = ti.abs(x) <= arm_width and ti.abs(y) <= scaled_t

        if horizontal or vertical:
            data[i, j] = color