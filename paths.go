package presdo

import (
    "path"
    "path/filepath"
    "strings"
)

type PathStruct struct { }

var paths PathStruct = PathStruct{}

// Concat public directory with request path
func (p *PathStruct) Public(requestPath string) string {
    return path.Join(websiteConfig.PublicDir, requestPath)
}

// Return correct request path, addin index file in directory request
func (p *PathStruct) Request(requestPath string) string {
    if filepath.Ext(requestPath) == "" {
        return path.Join(requestPath, "index" + websiteConfig.Ext)
    }

    return requestPath
}

// Transform markdown path in public request
func (p *PathStruct) Page(markdownPath string) string {
    url := markdownPath
    url = strings.Replace(url, "markdown/", "", 1)
    url = strings.Replace(url, ".md", websiteConfig.Ext, 1)

    return "/" + url
}

// Concat markdown directory with request path
func (p *PathStruct) Markdown(requestPath string) string {
    return path.Join("markdown", strings.Replace(requestPath, path.Ext(requestPath), ".md", 1))
}

func (p *PathStruct) Index(requestPath string) string {
    return path.Join("markdown", strings.Replace(requestPath, "index" + websiteConfig.Ext, "", 1), "presdo.index")
}

func (p *PathStruct) IndexPath(requestPath string) string {
    if filepath.Ext(requestPath) == "" {
        return path.Join("markdown", requestPath)
    } else {
        return path.Join("markdown", path.Dir(requestPath))
    }
}
