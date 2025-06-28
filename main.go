package main

import (
	"encoding/json"
	filltemplates "gen-color-scheme/fill_templates"
	gen_color_scheme "gen-color-scheme/get_colors"
	"os"
	"strings"
)

func main() {

	path := filltemplates.GetWaybarPath()

	oldWall, err := os.ReadFile(path + "wallpapers.json")

	images := gen_color_scheme.GetJsonFromWaytrogen()

	currWall, _ := json.Marshal(images)

	if err == nil && string(oldWall) == string(currWall) {
		os.Exit(0)
	}

	colorDefinitions := make([]string, 0, len(*images))
	templates := make([]string, len(*images))

	for _, image := range *images {
		wallpaperColors := filltemplates.WallpaperColors{WallpaperAndMonitor: image, AvgColor: gen_color_scheme.GetImageColor(image)}
		colorDefinitions = append(colorDefinitions, wallpaperColors.GetColorStr())
		templates = append(templates, wallpaperColors.FillInTemplate())
	}

	defines := strings.Join(colorDefinitions, "\n")
	filledTemplates := strings.Join(templates, "\n")

	baseTemplate, err := os.ReadFile("templates/base.css")

	if err != nil {
		panic("could not load base template")
	}

	filledComponents := []string{string(baseTemplate), filledTemplates}

	styleSheet := strings.Join(filledComponents, "\n\n")

	err = os.WriteFile(path+"colors.css", []byte(defines), 0644)

	if err != nil {
		panic("could not write colors.css")
	}

	err = os.WriteFile(path+"style.css", []byte(styleSheet), 0644)

	if err != nil {
		panic("could not write style.css")
	}

	err = os.WriteFile(path+"wallpapers.json", currWall, 0644)

	if err != nil {
		panic("could not write wallpapers.json")
	}
}
