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
        "isPost": func (page Page) bool {
          return page.Layout == "post"
        },

        "isCategory": func (page Page) bool {
          return page.Layout == "category"
        },
    }

    var err error

    ts = template.New("")
    if ts, err = ts.Funcs(funcMap).ParseGlob(websiteConfig.TemplateDir + "*.html"); err != nil {
        log.Fatalln("Template error:\n", err)
    }
}

func (v *ViewStruct) HTML(page *Page) template.HTML {
    buffer := new(bytes.Buffer)
    ts.ExecuteTemplate(buffer, page.Layout + ".html", page)

    return template.HTML(buffer.String())
}

