import taichi as ti

from ti.kernel.sample import *

@ti.kernel
def resize_nearest(
        src: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        dst: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        src_width: ti.f32, src_height: ti.f32,
        dst_width: ti.f32, dst_height: ti.f32,
        scale_x: ti.f32, scale_y: ti.f32,
        offset_x: ti.f32, offset_y: ti.f32,
):
    """最近邻缩放"""
    for x, y in ti.ndrange(int(dst_width), int(dst_height)):
        # 反向映射到源坐标
        src_x = (ti.f32(x) - offset_x) / scale_x
        src_y = (ti.f32(y) - offset_y) / scale_y

        # 边界检查
        if 0 <= src_x < src_width and 0 <= src_y < src_height:
            dst[x, y] = src[ti.cast(src_x, ti.i32), ti.cast(src_y, ti.i32)]


@ti.kernel
def resize_bilinear(
        src: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        dst: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        src_width: ti.f32, src_height: ti.f32,
        dst_width: ti.f32, dst_height: ti.f32,
        scale_x: ti.f32, scale_y: ti.f32,
        offset_x: ti.f32, offset_y: ti.f32,
):
    """双线性插值缩放"""
    for x, y in ti.ndrange(int(dst_width), int(dst_height)):
        # 反向映射到源坐标
        src_x = (ti.f32(x) - offset_x) / scale_x
        src_y = (ti.f32(y) - offset_y) / scale_y

        # 边界检查
        if 0 <= src_x < src_width and 0 <= src_y < src_height:
            dst[x, y] = bilinear_sample(src, src_x, src_y, src_width, src_height)


@ti.kernel
def resize_bicubic(
        src: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        dst: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        src_width: ti.f32, src_height: ti.f32,
        dst_width: ti.f32, dst_height: ti.f32,
        scale_x: ti.f32, scale_y: ti.f32,
        offset_x: ti.f32, offset_y: ti.f32,
):
    """双三次插值缩放"""
    for x, y in ti.ndrange(int(dst_width), int(dst_height)):
        # 反向映射到源坐标
        src_x = (ti.f32(x) - offset_x) / scale_x
        src_y = (ti.f32(y) - offset_y) / scale_y

        # 边界检查
        if 0 <= src_x < src_width and 0 <= src_y < src_height:
            dst[x, y] = bicubic_sample(src, src_x, src_y, src_width, src_height)


@ti.kernel
def resize_lanczos(
        src: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        dst: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        src_width: ti.f32, src_height: ti.f32,
        dst_width: ti.f32, dst_height: ti.f32,
        scale_x: ti.f32, scale_y: ti.f32,
        offset_x: ti.f32, offset_y: ti.f32,
):
    """Lanczos4 插值缩放（质量最高）"""
    for x, y in ti.ndrange(int(dst_width), int(dst_height)):
        # 反向映射到源坐标
        src_x = (ti.f32(x) - offset_x) / scale_x
        src_y = (ti.f32(y) - offset_y) / scale_y

        # 边界检查
        if 0 <= src_x < src_width and 0 <= src_y < src_height:
            dst[x, y] = lanczos4_sample(src, src_x, src_y, src_width, src_height)

