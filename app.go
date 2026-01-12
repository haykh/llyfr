package main

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type App struct {
	ctx         context.Context
	pdfviewer   string
	bibfile     string
	libdir      string
	openingfile bool
}

func NewApp(pdfviewer, bibfile, libdir string) *App {
	return &App{
		pdfviewer:   pdfviewer,
		bibfile:     bibfile,
		libdir:      libdir,
		openingfile: false,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.LogSetLogLevel(a.ctx, logger.WARNING)
}

func (a *App) Exit() {
	runtime.Quit(a.ctx)
}

func (a *App) OpenPDF(filename string) error {
	if a.openingfile {
		return nil
	}
	a.openingfile = true
	fpath := filepath.Join(a.libdir, filename)

	if err := exec.Command(a.pdfviewer, fpath).Start(); err != nil {
		return err
	} else {
		a.Exit()
		return nil
	}
}

func (a *App) GetLiterature() []map[string]string {
	if file, err := os.Open(a.bibfile); err != nil {
		runtime.LogError(a.ctx, err.Error())
		a.Exit()
		return nil
	} else {
		defer file.Close()
		linebuffer := bufio.NewScanner(file)
		if err := linebuffer.Err(); err != nil {
			a.Exit()
			runtime.LogError(a.ctx, err.Error())
		}

		var elements []map[string]string

		for linebuffer.Scan() {
			line := linebuffer.Text()
			if newentry, err := regexp.MatchString("^@", line); err != nil {
				runtime.LogFatal(a.ctx, err.Error())
				a.Exit()
			} else if newentry {
				pubtype := regexp.MustCompile(`^@([a-zA-Z0-9_]+)\{`).FindStringSubmatch(line)
				if strings.ToLower(pubtype[1]) != "comment" {
					elements = append(elements, map[string]string{
						"type": pubtype[1],
					})
				}
			} else {
				re := regexp.MustCompile(`^\s*([a-zA-Z0-9_]+)\s*=\s*{(.+)},$`)
				if matches := re.FindStringSubmatch(line); len(matches) > 0 {
					if len(matches) == 3 {
						if len(elements) == 0 {
							runtime.LogFatal(a.ctx, "Elements not initialized")
							a.Exit()
							return nil
						}
						elements[len(elements)-1][matches[1]] = matches[2]
					} else {
						runtime.LogFatalf(a.ctx, "# matches: %d", len(matches))
						a.Exit()
					}
				}
			}
		}
		for i := 0; i < len(elements); i++ {
			for k, v := range elements[i] {
				value := v
				value = regexp.MustCompile(`\{\\['`+"`"+`"^]\{?\\?([a-zA-Z])\}?\}`).ReplaceAllString(value, "$1")
				value = regexp.MustCompile(`\\c\{?c\}?`).ReplaceAllString(value, "c")
				value = regexp.MustCompile(`\\l`).ReplaceAllString(value, "l")
				value = regexp.MustCompile(`{|}|~`).ReplaceAllString(value, "")
				value = regexp.MustCompile(`\n`).ReplaceAllString(value, " ")
				value = regexp.MustCompile(`\s+`).ReplaceAllString(value, " ")
				value = regexp.MustCompile(`\$\\gamma\$`).ReplaceAllString(value, "gamma")

				switch k {
				case "author":
					authorlist := regexp.MustCompile(" and ").Split(value, -1)
					if strings.Contains(strings.ToLower(authorlist[0]), "collaboration") {
						authorlist = authorlist[:1]
					}
					// take at most 3
					if len(authorlist) > 3 {
						authorlist = authorlist[:3]
						authorlist = append(authorlist, "et al.")
					}
					value = ""
					for j := 0; j < len(authorlist); j++ {
						if strings.Contains(authorlist[j], ",") {
							authorlist[j] = strings.Split(authorlist[j], ",")[0]
						}
						value += authorlist[j]
						if j < len(authorlist)-1 {
							value += ", "
						}
					}
				case "journal":
					switch value {
					case "Astrophysical Journal", "The Astrophysical Journal":
						value = "ApJ"
					case "Astrophysical Journal Letters":
						value = "ApJL"
					case "Astrophysical Journal: Supplement", "Astrophysical Journal Supplement":
						value = "ApJS"
					case "Astronomy and Astrophysics", "Astronomy & Astrophysics", "Astronomy \\& Astrophysics":
						value = "A&A"
					case "Monthly Notices of the Royal Astronomical Society":
						value = "MNRAS"
					case "Physical Review D":
						value = "PRD"
					case "Physical Review Letters":
						value = "PRL"
					case "Journal of Plasma Physics":
						value = "JPP"
					}
				case "file":
					value = regexp.MustCompile(":(.+)?:").FindStringSubmatch(value)[1]
				}

				elements[i][k] = value
			}
		}
		caser := cases.Title(language.English)

		for i := 0; i < len(elements); i++ {
			if strings.ToLower(elements[i]["type"]) != "article" {
				if strings.Contains(strings.ToLower(elements[i]["type"]), "thesis") {
					elements[i]["journal"] = "Thesis"
				} else {
					elements[i]["journal"] = caser.String(strings.ToLower(elements[i]["type"]))
				}
			}
		}
		for i := len(elements) - 1; i >= 0; i-- {
			if len(elements[i]) == 0 {
				elements = append(elements[:i], elements[i+1:]...)
			}
		}

		return elements
	}
}
