package main

import (
	"fmt"
	"image/color"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/PWND0U/dztool/dzAi/dzOcr"
	"github.com/up-zero/gotool/imageutil"
)

func main() {
	fmt.Println("NumCPU:", runtime.NumCPU())
	// 按实际情况配置下述路径
	config := dzOcr.Config{
		OnnxRuntimeLibPath:    "G:\\code\\go\\dztool\\dzAi\\dzOcr/dist/cpuLib/onnxruntime.dll",
		GPUOnnxRuntimeLibPath: "G:\\code\\go\\dztool\\dzAi\\dzOcr/dist/gpuLib/onnxruntime.dll",
		DetModelPath:          "G:\\code\\go\\dztool\\dzAi\\dzOcr/dist/model/det.onnx",
		RecModelPath:          "G:\\code\\go\\dztool\\dzAi\\dzOcr/dist/model/recV5.onnx",
		DictPath:              "G:\\code\\go\\dztool\\dzAi\\dzOcr/dist/model/dict.txt",
		DetOutsideExpandPix:   10,
		UseCuda:               true,
		//UseTensorrt:           true,
		GPUDeviceID: 0,
		//UseDirectML:           true,
		NumThreads: int(float64(runtime.NumCPU()) * 0.8),
	}

	// 初始化引擎
	var engine dzOcr.Engine
	engine, err := dzOcr.NewDzPaddleOcrEngine(config)
	if err != nil {
		log.Fatalf("创建 OCR 引擎失败: %v\n", err)
	}
	defer engine.Destroy()

	// 打开图像
	imagePath := "G:\\code\\go\\dztool\\dzAi\\dzOcr/dist/img_1.png"
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
	drawImg := dzOcr.DrawBoxes(img, boxes, &color.RGBA{
		R: 237,
		G: 81,
		B: 38,
		A: 255,
	})
	// 保存绘制后的图像
	imageutil.Save("G:\\code\\go\\dztool\\dzAi\\dzOcr\\main/test_box.png", drawImg, 100)
	//fmt.Println("-----------------")
	//results, err = engine.RunOCRByFile("G:\\code\\go\\dztool\\dzAi\\dzOcr/dist/test2.jpg")
	//if err != nil {
	//	log.Fatalf("运行 OCR 失败: %v\n", err)
	//}
	//for _, result := range results {
	//	log.Printf("识别结果: %v\n", result)
	//}
	//fmt.Println("-----------------")
	//results, err = engine.RunOCRByFile("G:\\code\\go\\dztool\\dzAi\\dzOcr/dist/test3.jpg")
	//if err != nil {
	//	log.Fatalf("运行 OCR 失败: %v\n", err)
	//}
	//for _, result := range results {
	//	log.Printf("识别结果: %v\n", result)
	//}
	//fmt.Println("-----------------")
	//results, err = engine.RunOCRByFile("G:\\code\\go\\dztool\\dzAi\\dzOcr/dist/test4.jpg")
	//if err != nil {
	//	log.Fatalf("运行 OCR 失败: %v\n", err)
	//}
	//for _, result := range results {
	//	log.Printf("识别结果: %v\n", result)
	//}
	//fmt.Println("-----------------")
	//results, err = engine.RunOCRByFile("G:\\code\\go\\dztool\\dzAi\\dzOcr/dist/test5.jpg")
	//if err != nil {
	//	log.Fatalf("运行 OCR 失败: %v\n", err)
	//}
	//for _, result := range results {
	//	log.Printf("识别结果: %v\n", result)
	//}
	now = time.Now()
	imagePaths := []string{
		"G:\\code\\go\\dztool\\dzAi\\dzOcr/dist/img.png",
		"G:\\code\\go\\dztool\\dzAi\\dzOcr/dist/test.jpg",
	}
	wg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			nowi := time.Now()
			_, err := engine.RunOCRByFile(imagePaths[i%2])
			fmt.Printf("第%d次OCR识别耗时: %v\n", i, time.Since(nowi))
			if err != nil {
				log.Fatalf("运行 OCR 失败: %v\n", err)
			}
		}()
	}
	wg.Wait()
	elapsed = time.Since(now)
	fmt.Printf("100次OCR识别耗时: %v\n", elapsed)
}
