package filltemplates

import (
	"fmt"
	gen_color_scheme "gen-color-scheme/get_colors"
	"os"
	// "os/exec"
	"regexp"
)

type WallpaperColors struct {
	WallpaperAndMonitor gen_color_scheme.WallpaperAndMonitor
	AvgColor            gen_color_scheme.RGBA
}

func (wallColors *WallpaperColors) GetColorStr() string {
	avg := wallColors.AvgColor
	inverse := avg.GetInverse()
	monitor := wallColors.WallpaperAndMonitor.Monitor
	return fmt.Sprint(avg.PrintFormatted(false, monitor), "\n", inverse.PrintFormatted(true, monitor))
}

func (wallColors *WallpaperColors) FillInTemplate() string {
	// cmd := exec.Command("echo $PWD")
	// stdOut, _ := cmd.Output()
	// fmt.Println(string(stdOut))
	file, err := os.ReadFile("./templates/monitor_template.css")

	if err != nil {
		panic("could not find monitor template")
	}

	re := regexp.MustCompile(`__MONITOR__`)

	return string(re.ReplaceAll(file, []byte(wallColors.WallpaperAndMonitor.Monitor)))

}

func GetWaybarPath() string { return fmt.Sprintf(`/home/%v/.config/waybar/`, os.Getenv("USER")) }
