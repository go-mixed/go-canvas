package main

import (
	"fmt"
	"image/color"
	"path/filepath"
	"time"

	"github.com/go-mixed/go-canvas/misc"
	"github.com/go-mixed/go-canvas/render"
	"github.com/go-mixed/go-canvas/ti"
	"github.com/go-mixed/go-taichi/taichi"
)

func main() {
	t := time.Now()
	cd := misc.GetCurrentDir()

	stage, err := render.NewStage(taichi.ArchCuda, 720, 1280)
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
	background.FillTexture(color.White)
	fmt.Printf("init background: %v\n", time.Since(t))

	t = time.Now()
	img, err := render.NewImageSprite(stage, filepath.Join(cd, "examples", "1.jpg"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("init image: %v\n", time.Since(t))

	t = time.Now()
	img.ResizeTo(720, 1280)
	fmt.Printf("resize image: %v\n", time.Since(t))

	img.SetX(-200)

	t = time.Now()
	mask, err := render.NewShapeMask(img, 720, 1280, 720*0.5, 1280*0.5)
	if err != nil {
		panic(err)
	}
	mask.DrawShape(ti.ShapeTypeCircle, 0.5)

	_, err = render.NewTextSprite(stage, "<text font-size='50'>Hello</text>\n <text font-size='60' color='#ff0000'>World</text>!", 720, 1280, ti.Align{HAlign: ti.HAlignCenter, VAlign: ti.VAlignMiddle})
	if err != nil {
		panic(err)
	}

	stage.Render()
	fmt.Printf("render: %v\n", time.Since(t))
	t = time.Now()

	err = ti.SaveTiImageToFile(stage.Texture(), filepath.Join(misc.GetCurrentDir(), "out.png"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("save: %v\n", time.Since(t))

}
