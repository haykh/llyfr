package main

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	args := os.Args[1:]

	bibfile := "refs.bib"
	libdir := "./"
	if len(args) > 0 {
		bibfile = args[0]
	}
	if len(args) > 1 {
		libdir = args[1]
	} else if len(args) > 0 {
		libdir = filepath.Dir(bibfile)
	}
	app := NewApp(bibfile, libdir)

	if err := wails.Run(&options.App{
		Title:  "guitest",
		Width:  1280,
		Height: 960,
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
