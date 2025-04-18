package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
	"thumbnailinago/pkg/config"
	"thumbnailinago/pkg/locales"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
)

// App struct
type App struct {
	ctx context.Context
}

// stores all the settings
type SettingsStruct struct {
	SettingsYAML
	SVG string
	Templates
}

type SettingsYAML struct {
	Frontend FrontendSettings `yaml:"frontend"`
	Paths    PathSettings     `yaml:"paths"`
}

// settings of the frontend
type FrontendSettings struct {
	Locale         string   `json:"locale" yaml:"locale"`
	Days           []string `json:"days" yaml:"days"`
	ReplacementKey string   `json:"replacementKey" yaml:"replacementKey"`
	DateFormat     string   `json:"dateFormat" yaml:"dateFormat"`
}

// last paths of the File-dialogues
type PathSettings struct {
	SVG      string `yaml:"svg"`
	Export   string `yaml:"export"`
	Inkscape string `yaml:"inkscape"`
}

type Templates struct {
	DateTemplate *template.Template
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) beforeClose(ctx context.Context) bool {
	Settings.Save()

	return false
}

var Settings SettingsStruct

// sends the settings to the frontend
func (a *App) GetSettings() FrontendSettings {
	return Settings.SettingsYAML.Frontend
}

// stores the settings
func (a *App) SetSettings(newSettings FrontendSettings) error {
	runtime.LogDebugf(a.ctx, "%v:", Settings.SettingsYAML)

	Settings.SettingsYAML.Frontend = newSettings

	// recreate the templates because there may be a new date-format string
	if err := Settings.RecreateTemplates(); err != nil {
		return err
	}

	runtime.LogDebugf(a.ctx, "%v:", Settings.SettingsYAML)

	return nil
}

// recreates the templates
func (s *SettingsStruct) RecreateTemplates() error {
	if _, err := s.Templates.DateTemplate.Parse(s.Frontend.DateFormat); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *SettingsStruct) Save() error {
	// check wether the parent-directory exists
	if configPath, err := config.GetConfigPath(); err != nil {
		return err
	} else {
		parent := filepath.Dir(configPath)

		if _, err := os.Stat(parent); err != nil {
			os.MkdirAll(parent, 0777)
		}

		if data, err := yaml.Marshal(s.SettingsYAML); err != nil {
			return err
		} else {
			if err := os.WriteFile(configPath, data, 0777); err != nil {
				return err
			}
		}

		return nil
	}
}

type FrontendTemplate struct {
	SVG  string
	Name string `json:"name"`
}

// opens and loads a template-file
func (a *App) OpenTemplate() (FrontendTemplate, error) {
	// check, wether the directory of the last svg-file exists
	dir := filepath.Dir(Settings.Paths.SVG)

	if _, err := os.Stat(dir); err != nil {
		dir = ""
	}

	if file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Open Template",
		DefaultDirectory: dir,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "SVG-File (*.svg)",
				Pattern:     "*.svg",
			},
		},
	}); err != nil {
		return FrontendTemplate{}, err

		// check for an empty-file value (=file-selection cancelled)
	} else if len(file) == 0 {
		return FrontendTemplate{}, nil
	} else {
		// read the file-content
		if content, err := os.ReadFile(file); err != nil {
			return FrontendTemplate{}, err
		} else {
			// store the template
			Settings.SVG = string(content)

			// store the path
			Settings.Paths.SVG = file

			return a.RefreshPreview()
		}
	}
}

func (a *App) RefreshPreview() (FrontendTemplate, error) {
	var buf bytes.Buffer
	now := time.Now()

	if err := Settings.Templates.DateTemplate.Execute(&buf, locales.Create(Settings.Frontend.Locale, &now)); err != nil {
		return FrontendTemplate{}, err
	}

	preview := strings.ReplaceAll(Settings.SVG, Settings.Frontend.ReplacementKey, buf.String())

	return FrontendTemplate{SVG: preview, Name: filepath.Base(Settings.Paths.SVG)}, nil
}

type GenerateThumbnailsJob struct {
	From        string   `json:"from"`
	To          string   `json:"to"`
	Time        string   `json:"time"`
	CustomDates []string `json:"customDates"`
}

func (a *App) GenerateThumbnails(job GenerateThumbnailsJob) (int, error) {
	exportCount := 0

	// get the weekday of the start-date
	if startDate, err := time.Parse(time.DateOnly, job.From); err != nil {
		return exportCount, err
	} else if endDate, err := time.Parse(time.DateOnly, job.To); err != nil {
		return exportCount, err
	} else {
		dir := Settings.SettingsYAML.Paths.Export

		if _, err := os.Stat(dir); err != nil {
			dir = ""
		}

		// ask for the output-directory
		if outDir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
			DefaultDirectory:     dir,
			Title:                "Generate Thumbnails",
			CanCreateDirectories: true,
		}); err != nil {
			return exportCount, err

			// empty dir-string == dialog cancelled)
		} else if len(outDir) == 0 {
			return exportCount, nil
		} else {
			// store the path
			Settings.Paths.Export = outDir

			// sort the days-slice
			slices.Sort(Settings.Frontend.Days)

			// format the time for the output
			job.Time = strings.ReplaceAll(job.Time, ":", "-")

			currentDate := startDate

			// iterate through the days until the endDate is reached
			for currentDate.Compare(endDate) <= 0 {
				// check wether the current-weekday is included in days
				if slices.Contains(Settings.Frontend.Days, currentDate.Weekday().String()) {
					data := ThumbnailData{
						Date:   currentDate,
						Time:   &job.Time,
						SVG:    &Settings.SVG,
						OutDir: &outDir,
					}

					if err := data.generate(); err != nil {
						return exportCount, err
					} else {
						exportCount++
					}
				}

				currentDate = currentDate.AddDate(0, 0, 1)
			}

			// iterate through the custom-dates
			for _, customDate := range job.CustomDates {
				// try to parse the date
				if currentDate, err := time.Parse(time.DateOnly, customDate); err != nil {
					return exportCount, err
				} else {
					data := ThumbnailData{
						Date:   currentDate,
						Time:   &job.Time,
						SVG:    &Settings.SVG,
						OutDir: &outDir,
					}

					if err := data.generate(); err != nil {
						return exportCount, err
					} else {
						exportCount++
					}
				}
			}
		}
	}

	return exportCount, nil
}

type ThumbnailData struct {
	Date   time.Time
	Time   *string
	SVG    *string
	OutDir *string
}

// generate thumbnail
func (d *ThumbnailData) generate() error {
	filename := filepath.Join(*d.OutDir, fmt.Sprintf("%s.%s.png", d.Date.Format(time.DateOnly), *d.Time))

	dateMap := locales.Create(Settings.Frontend.Locale, &d.Date)

	var date bytes.Buffer
	if err := Settings.Templates.DateTemplate.Execute(&date, dateMap); err != nil {
		return err

		// create a copy of the template with the date inserted
	} else if file, err := os.CreateTemp(os.TempDir(), "*.svg"); err != nil {
		return err
	} else {
		defer os.Remove(file.Name())

		file.WriteString(strings.ReplaceAll(*d.SVG, Settings.Frontend.ReplacementKey, date.String()))

		inkscapeCommand := exec.Command(Settings.Paths.Inkscape, file.Name(), "-o", filename)

		if err := inkscapeCommand.Run(); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	// try to open the settings-file
	if configPath, err := config.GetConfigPath(); err != nil {
		panic(err)
	} else {
		// check wether the config-path exists
		if _, err := os.Stat(configPath); err != nil {
			// the config-file doesn´t exist -> create it with the default
			Settings = SettingsStruct{
				SettingsYAML: SettingsYAML{
					Frontend: FrontendSettings{
						Locale:         "en",
						Days:           []string{"sunday"},
						ReplacementKey: "SUNDAY_DATE",
						DateFormat:     "{{.Day}}. {{.Month}} {{.Year}}",
					},
					Paths: PathSettings{
						Inkscape: config.GetInkscapePath(),
					},
				},
			}

			// save the config
			Settings.Save()
		} else {
			// load the config
			if content, err := os.ReadFile(configPath); err != nil {
				panic(err)
			} else if err := yaml.Unmarshal(content, &Settings.SettingsYAML); err != nil {
				panic(err)
			}
		}

		// initialize the template
		Settings.Templates.DateTemplate = template.New("dateTemplate")
		Settings.RecreateTemplates()
	}
}
