package main

import (
	"fmt"
	"image"
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

	stage, err := render.NewStage(renderer, 720, 1280)
	if err != nil {
		panic(err)
	}
	defer stage.Release() // 必须释放

	fmt.Printf("init stage: %v\n", time.Since(t))

	t = time.Now()
	background, err := render.NewBlockSprite(stage, 720, 1280)
	if err != nil {
		panic(err)
	}
	background.Fill(color.White)
	fmt.Printf("init background: %v\n", time.Since(t))

	t = time.Now()
	img, err := render.NewImageSprite(stage, filepath.Join(cd, "examples", "1.jpg"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("init image: %v\n", time.Since(t))

	t = time.Now()
	img.Resize(720, 1280)
	//img.Blur(ti.BlurModeMosaic, 20)
	fmt.Printf("resize image: %v\n", time.Since(t))

	img.SetX(-200)

	t = time.Now()
	mask, err := render.NewShapeMask(img, 720, 1280, 720*0.5, 1280*0.5)
	if err != nil {
		panic(err)
	}
	mask.DrawShape(ti.ShapeTypeCircle, 0.5)

	fontLibrary := font.NewFontLibrary()

	text, err := render.NewTextSprite(stage, fontLibrary, 720, 1280, font.WithAlign(ti.HAlignCenter, ti.VAlignMiddle))
	if err != nil {
		panic(err)
	}
	text.SetText("<text font-size='50'>Hello</text>\n <text font-size='60' color='#ff0000'>World</text>!")

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

	img2 := image.NewRGBA(image.Rect(0, 0, 720, 1080))
	t = time.Now()
	for i, b := range buf {
		img2.Set(i%720, i/720, ti.ARGB(b))
	}
	fmt.Printf("save out2: %v\n", time.Since(t))
	misc.SaveImage(img2, filepath.Join(misc.GetCurrentDir(), "out2.png"))

}
