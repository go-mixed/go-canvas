package main

import (
	"fmt"
	"image/color"
	"path/filepath"
	"time"

	"github.com/go-mixed/go-canvas/font"
	"github.com/go-mixed/go-canvas/misc"
	"github.com/go-mixed/go-canvas/render"
	"github.com/go-mixed/go-canvas/ti"
	"github.com/go-mixed/go-taichi/taichi"
)

func main() {
	t := time.Now()
	cd := misc.GetCurrentDir()

	renderer, err := render.NewRenderer(taichi.ArchCuda)
	if err != nil {
		panic(err)
	}
	defer renderer.Release() // 必须释放

	rect := ti.RectWH[int](0, 0, 720, 1280)

	stage, err := render.NewStage(renderer, rect.Width(), rect.Height())
	if err != nil {
		panic(err)
	}
	defer stage.Release() // 必须释放

	fmt.Printf("init stage: %v\n", time.Since(t))

	t = time.Now()
	background, err := render.NewBlockSprite(stage, ti.Attr().SetRect(rect))
	if err != nil {
		panic(err)
	}
	background.Fill(color.White)
	fmt.Printf("init background: %v\n", time.Since(t))

	fontLibrary := font.NewFontLibrary()

	container, err := render.NewContainer(stage, ti.Attr().SetRect(rect))
	if err != nil {
		panic(err)
	}

	//t = time.Now()
	//mask, err := render.NewShapeMask(img, 720, 1280, 720*0.5, 1280*0.5)
	//if err != nil {
	//	panic(err)
	//}
	//mask.DrawShape(ti.ShapeTypeCircle, 0.5)

	text, err := render.NewTextSprite(container, fontLibrary, ti.Attr().SetWidth(120), font.WithAlign(ti.HAlignCenter, ti.VAlignMiddle))
	if err != nil {
		panic(err)
	}
	text.SetText("<text font-size='50'>Hello1111</text>\n <text font-size='60' color='#ff0000'>World</text>!")

	t = time.Now()
	img, err := render.NewImageSprite(container, ti.Attr().SetXY(-50, 300), filepath.Join(cd, "examples", "1.jpg"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("init image: %v\n", time.Since(t))

	t = time.Now()
	img.Resize(300, 300)
	//img.Blur(ti.BlurModeMosaic, 20)
	fmt.Printf("resize image: %v\n", time.Since(t))

	text2, err := render.NewTextSprite(container, fontLibrary, ti.Attr().SetY(text.Height()+20), font.WithAlign(ti.HAlignCenter, ti.VAlignMiddle))
	if err != nil {
		panic(err)
	}

	text2.SetText("<text font-size='50' font-color=‘#00ff00’>你</text>\n <text font-size='60' color='#ff0000'>好adada</text>\n<text  font-color=‘#00ff00’>呢啊啊</text>")

	ti.SaveTiImageToFile(text2.Texture(), filepath.Join(misc.GetCurrentDir(), "out2.png"))
	stage.Render()
	fmt.Printf("render: %v\n", time.Since(t))

	t = time.Now()
	err = ti.SaveTiImageToFile(stage.Texture(), filepath.Join(misc.GetCurrentDir(), "out.png"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("save out1: %v\n", time.Since(t))

	buf := make([]uint32, 720*1280)

	t = time.Now()
	stage.ToBgraImage(buf)
	fmt.Printf("to bgra image: %v\n", time.Since(t))

}
