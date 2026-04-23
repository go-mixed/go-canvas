import matplotlib.pyplot as plt
import matplotlib.ticker as ticker
import numpy as np
from matplotlib.patches import Rectangle

# --- 严格参数 设置 ---
width_px = 720
height_px = 1280
dpi = 100

# 1像素边框的物理磅值。在 100 DPI 下，1px = 1 point。
# 增加一点 linewidth (1.5) 可以确保在所有查看器中都能看到清晰的 1px。
border_lw_pts = 1.0
border_color = '#ff0000' # 深绿色

# 计算中心点
center_x = width_px // 2
center_y = height_px // 2

# 创建画布，设置浅绿色背景
fig = plt.figure(figsize=(width_px / dpi, height_px / dpi), dpi=dpi, facecolor='#dcfce7')
ax = fig.add_axes([0, 0, 1, 1]) # 占满全图

# 1. 设置 Axes 背景色
ax.set_facecolor('#dcfce7')

# 2. **核心修正：精密绘制 1px 外边框**
# 关键在于：为了得到完美的 1px 物理边框，我们需要让 Rectangle 的
# 物理中心正好位于画布的最外沿。我们通过将坐标轴范围向外推半个像素来实现。
ax.set_xlim(-center_x - 0.5, center_x + 0.5)
ax.set_ylim(center_y + 0.5, -center_y - 0.5) # Y轴向下为正

# 创建一个占满全图的空心矩形。
# x, y 原点是 -center_x - 0.5，宽度是 width_px
rect = Rectangle((-center_x - 0.5, -center_y - 0.5), width_px, height_px,
                 fill=False, color=border_color, linewidth=border_lw_pts, zorder=100)
ax.add_patch(rect)

# 3. 将坐标轴移动到中心 (0, 0)
ax.spines['left'].set_position('center')
ax.spines['bottom'].set_position('center')
ax.spines['left'].set_color(border_color)   # 设置中心轴颜色
ax.spines['bottom'].set_color(border_color)
ax.spines['right'].set_color('none')
ax.spines['top'].set_color('none')

# 4. 配置刻度逻辑
ax.xaxis.set_major_locator(ticker.MultipleLocator(100))
ax.yaxis.set_major_locator(ticker.MultipleLocator(100))
ax.xaxis.set_minor_locator(ticker.MultipleLocator(10))
ax.yaxis.set_minor_locator(ticker.MultipleLocator(10))

# 5. 刻度样式
ax.tick_params(axis='both', which='major', direction='inout', length=20, width=1.5, color=border_color)
ax.tick_params(axis='both', which='minor', direction='inout', length=10, width=1.0, color=border_color)

# 6. 绘制 1px 极细刻度
# 使用 plot 绘制，这里可以设置 alpha
for x in range(-center_x, center_x, 1):
    if x % 10 != 0:
        ax.plot([x, x], [-5, 5], color=border_color, linewidth=0.5, alpha=0.2)

for y in range(-center_y, center_y, 1):
    if y % 10 != 0:
        ax.plot([-5, 5], [y, y], color=border_color, linewidth=0.5, alpha=0.2)

# 7. 标注数字 (100, 200...)
label_color = '#166534'
for x in range(-center_x + 100, center_x, 100):
    if x != 0:
        ax.text(x, 45, f'{x}', color=label_color, fontsize=9, ha='center', fontweight='bold')

for y in range(-center_y + 100, center_y, 100):
    if y != 0:
        ax.text(15, y, f'{y}', color=label_color, fontsize=9, va='center', fontweight='bold')

# 标注原点
ax.text(12, 35, '0,0', color=label_color, fontsize=10, fontweight='bold')

# 8. 彻底隐藏默认标签
ax.set_xticklabels([])
ax.set_yticklabels([])

# 9. 保存
plt.savefig('ruler.png', dpi=dpi, pad_inches=0)
print("带有物理 1px 外边框的中心对齐像素尺已生成：ruler.png")
print("如果看不清边框，请将图片放大到 100% 物理像素查看。")
plt.show()