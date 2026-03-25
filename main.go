package main

import (
	"fmt"
	"path/filepath"
	"slideshow/misc"
	"slideshow/render"
	"slideshow/ti"
	"time"

	"github.com/go-mixed/go-taichi/taichi"
)

func main() {
	t := time.Now()
	cd := misc.GetCurrentDir()
	runtime, err := taichi.NewRuntime(taichi.ArchCuda, taichi.WithCacheTcm(true))
	if err != nil {
		panic(err)
	}
	defer runtime.Release() // 必须释放

	fmt.Printf("init taichi runtime: %v\n", time.Since(t))

	t = time.Now()
	renderer, err := render.NewRenderer(runtime)
	if err != nil {
		panic(err)
	}
	defer renderer.Release() // 必须释放

	fmt.Printf("init renderer: %v\n", time.Since(t))

	t = time.Now()
	// 创建舞台
	stage, err := render.NewStage(renderer, 720, 1280)
	if err != nil {
		panic(err)
	}

	defer stage.Release()

	fmt.Printf("init stage: %v\n", time.Since(t))

	t = time.Now()
	//
	img, err := render.NewImageSprite(renderer, filepath.Join(cd, "examples", "1.jpg"))
	if err != nil {
		panic(err)
	}
	defer img.Release()
	fmt.Printf("init image: %v\n", time.Since(t))

	t = time.Now()

	img.ResizeTo(720, 1280)

	fmt.Printf("resize image: %v\n", time.Since(t))

	stage.Add(img)

	t = time.Now()

	mask, err := render.NewShapeMask(renderer, 720, 1280, 720*0.5, 1280*0.5)
	if err != nil {
		panic(err)
	}
	defer mask.Release()

	mask.DrawShape(ti.ShapeTypeCircle, 0.5)

	img.SetMask(mask)

	fmt.Printf("init mask: %v\n", time.Since(t))

	t = time.Now()
	stage.Render()
	screen := stage.Texture()

	fmt.Printf("render: %v\n", time.Since(t))
	t = time.Now()

	err = ti.SaveTiImageToFile(screen, filepath.Join(misc.GetCurrentDir(), "out.png"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("save: %v\n", time.Since(t))

}
