package main

import (
	"bufio"
	"context"
	"errors"
	"os"
	"os/exec"
	"regexp"
	sysruntime "runtime"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	jsonfile    string
	libdir      string
	openingfile bool
}

// NewApp creates a new App application struct
func NewApp(jsonfile, libdir string) *App {
	return &App{
		jsonfile:    jsonfile,
		libdir:      libdir,
		openingfile: false,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
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
	var cmd *exec.Cmd
	filepath := a.libdir + filename

	switch sysruntime.GOOS {
	case "darwin":
		cmd = exec.Command("open", filepath)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", filepath)
	case "linux":
		cmd = exec.Command("xdg-open", filepath)
	default:
		return errors.New("unsupported platform")
	}

	if err := cmd.Start(); err != nil {
		return err
	} else {
		a.Exit()
		return nil
	}
}

func (a *App) GetLiterature() []map[string]string {
	if file, err := os.Open(a.jsonfile); err != nil {
		runtime.LogError(a.ctx, err.Error())
		return nil
	} else {
		defer file.Close()
		linebuffer := bufio.NewScanner(file)
		if err := linebuffer.Err(); err != nil {
			runtime.LogError(a.ctx, err.Error())
		}

		var elements []map[string]string

		for linebuffer.Scan() {
			line := linebuffer.Text()
			if newentry, err := regexp.MatchString("^@", line); err != nil {
				runtime.LogFatal(a.ctx, err.Error())
			} else if newentry {
				pubtype := regexp.MustCompile("^@([a-zA-Z0-9_]+)\\{").FindStringSubmatch(line)
				if strings.ToLower(pubtype[1]) != "comment" {
					elements = append(elements, map[string]string{
						"type": pubtype[1],
					})
				}
			} else {
				re := regexp.MustCompile("^\\s*([a-zA-Z0-9_]+)\\s*=\\s*{(.+)},$")
				if matches := re.FindStringSubmatch(line); len(matches) > 0 {
					if len(matches) == 3 {
						if len(elements) == 0 {
							runtime.LogFatal(a.ctx, "Elements not initialized")
							return nil
						}
						elements[len(elements)-1][matches[1]] = matches[2]
					} else {
						runtime.LogFatalf(a.ctx, "# matches: %d", len(matches))
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
				value = regexp.MustCompile("{|}|~").ReplaceAllString(value, "")
				value = regexp.MustCompile("\\n").ReplaceAllString(value, " ")
				value = regexp.MustCompile("\\s+").ReplaceAllString(value, " ")
				value = regexp.MustCompile(`\$\\gamma\$`).ReplaceAllString(value, "gamma")

				if k == "author" {
					authorlist := regexp.MustCompile(" and ").Split(value, -1)
					if strings.Contains(strings.ToLower(authorlist[0]), "collaboration") {
						authorlist = authorlist[:1]
					}
					// take at most 3
					if len(authorlist) > 3 {
						authorlist = authorlist[:3]
						authorlist = append(authorlist, "et al.")
					}
					// separate with commas
					value = ""
					for j := 0; j < len(authorlist); j++ {
						// remove after ,
						if strings.Contains(authorlist[j], ",") {
							authorlist[j] = strings.Split(authorlist[j], ",")[0]
						}
						value += authorlist[j]
						if j < len(authorlist)-1 {
							value += ", "
						}
					}
				} else if k == "journal" {
					if value == "Astrophysical Journal" || value == "The Astrophysical Journal" {
						value = "ApJ"
					} else if value == "Astrophysical Journal Letters" {
						value = "ApJL"
					} else if value == "Astrophysical Journal: Supplement" || value == "Astrophysical Journal Supplement" {
						value = "ApJS"
					} else if value == "Astronomy and Astrophysics" || value == "Astronomy & Astrophysics" || value == "Astronomy \\& Astrophysics" {
						value = "A&A"
					} else if value == "Monthly Notices of the Royal Astronomical Society" {
						value = "MNRAS"
					} else if value == "Physical Review D" {
						value = "PRD"
					} else if value == "Physical Review Letters" {
						value = "PRL"
					} else if value == "Journal of Plasma Physics" {
						value = "JPP"
					}
				} else if k == "file" {
					value = regexp.MustCompile(":(.+)?:").FindStringSubmatch(value)[1]
				}

				elements[i][k] = value
			}
		}
		for i := 0; i < len(elements); i++ {
			if strings.ToLower(elements[i]["type"]) != "article" {
				if strings.Contains(strings.ToLower(elements[i]["type"]), "thesis") {
					elements[i]["journal"] = "Thesis"
				} else {
					elements[i]["journal"] = strings.Title(strings.ToLower(elements[i]["type"]))
				}
			}
		}
		// remove empty elements
		for i := len(elements) - 1; i >= 0; i-- {
			if len(elements[i]) == 0 {
				elements = append(elements[:i], elements[i+1:]...)
			}
		}

		return elements
	}
}
