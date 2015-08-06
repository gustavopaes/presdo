/*
 * Presdo. An static markdown web server.
 *
 * Original basead on Lanyon, by Marcus Kazmierczak (http://github.com/mkaz/lanyon)
 *
 * Author: Gustavo Paes <gustavo.paes@gmail.com>
 * Source: http://github.com/gustavopaees/
 */

package presdo

import (
    "log"
    "os"
    "io/ioutil"
    "encoding/json"
    "time"
    "net/http"
    "html/template"
//    "compress/gzip"

    "github.com/russross/blackfriday"
)

const VERSION = "0.0.1"

// File info struct with full relative path and os.Stat()
type FileInfo struct {
    FullPath string
    Stat     os.FileInfo
}

// Time, in days, for each content type expire. It's sent to cache-control http response.
type cacheTimes struct {
    Html       int
    Css        int
    Javascript int
    Image      int
    Index      int
}

type Page struct {
  Title, Category, Layout, Url string
  Content                      template.HTML
  Date                         time.Time
  Updated                      time.Time
  Params                       map[string]string
}

// Basic website information
var websiteConfig struct {
    Domain      string
    PortNumber  int
    PublicDir   string
    TemplateDir string
    DateFormat  string
    Ext         string
    Encode      string
    NotFound    string
    CacheConfig cacheTimes
}

// Read config file and start server
func Run() {
    if code, err := readConfigFile("./Presdofile.json"); code != 200 {
        log.Fatalln("Error", code, err)
    }

    if code, err := setDefaultValues(); code != 200 {
        log.Fatalln("Error", code, err)
    }

    view.Init()
    server.Run()
}

// Read file config Presdofile.json
func readConfigFile(configFile string) (int, error) {
    if _, err := os.Stat(configFile); os.IsNotExist(err) {
        //log.Fatalln("Config file not found. Please, create 'Presdofile.json' and try again.")
        return 404, err
    }

    file, err := ioutil.ReadFile(configFile)

    if err != nil {
        //log.Fatalln("Error on reading config file: ", err)
        return 403, err
    }

    if err := json.Unmarshal(file, &websiteConfig); err != nil {
        //log.Fatalln("Error on parsing json. Please, review your 'Presdofile.json': ", err)
        return 500, err
    }

    return 200, nil
}

// set default values for each website config.
func setDefaultValues() (int, error) {
    if websiteConfig.Domain == "" {
        websiteConfig.Domain = "localhost"
    }

    if websiteConfig.PortNumber == 0 {
        websiteConfig.PortNumber = 8080
    }

    if websiteConfig.PublicDir == "" {
        websiteConfig.PublicDir = "public/"
    }

    if websiteConfig.TemplateDir == "" {
        websiteConfig.TemplateDir = "view/default/"
    }

    if websiteConfig.DateFormat == "" {
        websiteConfig.DateFormat = "2006-01-02 03:04pm"
    }

    if websiteConfig.Ext == "" {
        websiteConfig.Ext = ".htm"
    }

    if websiteConfig.Encode == "" {
        websiteConfig.Encode = "utf-8"
    }

    if websiteConfig.NotFound == "" {
        websiteConfig.NotFound = "404.html"
    }

    //if _, err := os.Stat(websiteConfig.PublicDir); os.IsNotExist(err) {
    //    //log.Fatalln(err)
    //    return 404, err
    //}

    //if _, err := os.Stat(websiteConfig.TemplateDir); os.IsNotExist(err) {
    //    //log.Fatalln(err)
    //    return 404, err
    //}

    return 200, nil
}

// Check if request is to a static file or not.
// If request is a directory, check if there is an index file
func getRequestURIStaticFile(requestURI string) (interface{}, error) {
    publicPathRequest := paths.Public(paths.Request(requestURI))
    fileStat, err := os.Stat(publicPathRequest)

    if err != nil {
        return nil, err
    }

    return FileInfo{
        FullPath: requestURI,
        Stat    : fileStat,
    }, nil
}

func recreateCachedFileIfChanged(requestFilePath string, markdownPath string, cachedPath string) {
    markdownStat, err1 := os.Stat(markdownPath)
    cachedStat, err2   := os.Stat(cachedPath)

    if err1 == nil && err2 == nil {
        if cachedStat.ModTime() != markdownStat.ModTime() {
            parseMarkdownFile(requestFilePath, markdownPath)
        }
    }
}

// Look for an markdown file
func responseMarkdownFile(w http.ResponseWriter, r *http.Request) error {
    requestFilePath := paths.Request(r.RequestURI)
    cachedPath      := paths.Cache(requestFilePath)
    markdownPath    := paths.Markdown(requestFilePath)

    // define headers to send to client
    w.Header().Add("Content-Type", "text/html; charset=" + websiteConfig.Encode)

    // if markdown not exist
    if _, err := os.Stat(markdownPath); err != nil {
        return err
    }

    var cachedStat os.FileInfo
    var errStat error

    // if not exist cached file, parse markdown content and create HTML file synchronous
    if cachedStat, errStat = os.Stat(cachedPath); errStat != nil {
        parseMarkdownFile(requestFilePath, markdownPath)
    }

    server.Content(w, r, cachedPath, cachedStat.ModTime())

    go recreateCachedFileIfChanged(requestFilePath, markdownPath, cachedPath)

    //if cachedStat, err := os.Stat(cachedPath); err == nil {
    //    // send cached content to user request
    //    if cacheFileContent, err := os.Open(cachedPath); err == nil {
    //        server.Content(w, r, cachedPath, cachedStat.ModTime(), cacheFileContent)
    //
    //        // if original content was changed, parse markdown file (async)
    //        if cachedStat.ModTime() != markdownStat.ModTime() {
    //            go parseMarkdownFile(requestFilePath, markdownPath)
    //        }
    //    }
    //} else {
    //    pageContent = parseMarkdownFile(requestFilePath, markdownPath)
    //    fmt.Fprintf(w, string(pageContent))
    //}

    return nil
}

func parseMarkdownFile(requestFilePath string, markdownPath string) template.HTML {
    markdownContent, _ := ioutil.ReadFile(markdownPath)

    pageHtml := view.HTML(&Page{
        Layout:   "page",
        Title:    "Página de teste!",
        Content:  markdownRender(markdownContent),
        Date:     time.Now(),
        Updated:  time.Now(),
        Params:   make(map[string]string),
    })

    log.Println("pageHtml", pageHtml)

    saveParsedMarkdownFile(requestFilePath, markdownPath, pageHtml)

    return pageHtml
}

func saveParsedMarkdownFile(requestFilePath string, markdownPath string, parsedMarkdown template.HTML) {
    cachedPath := paths.Cache(requestFilePath)

    if err := ioutil.WriteFile(cachedPath, []byte(parsedMarkdown), 0644); err != nil {
        log.Println("Error on create markdown cache file: ", err)
    } else {
        if markdownStat, err := os.Stat(markdownPath); err == nil {
            os.Chtimes(cachedPath, markdownStat.ModTime(), markdownStat.ModTime())
        }
    }
}

// configure markdown render options
// See blackfriday markdown source for details
func markdownRender(content []byte) template.HTML {
  htmlFlags := 0
  //htmlFlags |= blackfriday.HTML_SKIP_SCRIPT
  htmlFlags |= blackfriday.HTML_USE_XHTML
  htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
  htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
  htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
  renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

  extensions := 0
  extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
  extensions |= blackfriday.EXTENSION_TABLES
  extensions |= blackfriday.EXTENSION_FENCED_CODE
  extensions |= blackfriday.EXTENSION_AUTOLINK
  extensions |= blackfriday.EXTENSION_STRIKETHROUGH
  extensions |= blackfriday.EXTENSION_SPACE_HEADERS

  return template.HTML(blackfriday.Markdown([]byte(content), renderer, extensions))
}
