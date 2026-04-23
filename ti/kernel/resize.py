import taichi as ti

from ti.kernel.sample import *

@ti.kernel
def resize_nearest(
        src: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        dst: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        src_x: ti.f32, src_y: ti.f32, src_w: ti.f32, src_h: ti.f32,
        dst_x: ti.f32, dst_y: ti.f32, dst_w: ti.f32, dst_h: ti.f32,
):
    """最近邻缩放：将 src 的 [src_x, src_y, src_w, src_h] 区块缩放到 dst 的 [dst_x, dst_y, dst_w, dst_h] 区块"""
    scale_x = src_w / dst_w
    scale_y = src_h / dst_h
    dst_xi = ti.cast(dst_x, ti.i32)
    dst_yi = ti.cast(dst_y, ti.i32)
    src_x_end = src_x + src_w
    src_y_end = src_y + src_h
    for x, y in ti.ndrange(int(dst_w), int(dst_h)):
        sx = src_x + ti.f32(x) * scale_x
        sy = src_y + ti.f32(y) * scale_y
        if 0 <= sx < src_x_end and 0 <= sy < src_y_end:
            dst[dst_xi + x, dst_yi + y] = src[ti.cast(sx, ti.i32), ti.cast(sy, ti.i32)]


@ti.kernel
def resize_bilinear(
        src: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        dst: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        src_x: ti.f32, src_y: ti.f32, src_w: ti.f32, src_h: ti.f32,
        dst_x: ti.f32, dst_y: ti.f32, dst_w: ti.f32, dst_h: ti.f32,
):
    """双线性插值缩放：将 src 的 [src_x, src_y, src_w, src_h] 区块缩放到 dst 的 [dst_x, dst_y, dst_w, dst_h] 区块"""
    scale_x = src_w / dst_w
    scale_y = src_h / dst_h
    dst_xi = ti.cast(dst_x, ti.i32)
    dst_yi = ti.cast(dst_y, ti.i32)
    src_x_end = src_x + src_w
    src_y_end = src_y + src_h
    src_total_w = ti.cast(src.shape[0], ti.f32)
    src_total_h = ti.cast(src.shape[1], ti.f32)
    for x, y in ti.ndrange(int(dst_w), int(dst_h)):
        sx = src_x + ti.f32(x) * scale_x
        sy = src_y + ti.f32(y) * scale_y
        if 0 <= sx < src_x_end and 0 <= sy < src_y_end:
            dst[dst_xi + x, dst_yi + y] = bilinear_sample(src, sx, sy, src_total_w, src_total_h)


@ti.kernel
def resize_bicubic(
        src: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        dst: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        src_x: ti.f32, src_y: ti.f32, src_w: ti.f32, src_h: ti.f32,
        dst_x: ti.f32, dst_y: ti.f32, dst_w: ti.f32, dst_h: ti.f32,
):
    """双三次插值缩放：将 src 的 [src_x, src_y, src_w, src_h] 区块缩放到 dst 的 [dst_x, dst_y, dst_w, dst_h] 区块"""
    scale_x = src_w / dst_w
    scale_y = src_h / dst_h
    dst_xi = ti.cast(dst_x, ti.i32)
    dst_yi = ti.cast(dst_y, ti.i32)
    src_x_end = src_x + src_w
    src_y_end = src_y + src_h
    src_total_w = ti.cast(src.shape[0], ti.f32)
    src_total_h = ti.cast(src.shape[1], ti.f32)
    for x, y in ti.ndrange(int(dst_w), int(dst_h)):
        sx = src_x + ti.f32(x) * scale_x
        sy = src_y + ti.f32(y) * scale_y
        if 0 <= sx < src_x_end and 0 <= sy < src_y_end:
            dst[dst_xi + x, dst_yi + y] = bicubic_sample(src, sx, sy, src_total_w, src_total_h)


@ti.kernel
def resize_lanczos(
        src: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        dst: ti.types.ndarray(element_shape=(4,), dtype=ti.f32, ndim=2),
        src_x: ti.f32, src_y: ti.f32, src_w: ti.f32, src_h: ti.f32,
        dst_x: ti.f32, dst_y: ti.f32, dst_w: ti.f32, dst_h: ti.f32,
):
    """Lanczos4 插值缩放（质量最高）：将 src 的 [src_x, src_y, src_w, src_h] 区块缩放到 dst 的 [dst_x, dst_y, dst_w, dst_h] 区块"""
    scale_x = src_w / dst_w
    scale_y = src_h / dst_h
    dst_xi = ti.cast(dst_x, ti.i32)
    dst_yi = ti.cast(dst_y, ti.i32)
    src_x_end = src_x + src_w
    src_y_end = src_y + src_h
    src_total_w = ti.cast(src.shape[0], ti.f32)
    src_total_h = ti.cast(src.shape[1], ti.f32)
    for x, y in ti.ndrange(int(dst_w), int(dst_h)):
        sx = src_x + ti.f32(x) * scale_x
        sy = src_y + ti.f32(y) * scale_y
        if 0 <= sx < src_x_end and 0 <= sy < src_y_end:
            dst[dst_xi + x, dst_yi + y] = lanczos4_sample(src, sx, sy, src_total_w, src_total_h)
