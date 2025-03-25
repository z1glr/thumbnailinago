package main

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sqweek/dialog"
)

type responseMessage struct {
	Status  int
	Message *string
	Data    any
}

func (result responseMessage) send(c *fiber.Ctx) error {
	if result.Status >= 400 {
		if result.Message != nil {
			return fiber.NewError(result.Status, *result.Message)
		} else {
			return fiber.NewError(result.Status)
		}
	} else {
		if result.Data != nil {
			c.JSON(result.Data)
		} else {
			if result.Message != nil {
				c.SendString(*result.Message)
			}
		}

		return c.SendStatus(result.Status)
	}
}

type exportBody struct {
	Time   string   `json:"time"`
	From   string   `json:"from"`
	Until  string   `json:"until"`
	Custom []string `json:"custom"`
	Type   string   `json:"type"`
}

func postExport(c *fiber.Ctx) responseMessage {
	var response responseMessage

	var body exportBody

	if err := c.BodyParser(&body); err != nil {
		fmt.Printf("can't parse export-body: %v", err)
		response.Status = http.StatusBadRequest
	} else {
		if dir, err := dialog.Directory().Title("Export directory").SetStartDir(Config.LastExport).Browse(); err == nil {
			Config.LastExport = dir

			if dtFrom, err := time.Parse(time.DateTime, fmt.Sprintf("%s %s", body.From, body.Time)); err != nil {
				response.Status = http.StatusBadRequest
			} else if dtUntil, err := time.Parse(time.DateTime, fmt.Sprintf("%s %s", body.Until, body.Time)); err != nil {
				response.Status = http.StatusBadRequest
			} else {
				if exportDays, err := strucToMap(Config.ClientSettings.ExportDays); err != nil {
					response.Status = http.StatusInternalServerError
				} else {
					errors := map[time.Time]error{}

					// export the range
					for dtUntil.Sub(dtFrom) >= 0 {
						if exportDays[dtFrom.Weekday().String()] == true {
							if err := exportSvg(dtFrom, dir, body.Type); err != nil {
								errors[dtFrom] = err
							}
						}

						dtFrom = dtFrom.AddDate(0, 0, 1)
					}

					// export custom date
					for _, strDt := range body.Custom {
						if dt, err := time.Parse(time.DateTime, fmt.Sprintf("%s %s", strDt, body.Time)); err != nil {
							fmt.Printf("can't parse %q as date", strDt)
						} else if err := exportSvg(dt, dir, body.Type); err != nil {
							errors[dtFrom] = err
						}
					}

					if len(errors) > 0 {
						for dt, err := range errors {
							fmt.Printf("thumbnail creation failed for %q: %v", dt.Format(time.DateTime), err)
						}

						response.Status = http.StatusInternalServerError
					}
				}
			}
		}
	}

	return response
}

var dateFormatMap = map[string]string{
	"<d>":    "2",
	"<dd>":   "02",
	"<D>":    "Mon",
	"<DD>":   "Monday",
	"<m>":    "1",
	"<mm>":   "01",
	"<M>":    "Jan",
	"<MM>":   "January",
	"<yy>":   "06",
	"<yyyy>": "2006",
}

var germanMonths = map[time.Month]string{
	time.January:   "Januar",
	time.February:  "Februar",
	time.March:     "März",
	time.April:     "April",
	time.May:       "Mai",
	time.June:      "Juni",
	time.July:      "Juli",
	time.August:    "August",
	time.September: "September",
	time.October:   "Oktober",
	time.November:  "November",
	time.December:  "Dezember",
}

func formatDate(dt time.Time) string {
	dateString := Config.ClientSettings.DateFormat

	for k, v := range dateFormatMap {
		var formatted string
		switch k {
		case "<M>":
			formatted = germanMonths[dt.Local().Month()][:3]
		case "<MM>":
			formatted = germanMonths[dt.Local().Month()]
		default:
			formatted = dt.Local().Format(v)
		}

		dateString = strings.ReplaceAll(dateString, k, formatted)
	}

	return dateString
}

func exportSvg(dt time.Time, dir, ext string) error {
	if tempFile, err := os.CreateTemp(os.TempDir(), "thumbnailnago.*.svg"); err != nil {
		return fmt.Errorf("can't create tempfile: %v", err)
	} else {
		defer os.Remove(tempFile.Name())
		defer tempFile.Close()

		if err := os.WriteFile(tempFile.Name(), []byte(strings.ReplaceAll(svg.Str, Config.ClientSettings.ReplacementPattern, formatDate(dt))), 0o644); err != nil {
			return fmt.Errorf("can't write svg: %v", err)
		} else {
			filename := dt.Format(fmt.Sprintf("2006-01-02.15-04-05.%s", ext))
			var exportFile string
			if ext == "png" {
				exportFile = filename
			} else {
				exportFile = strings.ReplaceAll(tempFile.Name(), ".svg", ".png")
			}

			command := exec.Command("inkscape", tempFile.Name(), "-o", exportFile)

			if err := command.Run(); err != nil {
				return fmt.Errorf("can't run inkscape: %v", err)
			} else {
				if ext != "png" {
					defer os.Remove(exportFile)

					options := jpeg.Options{
						Quality: 95,
					}

					if fIn, err := os.Open(exportFile); err != nil {
						return fmt.Errorf("can't open exported png: %v", err)
					} else if fOut, err := os.Create(path.Join(dir, filename)); err != nil {
						return fmt.Errorf("can't create jpg: %v", err)
					} else if img, err := png.Decode(fIn); err != nil {
						return fmt.Errorf("can't decode png: %v", err)
					} else if err := jpeg.Encode(fOut, img, &options); err != nil {
						return fmt.Errorf("can't encode as jpg: %v", err)
					} else {
						defer fIn.Close()
						defer fOut.Close()

						return nil
					}
				} else {
					return nil
				}
			}
		}
	}
}

type ExportDays struct {
	Monday    bool `json:"Monday"`
	Tuesday   bool `json:"Tuesday"`
	Wednesday bool `json:"Wednesday"`
	Thursday  bool `json:"Thursday"`
	Friday    bool `json:"Friday"`
	Saturday  bool `json:"Saturday"`
	Sunday    bool `json:"Sunday"`
}

func getSettings(c *fiber.Ctx) responseMessage {
	var response responseMessage

	response.Data = Config.ClientSettings

	return response
}

func postSettings(c *fiber.Ctx) responseMessage {
	var response responseMessage

	var body ClientSettings

	if err := c.BodyParser(&body); err != nil {
		fmt.Printf("can't parse settings-body")
		response.Status = http.StatusBadRequest
	} else {
		Config.ClientSettings = body
	}

	return response
}

type svgFile struct {
	Name string `json:"name"`
	Str  string `json:"str"`
}

var svg = svgFile{
	Name: "default.svg",
	Str:  `<?xml version="1.0" encoding="utf-8" standalone="no"?><svg width="1920" height="1080" viewBox="0 0 507.99999 285.75001" version="1.1" id="svg8" inkscape:version="1.3.2 (091e20e, 2023-11-25, custom)" sodipodi:docname="default-template.svg" xmlns:inkscape="http://client.inkscape.org/namespaces/inkscape" xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd" xmlns:xlink="http://client.w3.org/1999/xlink" xmlns="http://client.w3.org/2000/svg" xmlns:svg="http://client.w3.org/2000/svg" xmlns:rdf="http://client.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:cc="http://creativecommons.org/ns#" xmlns:dc="http://purl.org/dc/elements/1.1/"><defs id="defs2"><pattern inkscape:collect="always" xlink:href="#Checkerboard" preserveAspectRatio="xMidYMid" id="pattern14" patternTransform="scale(15.9)" x="0" y="0" /><pattern inkscape:collect="always" style="fill:#000000" patternUnits="userSpaceOnUse" width="2" height="2" patternTransform="translate(0,0) scale(10,10)" id="Checkerboard" inkscape:stockid="Checkerboard" preserveAspectRatio="xMidYMid" inkscape:isstock="true" inkscape:label="Checkerboard"><rect style="stroke:none" x="0" y="0" width="1" height="1" id="rect209" /><rect style="stroke:none" x="1" y="1" width="1" height="1" id="rect211" /></pattern><filter inkscape:collect="always" style="color-interpolation-filters:sRGB" id="filter957" x="-0.011544351" width="1.0230887" y="-0.035689892" height="1.0713798"><feGaussianBlur inkscape:collect="always" stdDeviation="1.1219964" id="feGaussianBlur959" /></filter></defs><sodipodi:namedview id="base" pagecolor="#ffffff" bordercolor="#666666" borderopacity="1.0" inkscape:pageopacity="0.0" inkscape:pageshadow="2" inkscape:zoom="0.98994949" inkscape:cx="1802.1121" inkscape:cy="1120.2592" inkscape:document-units="px" inkscape:current-layer="layer1" inkscape:document-rotation="0" showgrid="false" units="px" inkscape:snap-page="true" inkscape:snap-bbox="true" inkscape:snap-bbox-midpoints="true" inkscape:bbox-nodes="true" inkscape:snap-global="true" inkscape:window-width="1918" inkscape:window-height="1008" inkscape:window-x="953" inkscape:window-y="1080" inkscape:window-maximized="0" inkscape:showpageshadow="0" inkscape:pagecheckerboard="1" inkscape:deskcolor="#d1d1d1" /><metadata id="metadata5"><rdf:RDF><cc:Work rdf:about=""><dc:format>image/svg+xml</dc:format><dc:type rdf:resource="http://purl.org/dc/dcmitype/StillImage" /></cc:Work></rdf:RDF></metadata><g inkscape:label="Layer 1" inkscape:groupmode="layer" id="layer1"><rect style="fill:#ff00ff;fill-rule:evenodd;stroke-width:16.9333;stroke-linecap:round;stroke-linejoin:round" id="rect14" width="508" height="285.75" x="0" y="0" /><rect style="fill:url(#pattern14);fill-rule:evenodd;stroke-width:16.9333;stroke-linecap:round;stroke-linejoin:round;fill-opacity:1" id="rect1" width="508" height="285.75" x="0" y="0" /><text xml:space="preserve" style="font-weight:bold;font-size:10.5833px;line-height:1.25;font-family:sans-serif;-inkscape-font-specification:'sans-serif Bold';stroke-width:0.264583" x="116.62633" y="-83.117203" id="text849"><tspan sodipodi:role="line" id="tspan847" x="116.62633" y="-83.117203" style="stroke-width:0.264583" /></text><text xml:space="preserve" style="font-weight:bold;font-size:35.2404px;line-height:1.25;font-family:'Noto Sans';-inkscape-font-specification:'Noto Sans Bold';text-align:center;text-anchor:middle;stroke-width:0.88101;filter:url(#filter957)" x="370.15988" y="41.800316" id="text853"><tspan sodipodi:role="line" id="tspan851" x="370.15988" y="41.800316" style="font-style:normal;font-variant:normal;font-weight:500;font-stretch:normal;font-family:'Noto Sans';-inkscape-font-specification:'Noto Sans Medium';text-align:center;text-anchor:middle;stroke-width:0.88101">Service on</tspan><tspan sodipodi:role="line" x="370.15988" y="85.850815" style="font-style:normal;font-variant:normal;font-weight:500;font-stretch:normal;font-family:'Noto Sans';-inkscape-font-specification:'Noto Sans Medium';text-align:center;text-anchor:middle;stroke-width:0.88101" id="tspan855">SUNDAY_DATE</tspan></text><text xml:space="preserve" style="font-weight:bold;font-size:35.2404px;line-height:1.25;font-family:'Noto Sans';-inkscape-font-specification:'Noto Sans Bold';text-align:center;text-anchor:middle;fill:#ffffff;fill-opacity:1;stroke-width:0.88101" x="368.04318" y="39.683647" id="text853-3"><tspan sodipodi:role="line" id="tspan851-6" x="368.04318" y="39.683647" style="font-style:normal;font-variant:normal;font-weight:500;font-stretch:normal;font-family:'Noto Sans';-inkscape-font-specification:'Noto Sans Medium';text-align:center;text-anchor:middle;fill:#ffffff;fill-opacity:1;stroke-width:0.88101">Service on</tspan><tspan sodipodi:role="line" x="368.04318" y="83.734146" style="font-style:normal;font-variant:normal;font-weight:500;font-stretch:normal;font-family:'Noto Sans';-inkscape-font-specification:'Noto Sans Medium';text-align:center;text-anchor:middle;fill:#ffffff;fill-opacity:1;stroke-width:0.88101" id="tspan855-7">SUNDAY_DATE</tspan></text></g></svg>`,
}

func getSvg(c *fiber.Ctx) responseMessage {
	var response responseMessage

	response.Data = svg

	return response
}

func patchSvg(c *fiber.Ctx) responseMessage {
	var response responseMessage

	if filename, err := dialog.File().Title("Open template").SetStartFile(Config.LastSvg).Filter("Scalable Vector Graphics (*.svg)", "svg").Load(); err != nil {
		fmt.Printf("can't open file: %v", err)
		response.Status = http.StatusInternalServerError
	} else {
		Config.LastSvg = filename

		if svgNew, err := os.ReadFile(filename); err != nil {
			fmt.Printf("can't read svg-file %q: %v", filename, err)
			response.Status = http.StatusInternalServerError
		} else {
			svg = svgFile{
				Str:  string(svgNew),
				Name: filename,
			}

			response = getSvg(c)
		}
	}

	return response
}

func exit(c *fiber.Ctx) responseMessage {
	writeConfig()

	// os.Exit(0)

	return responseMessage{}
}

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan

		writeConfig()

		os.Exit(0)
	}()

	app := fiber.New(fiber.Config{
		AppName:               "advent-server",
		DisableStartupMessage: true,
	})

	endpoints := map[string]map[string]func(*fiber.Ctx) responseMessage{
		"GET": {
			"settings": getSettings,
			"svg":      getSvg,
		},
		"POST": {
			"settings": postSettings,
			"export":   postExport,
			"exit":     exit,
		},
		"PATCH": {
			"svg": patchSvg,
		},
	}

	handleMethods := map[string]func(path string, handlers ...func(*fiber.Ctx) error) fiber.Router{
		"GET":    app.Get,
		"POST":   app.Post,
		"PATCH":  app.Patch,
		"DELETE": app.Delete,
	}

	for method, handlers := range endpoints {
		for address, handler := range handlers {
			handleMethods[method]("/api/"+address, func(c *fiber.Ctx) error {
				response := handler(c)

				return response.send(c)
			})
		}
	}

	app.Static("/", "./client")

	// launch the client-window
	url := "http://localhost:4217/index.html"
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	exec.Command(cmd, args...).Start()

	app.Listen(fmt.Sprintf("localhost:%d", 4217))
}
