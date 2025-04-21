package main

import (
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	args := os.Args[1:]

	jsonfile := "/home/hayk/Documents/Literature/refs.bib"
	libdir := "/home/hayk/Documents/Literature/"
	if len(args) > 0 {
		jsonfile = args[0]
	}
	if len(args) > 1 {
		libdir = args[1]
	}
	app := NewApp(jsonfile, libdir)

	if err := wails.Run(&options.App{
		Title:  "guitest",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Frameless:        true,
		BackgroundColour: &options.RGBA{R: 23, G: 23, B: 23, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	}); err != nil {
		println("Error:", err.Error())
	}
}
