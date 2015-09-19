package presdo

import (
    "net/http"
    "path/filepath"
    "mime"
    "fmt"
    "time"
    "io/ioutil"
    "log"
    "os"
    "strings"
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

    pathRequest := paths.Request(r.RequestURI)

    // add default headers
    w.Header().Add("Server", "Presdo " + VERSION)
    w.Header().Add("Vary", "Accept-Encoding")

    if fileStat, err := getRequestURIStaticFile( pathRequest ); err == nil {
        s.Static(w, r, fileStat.(FileInfo))
    } else {
        // check extention
        if websiteConfig.Ext != filepath.Ext(pathRequest) {
            s.NotFound(w, r)
        } else {
            // it is not an static file, look for markdown content
            errMarkdown := responseMarkdownFile(w, r)

            if os.IsNotExist(errMarkdown) {
                s.NotFound(w, r)
            }
        }
    }
}

func (s *ServerStruct) NotFound(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/" + websiteConfig.NotFound, 302)
}

// Send to client some static file
func (s *ServerStruct) Static(w http.ResponseWriter, r *http.Request, fileStat FileInfo) {
    statiFilePath := paths.Public(fileStat.FullPath)
    http.ServeFile(w, r, statiFilePath)
}

var unixEpochTime = time.Unix(0, 0)

// Copyright 2009 The Go Authors. All rights reserved.
// http://golang.org/src/net/http/fs.go#L261
func checkLastModified(w http.ResponseWriter, r *http.Request, modtime time.Time) bool {
    if modtime.IsZero() || modtime.Equal(unixEpochTime) {
        // If the file doesn't have a modtime (IsZero), or the modtime
        // is obviously garbage (Unix time == 0), then ignore modtimes
        // and don't process the If-Modified-Since header.
        return false
    }

    // The Date-Modified header truncates sub-second precision, so
    // use mtime < t+1s instead of mtime <= t to check for unmodified.
    if t, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since")); err == nil && modtime.Before(t.Add(1*time.Second)) {
        h := w.Header()
        delete(h, "Content-Type")
        delete(h, "Content-Length")
        w.WriteHeader(http.StatusNotModified)
        return true
    }

    w.Header().Set("Last-Modified", modtime.UTC().Format(http.TimeFormat))
    return false
}

// Send content response
func (s *ServerStruct) Content(w http.ResponseWriter, r *http.Request, filePath string, modTime time.Time) {
    ext := filepath.Ext(filePath)

    // define headers to send to client
    w.Header().Add("Content-Type", mime.TypeByExtension(ext))

    if fileContent, err := ioutil.ReadFile(filePath); err == nil {
        if strings.Contains(filePath, websiteConfig.NotFound) {
            w.WriteHeader(http.StatusNotFound)
        }

        checkLastModified(w, r, modTime)

        w.Write([]byte(fileContent))

        LogResponse(w, r)
    }
}
