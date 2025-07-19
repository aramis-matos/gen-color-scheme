package gen_color_scheme

import (
	"fmt"
	"image"
)

type RGBA struct {
	r uint64
	g uint64
	b uint64
	a float32
}

type Changer struct {
	Swww     string `json:"swww"`
	Mpvpaper string `json:"mpvpaper"`
}

type WallpaperAndMonitor struct {
	Monitor string  `json:"monitor"`
	Path    string  `json:"path"`
	Changer Changer `json:"changer"`
}

type FetcherFn func(chan uint64, func(imgData *image.Image, x, y int) (val, alpha uint32))

func (color *RGBA) GetInverse() RGBA {
	return RGBA{r: 255 - color.r, g: 255 - color.g, b: 255 - color.b, a: 1}
}

func (color *RGBA) PrintFormatted(isBg bool, monitor string) string {
	bg := ""
	if isBg {
		bg = "bg"
	}
	return fmt.Sprintf(`@define-color %vcolor%v rgba(%v,%v,%v,%v);`, bg, monitor, color.r, color.g, color.b, color.a)
}
