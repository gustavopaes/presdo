package presdo

import (
    "net/http"
    "path/filepath"
    "mime"
    "fmt"
    "time"
    "log"
    "os"
)

type ServerStruct struct { }

var server ServerStruct = ServerStruct{}

// Start http server
func (s *ServerStruct) Run() {
    http.HandleFunc("/", s.request)

    domainAndPort := fmt.Sprintf("%s:%d", websiteConfig.Domain, websiteConfig.PortNumber)
    
    fmt.Printf("   %s starting Presdo server on http://%s\n", time.Now(), domainAndPort)
    log.Fatal(http.ListenAndServe(domainAndPort, nil))
}

// Intercept all http requests
func (s *ServerStruct) request(w http.ResponseWriter, r *http.Request) {
    LogRequest(w, r)

    // add default headers
    w.Header().Add("Server", "Presdo " + VERSION)
    w.Header().Add("Vary", "Accept-Encoding")

    if fileStat, err := getRequestURIStaticFile( paths.Request(r.RequestURI) ); err == nil {
        s.Static(w, r, fileStat.(FileInfo))
    } else {
        // it is not an static file, look for markdown content
        err := responseMarkdownFile(w, r)

        if os.IsNotExist(err) {
            LogResponse(w, r)
            http.NotFound(w, r)
        }
    }
}

// Send to client some static file
func (s *ServerStruct) Static(w http.ResponseWriter, r *http.Request, fileStat FileInfo) {
    statiFilePath := paths.Public(fileStat.FullPath)

    s.Content(w, r, statiFilePath, fileStat.Stat.ModTime())
}

// Send content response
func (s *ServerStruct) Content(w http.ResponseWriter, r *http.Request, filePath string, modTime time.Time) {
    ext := filepath.Ext(filePath)

    // define headers to send to client
    w.Header().Add("Content-Type", mime.TypeByExtension(ext) + "; charset=" + websiteConfig.Encode)

    if fileContent, err := os.Open(filePath); err == nil {
        LogResponse(w, r)
        http.ServeContent(w, r, filePath, modTime, fileContent)
    }
}
