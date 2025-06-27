package gen_color_scheme

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"math"
	"os"
	"os/exec"
	"strings"
)

func GetJsonFromWaytrogen() *[]WallpaperAndMonitor {
	cmd := exec.Command("waytrogen", "-l")
	stdOut, err := cmd.Output()

	if err != nil {
		panic("waytrogen is not installed")
	}

	jsonOut := []WallpaperAndMonitor{}
	json.Unmarshal(stdOut, &jsonOut)

	return &jsonOut
}

func GetImageColor(wAndM WallpaperAndMonitor) RGBA {
	cmd := exec.Command("mktemp", "-d")
	stdOut, _ := cmd.Output()
	tempFolder := strings.Trim(string(stdOut[:]), "\n")
	fileName := strings.Split(wAndM.Path, "/")
	tempFileName := fmt.Sprintf("%v/%v_temp.jpg", tempFolder, fileName[len(fileName)-1])
	cmd = exec.Command("ffmpeg", "-ss", "00:00:01.00", "-i", wAndM.Path, "-vf", "scale=100:100:force_original_aspect_ratio=decrease", "-vframes", "1", tempFileName)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(tempFileName)
		panic(err)
	}
	tempFile, err := os.Open(tempFileName)

	if err != nil {
		panic("could not open image")
	}

	defer tempFile.Close()

	imgData, err := jpeg.Decode(tempFile)

	if err != nil {
		panic("cannot decode image")
	}

	maxX := imgData.Bounds().Max.X
	maxY := imgData.Bounds().Max.Y
	var totalPixels = uint64(maxX * maxY)

	sumField := getSum(maxX, maxY, &imgData)
	sum := make([]chan uint64, 3)

	for i := range sum {
		sum[i] = make(chan uint64)
	}

	go sumField(sum[0], func(imgData *image.Image, x, y int) (val, alpha uint32) {
		r, _, _, a := (*imgData).At(x, y).RGBA()
		return r, a
	})

	go sumField(sum[1], func(imgData *image.Image, x, y int) (val, alpha uint32) {
		_, g, _, a := (*imgData).At(x, y).RGBA()
		return g, a
	})

	go sumField(sum[2], func(imgData *image.Image, x, y int) (val, alpha uint32) {
		_, _, b, a := (*imgData).At(x, y).RGBA()
		return b, a
	})

	return RGBA{r: <-sum[0] / totalPixels, g: <-sum[1] / totalPixels, b: <-sum[2] / totalPixels, a: 0.7}
}

func getSum(maxX, maxY int, imgData *image.Image) FetcherFn {
	return func(sumChan chan uint64, fn func(imageData *image.Image, x, y int) (val, alpha uint32)) {
		var sum uint64 = 0
		for x := 0; x <= maxX; x++ {
			for y := 0; y < maxY; y++ {
				val, alpha := fn(imgData, x, y)
				removeAlpha := getColor(alpha)
				sum += removeAlpha(val)
			}
		}
		sumChan <- sum
	}
}

func getColor(alpha uint32) func(uint32) uint64 {
	return func(color uint32) uint64 {
		return uint64(math.Floor((float64(color) / float64(alpha)) * 100))
	}
}
