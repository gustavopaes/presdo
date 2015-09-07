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

// Concat markdown directory with request path
func (p *PathStruct) Markdown(requestPath string) string {
    return path.Join("markdown", strings.Replace(requestPath, path.Ext(requestPath), ".md", 1))
}

func (p *PathStruct) Index(requestPath string) string {
    return path.Join("markdown", strings.Replace(requestPath, "index" + websiteConfig.Ext, "", 1), "presdo.index")
}

// Concat cache directory with request path
func (p *PathStruct) Cache(requestPath string) string {
    return path.Join("cache", requestPath)
}
