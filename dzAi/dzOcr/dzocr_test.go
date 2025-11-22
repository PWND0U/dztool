package dzOcr

import (
	"fmt"
	"image/color"
	"log"
	"runtime"
	"testing"
	"time"

	"github.com/up-zero/gotool/imageutil"
)

func TestDzOcr(t *testing.T) {
	fmt.Println("NumCPU:", runtime.NumCPU())
	// 按实际情况配置下述路径
	config := Config{
		OnnxRuntimeLibPath:    "./dist/cpuLib/onnxruntime.dll",
		GPUOnnxRuntimeLibPath: "./dist/gpuLib/onnxruntime.dll",
		DetModelPath:          "./dist/model/det.onnx",
		RecModelPath:          "./dist/model/recV5.onnx",
		DictPath:              "./dist/model/dict.txt",
		DetOutsideExpandPix:   10,
		UseCuda:               true,
		//UseTensorrt:           true,
		//UseDirectML:           true,
		NumThreads: runtime.NumCPU(),
	}

	// 初始化引擎
	var engine Engine
	engine, err := NewDzPaddleOcrEngine(config)
	if err != nil {
		log.Fatalf("创建 OCR 引擎失败: %v\n", err)
	}
	defer engine.Destroy()

	// 打开图像
	imagePath := "./dist/img.png"
	img, err := imageutil.Open(imagePath)
	if err != nil {
		log.Fatalf("加载图像失败: %v\n", err)
	}

	// OCR识别
	now := time.Now()
	results, err := engine.RunOCR(img)
	elapsed := time.Since(now)
	fmt.Printf("OCR识别耗时: %v\n", elapsed)
	if err != nil {
		log.Fatalf("运行 OCR 失败: %v\n", err)
	}
	boxes := make([][4]int, 0)
	for _, result := range results {
		log.Printf("识别结果: %v\n", result)
		boxes = append(boxes, result.Box)
	}
	// 绘制识别区域
	drawImg := DrawBoxes(img, boxes, &color.RGBA{
		R: 237,
		G: 81,
		B: 38,
		A: 255,
	})
	// 保存绘制后的图像
	imageutil.Save("./dist/test_box.png", drawImg, 100)
	fmt.Println("-----------------")
	results, err = engine.RunOCRByFile("./dist/test2.jpg")
	if err != nil {
		log.Fatalf("运行 OCR 失败: %v\n", err)
	}
	for _, result := range results {
		log.Printf("识别结果: %v\n", result)
	}
	fmt.Println("-----------------")
	results, err = engine.RunOCRByFile("./dist/test3.jpg")
	if err != nil {
		log.Fatalf("运行 OCR 失败: %v\n", err)
	}
	for _, result := range results {
		log.Printf("识别结果: %v\n", result)
	}
	fmt.Println("-----------------")
	results, err = engine.RunOCRByFile("./dist/test4.jpg")
	if err != nil {
		log.Fatalf("运行 OCR 失败: %v\n", err)
	}
	for _, result := range results {
		log.Printf("识别结果: %v\n", result)
	}
	fmt.Println("-----------------")
	results, err = engine.RunOCRByFile("./dist/test5.jpg")
	if err != nil {
		log.Fatalf("运行 OCR 失败: %v\n", err)
	}
	for _, result := range results {
		log.Printf("识别结果: %v\n", result)
	}
	now = time.Now()
	imagePaths := []string{
		"./dist/img.png",
		"./dist/test.jpg",
	}
	for i := 0; i < 20; i++ {
		_, err := engine.RunOCRByFile(imagePaths[i%2])
		if err != nil {
			log.Fatalf("运行 OCR 失败: %v\n", err)
		}
	}
	elapsed = time.Since(now)
	fmt.Printf("100次OCR识别耗时: %v\n", elapsed)
}
