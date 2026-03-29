

"""
#!/usr/bin/env python3
# /// script
# requires-python = ">=3.8"
# dependencies = [
#     "taichi>=1.7.4",
# ]
# ///


TCM 生成器 - 使用 Compute Graph 方式

编译 Taichi kernel 为 AOT 模块（TCM 文件）
支持架构：CPU (x64/ARM64)、CUDA、Vulkan

使用方法:
    uv run ./ti/kernel/generate_tcm.py
"""

import sys
from pathlib import Path
import taichi

# 添加根目录到 Python 路径
parent_dir = Path(__file__).parent.parent.parent
if str(parent_dir) not in sys.path:
    sys.path.insert(0, str(parent_dir))

from ti.kernel.layer import render_layer_no_mask, render_layer_with_mask
from ti.kernel.image import (
    fill_color,
    blur_box,
    blur_gaussian,
    blur_mosaic,
)
from ti.kernel.resize import resize_lanczos, resize_bicubic, resize_bilinear, resize_nearest
from ti.kernel.sdf import (
    compute_normalized_coords,
    compute_circle,
    compute_diamond,
    compute_rect,
    compute_directional,
    compute_triangle,
    compute_star,
    compute_heart,
    compute_cross,
)
from ti.kernel.mask import (
    image_to_mask,
    compute_distance_field,
    feather_linear,
    feather_conic,
    feather_smoothstep,
    feather_sigmoid,
)

# 要导出的 kernel 列表
kernels = [
    render_layer_no_mask,
    render_layer_with_mask,
    fill_color,
    resize_lanczos,
    resize_bicubic,
    resize_bilinear,
    resize_nearest,
    blur_box,
    blur_gaussian,
    blur_mosaic,
    compute_normalized_coords,
    compute_circle,
    compute_diamond,
    compute_rect,
    compute_directional,
    compute_triangle,
    compute_star,
    compute_heart,
    compute_cross,

    image_to_mask,
    compute_distance_field,
    feather_linear,
    feather_conic,
    feather_smoothstep,
    feather_sigmoid,
]

# 架构名称映射
architectures = {
    taichi.cpu: "cpu",
    taichi.cuda: "cuda",
    taichi.vulkan: "vulkan",
}


def create_tcm_module(arch, output_file: Path):
    """
    为指定架构创建 TCM 模块

    Args:
        arch: Taichi 架构
        output_file: 输出文件路径
    """

    try:
        # 重置并初始化 Taichi
        taichi.reset()
        taichi.init(arch=arch)
        m = taichi.aot.Module(arch)

        for kernel in kernels:
            print(f" - {kernel.__name__}")
            m.add_kernel(kernel)

        # 保存为 TCM 文件
        m.archive(str(output_file))

        print(f"{output_file} 生成成功!")

    except Exception as e:
        print(f"[ERROR] 生成 {output_file} 失败: {e}")
        import traceback
        traceback.print_exc()


def main():
    """主函数"""
    print(f"Taichi {taichi.__version__} TCM 生成器")

    # 输出目录
    output_dir = Path(__file__).parent.parent / "tcm"
    print(f"输出目录: {output_dir}")


    # 为每个架构生成 TCM
    for arch, arch_name in architectures.items():
        print()
        output_file = output_dir / f"{arch_name}.tcm"

        create_tcm_module(arch, output_file)


if __name__ == "__main__":
    main()
