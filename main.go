package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/go-mixed/go-canvas/ctypes"
	"github.com/go-mixed/go-canvas/font"
	"github.com/go-mixed/go-canvas/internel/misc"
	"github.com/go-mixed/go-canvas/render"
	"github.com/go-mixed/go-canvas/ti"
	"github.com/go-mixed/go-taichi/taichi"
)

func main() {
	totalStart := time.Now()
	cd := misc.GetCurrentDir()

	fmt.Println("=== 初始化阶段 ===")

	// 1. 创建渲染器
	t := time.Now()
	renderer, err := render.NewRenderer(taichi.ArchCuda)
	if err != nil {
		panic(err)
	}
	defer renderer.Release()
	fmt.Printf("[1/8] 创建渲染器: %v\n", time.Since(t))

	rect := ctypes.RectWH[int](0, 0, 720, 1280)

	// 2. 创建舞台
	t = time.Now()
	stage, err := render.NewStage(renderer, rect.Width(), rect.Height(), render.WithRawImage(true))
	if err != nil {
		panic(err)
	}
	defer stage.Release()
	fmt.Printf("[2/8] 创建舞台: %v\n", time.Since(t))

	//3. 创建背景块
	t = time.Now()
	background, err := render.NewBlockSprite(stage, ctypes.Attr().SetRect(rect))
	if err != nil {
		panic(err)
	}
	background.Fill(ctypes.RGBA(0x00ff00ff))
	fmt.Printf("[3/8] 创建背景块: %v\n", time.Since(t))

	// 4. 创建容器
	//t = time.Now()
	//container, err := render.NewContainer(stage, ti.Attr().SetRect(rect))
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("[4/8] 创建容器: %v\n", time.Since(t))
	//
	fontLibrary, err := font.NewFontLibrary(font.FontOpt().SetLogger(log.Default()))
	if err != nil {
		panic(err)
	}

	// 5. 加载图片
	t = time.Now()
	img, err := render.NewImageSprite(stage, ctypes.Attr().SetRect(rect), filepath.Join(cd, "examples", "1.jpg"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("[5/8] 加载图片: %v\n", time.Since(t))
	_ = img
	// 6. 创建形状遮罩
	//t = time.Now()
	//mask, err := render.NewShapeMask(img, ti.Attr().SetWH(img.Width(), img.Height()))
	//if err != nil {
	//	panic(err)
	//}
	//mask.DrawShape(ti.ShapeTypeCircle, 0.5, ti.ShapeOpt())
	////img.Blur(ti.BlurModeMosaic, 20)
	//fmt.Printf("[6/8] 创建遮罩: %v\n", time.Since(t))

	// 7. 创建文字精灵 1
	t = time.Now()
	text1, err := render.NewTextSprite(stage, fontLibrary, ctypes.Attr().SetHeight(60), font.RTOpt().SetAlign(ctypes.HAlignCenter, ctypes.VAlignTop))
	if err != nil {
		panic(err)
	}
	text1.SetText("<text font-size='40'>Interesting</text> <text font-size='40' color='#ffffff'>qg</text><text font-size='40' color='#ff0000'>World!</text>")
	fmt.Printf("[7/9] 创建文字1: %v\n", time.Since(t))

	// 8. 创建文字精灵 2
	t = time.Now()
	text2, err := render.NewTextSprite(stage, fontLibrary, ctypes.Attr().SetY(text1.Height()+20).SetWidth(500), font.RTOpt().SetAlign(ctypes.HAlignCenter, ctypes.VAlignMiddle))
	if err != nil {
		panic(err)
	}

	text2.SetText("<text font-size='50' font-color='#00ff00'>你们</text>\n<text font-size='50' color='#ffffff'>在那里吃饭。</text>\n<text font-size='60' color='#ff0000'>现在好</text>\n<text  font-color='#00ff00'>吗？</text>")

	fmt.Printf("[8/9] 创建文字2: %v\n", time.Since(t))

	// 9. 渲染
	t = time.Now()
	stage.Render(0)
	fmt.Printf("[9/9] 渲染: %v\n", time.Since(t))

	// 保存
	fmt.Println("\n=== 保存阶段 ===")
	t = time.Now()
	err = ti.SaveTiImageToFile(stage.Texture(), filepath.Join(misc.GetCurrentDir(), "out.png"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("[保存] out.png: %v\n", time.Since(t))

	var buf = make([]uint32, stage.Width()*stage.Height())
	err = stage.GetBgraImage(buf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%x", buf[:16])

	fmt.Printf("\n=== 总耗时: %v ===\n", time.Since(totalStart))
}
