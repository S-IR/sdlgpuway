package main

import (
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Page struct {
	Path          string
	Title         string
	Description   string
	TemplatesPath []string
}

var Pages = []Page{
	Page{
		Path:          "/",
		Title:         "Learn SDL GPU",
		Description:   "Bla bla",
		TemplatesPath: []string{"index.html"},
	},
	Page{
		Path:          "/tutorials/sdl-gpu",
		Title:         "SDL GPU: Get Started",
		Description:   "Get started!",
		TemplatesPath: []string{"/tutorials/sdl-gpu/template.html", "/tutorials/sdl-gpu/sidebar.html", "/tutorials/sdl-gpu/get-started.html"},
	},
	Page{
		Path:          "/tutorials/sdl-gpu/get-started",
		Title:         "SDL GPU: Get Started",
		Description:   "Get started!",
		TemplatesPath: []string{"/tutorials/sdl-gpu/template.html", "/tutorials/sdl-gpu/sidebar.html", "/tutorials/sdl-gpu/get-started.html"},
	},
}

type Tutorial string

const (
	SDL_GPU_TUTORIAL = "SDL_GPU_TUTORIAL"
)

type ExistingLanguageForATutorial struct {
	Language string
	LogoPath string
	LogoAlt  string
	Tutorial string
}

var SDLGPUOdinPages = []Page{
	Page{
		Path:          "/tutorials/sdl-gpu/odin/prerequisites",
		Title:         "SDL GPU Way: Prerequisites",
		Description:   "Get started!",
		TemplatesPath: []string{"/tutorials/sdl-gpu/template.html", "/tutorials/sdl-gpu/sidebar.html", "/tutorials/sdl-gpu/odin/prerequisites/prerequisites.html"},
	},
	Page{
		Path:          "/tutorials/sdl-gpu/odin/chapter-1-hello-gpu",
		Title:         "SDL GPU Way: Chapter 1: Hello GPU (Odin)",
		Description:   "This chapter introduces you to basic concepts of modern GPU API's.",
		TemplatesPath: []string{"/tutorials/sdl-gpu/template.html", "/tutorials/sdl-gpu/sidebar.html", "/tutorials/sdl-gpu/odin/chapter-1-hello-gpu/index.html"},
	},
	Page{
		Path:          "/tutorials/sdl-gpu/odin/chapter-1-hello-gpu/1-1-hello_sdl.html",
		Title:         "SDL GPU Way: 1:1 Hello SDL",
		Description:   "Open a window with SDL3.",
		TemplatesPath: []string{"/tutorials/sdl-gpu/template.html", "/tutorials/sdl-gpu/sidebar.html", "/tutorials/sdl-gpu/odin/chapter-1-hello-gpu/1-1-hello_sdl.html"},
	},
}
var existingLanguages = []string{"odin"}

func main() {
	Pages = append(Pages, SDLGPUOdinPages...)
	// devMode := flag.Bool("dev", false, "Run development server")
	buildMode := flag.Bool("build", false, "Build static site")
	flag.Parse()

	Pages = append(Pages, SDLGPUOdinPages...)

	if *buildMode {
		generateStaticSite()
		log.Println("Static site generated in output/ directory")
	} else {
		startDevServer()
	}
}

func renderTemplate(w io.Writer, page Page) {
	tpl := template.New(filepath.Base(page.TemplatesPath[0])).Funcs(template.FuncMap{
		"hasPrefix": func(s, prefix string) bool {
			return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
		},
		"eq": func(a, b string) bool { return a == b },
	})

	for _, tmplPath := range page.TemplatesPath {
		fullPath := filepath.Join("templates", tmplPath)
		_, err := tpl.ParseFiles(fullPath)
		if err != nil {
			panic(err)
		}
	}
	type ChosenLanguage =  string
	const (
		C    ChosenLanguage = "C"
		Odin ChosenLanguage = "Odin"
	)
	chosenLanguage := Odin
	if strings.Contains(page.Path, "/odin/") {
		chosenLanguage = Odin
	} else if strings.Contains(page.Path, "/c/") {
		chosenLanguage = C
	} else {
		chosenLanguage = ""
	}

	data := struct {
		Page
		CurrentPage    string
		Title          string
		Description    string
		ChosenLanguage ChosenLanguage
	}{
		Page:        page,
		CurrentPage: page.Path,
		Title:       page.Title,
		Description: page.Description,
		ChosenLanguage: chosenLanguage,
	}

	err := tpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
func startDevServer() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestedPath := r.URL.Path

		var page *Page
		for _, p := range Pages {
			if p.Path == requestedPath {
				page = &p
				break
			}
		}

		if page == nil {
			http.NotFound(w, r)
			return
		}

		renderTemplate(w, *page)
	})

	log.Println("Development server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func generateStaticSite() {
	os.RemoveAll("public12")
	os.MkdirAll("public12", 0755)

	for _, page := range Pages {

		outputPath := filepath.Join("public12", page.Path)
		if filepath.Ext(outputPath) == "" {
			outputPath = filepath.Join(outputPath, "index.html")
		}

		os.MkdirAll(filepath.Dir(outputPath), 0755)

		file, err := os.Create(outputPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		renderTemplate(file, page)

		log.Printf("Generated: %s", outputPath)
	}
	CopyDir("./static/", "./output/static/")
}
