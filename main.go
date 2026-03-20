package main

import (
	"path/filepath"
	"slideshow/misc"
	"slideshow/render"
	"slideshow/ti"
	"time"

	"github.com/go-mixed/go-taichi/taichi"
)

func main() {

	cd := misc.GetCurrentDir()
	runtime, err := taichi.NewRuntime(taichi.ArchCuda, filepath.Join(cd, "lib"))
	if err != nil {
		panic(err)
	}
	defer runtime.Release() // 必须释放

	renderer, err := render.NewRenderer(runtime)
	if err != nil {
		panic(err)
	}
	defer renderer.Release() // 必须释放

	// 创建舞台
	stage, err := render.NewStage(renderer, 1920, 1080)
	if err != nil {
		panic(err)
	}

	defer stage.Release()

	//
	img, err := render.NewImageSprite(renderer, filepath.Join(cd, "examples", "1.jpg"))
	if err != nil {
		panic(err)
	}
	defer img.Release()

	time.Sleep(2 * time.Second)

	//
	stage.Render()
	screen := stage.Texture()

	err = ti.SaveTiImageToFile(screen, filepath.Join(misc.GetCurrentDir(), "out.png"))
	if err != nil {
		panic(err)
	}
}
