package presdo

import (
    "html/template"
    "bytes"
    "log"
)

var ts *template.Template

type ViewStruct struct {
    
}

var view ViewStruct = ViewStruct{}

func (v *ViewStruct) Init() {
    funcMap := template.FuncMap{
        "isPost": func (layout string) bool {
          return layout == "post"
        },

        "isCategory": func (layout string) bool {
          return layout == "category"
        },
    }

    var err error

    ts = template.New("")
    LogAction("reading templates on " + websiteConfig.TemplateDir)
    if ts, err = ts.Funcs(funcMap).ParseGlob(websiteConfig.TemplateDir + "*.html"); err != nil {
        log.Fatalln("Template error:\n", err)
    }
}

func (v *ViewStruct) HTML(layout string, page *Page) template.HTML {
    buffer := new(bytes.Buffer)
    ts.ExecuteTemplate(buffer, layout + ".html", page)

    return template.HTML(buffer.String())
}

