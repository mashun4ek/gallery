package views

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/csrf"

	"github.com/mashun4ek/webdevcalhoun/gallery/context"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.New("").Funcs(template.FuncMap{"csrfField": func() (template.HTML, error) { return "", errors.New("csrfField is not implemented") }}).ParseFiles(files...)

	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

// Render is used to render the view with predefined layout
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
		// do nothing
	default:
		vd = Data{
			Yield: data,
		}
	}
	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})
	if err := tpl.ExecuteTemplate(&buf, v.Layout, data); err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong. If the problem persists, please email mashun4ek@gmail.com", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

// slice of strings layout files
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// "views/layouts/footer.gohtml",
// "views/layouts/navbar.gohtml",
// "views/layouts/bootstrap.gohtml",

// takes slice of string (filepath for templates) and prepends TemplateDir directory
// eg the input {"home"} => {"views/home"} if TemplateDir=="views/"
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

// appends extention to string
// Eg {"home"} => {"home.gohtml"}
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
