package dzOcr

import (
	"log"
	"testing"

	"github.com/up-zero/gotool/imageutil"
)

func TestDzOcr(t *testing.T) {
	// 按实际情况配置下述路径
	config := Config{
		OnnxRuntimeLibPath:    "./dist/cpuLib/onnxruntime.dll",
		GPUOnnxRuntimeLibPath: "./dist/gpuLib/onnxruntime.dll",
		DetModelPath:          "./dist/model/det.onnx",
		RecModelPath:          "./dist/model/recV5.onnx",
		DictPath:              "./dist/model/dict.txt",
		UseCuda:               true,
	}

	// 初始化引擎
	var engine Engine
	engine, err := NewDzPaddleOcrEngine(config)
	if err != nil {
		log.Fatalf("创建 OCR 引擎失败: %v\n", err)
	}
	defer engine.Destroy()

	// 打开图像
	imagePath := "./dist/test4.jpg"
	img, err := imageutil.Open(imagePath)
	if err != nil {
		log.Fatalf("加载图像失败: %v\n", err)
	}

	// OCR识别
	results, err := engine.RunOCR(img)
	if err != nil {
		log.Fatalf("运行 OCR 失败: %v\n", err)
	}
	boxes := make([][4]int, 0)
	for _, result := range results {
		log.Printf("识别结果: %v\n", result)
		boxes = append(boxes, result.Box)
	}
	// 绘制识别区域
	drawImg := DrawBoxes(img, boxes, nil)
	// 保存绘制后的图像
	imageutil.Save("./dist/test_box.jpg", drawImg, 100)
}
