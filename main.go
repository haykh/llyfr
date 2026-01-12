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
	pdfviewer := "zathura"

	if len(args) > 0 {
		pdfviewer = args[0]
	}

	bibfile := "refs.bib"
	libdir := "./"
	if len(args) > 1 {
		bibfile = args[1]
	}
	if len(args) > 2 {
		libdir = args[2]
	} else if len(args) > 0 {
		libdir = filepath.Dir(bibfile)
	}
	app := NewApp(pdfviewer, bibfile, libdir)

	if err := wails.Run(&options.App{
		Title:  "llyfr",
		Width:  1280,
		Height: 960,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Frameless:        true,
		BackgroundColour: &options.RGBA{R: 23, G: 23, B: 23, A: 1},
		OnStartup:        app.startup,
		Bind: []any{
			app,
		},
	}); err != nil {
		println("Error:", err.Error())
	}
}
