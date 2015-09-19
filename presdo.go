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
    "path/filepath"
    "io/ioutil"
    "encoding/json"
    "time"
    "net/http"
    "html/template"
    "sort"
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

// Page struct.
type Page struct {
    Title, Category, Layout, Url string
    Index                        bool
    Content                      template.HTML
    Date                         time.Time
    Updated                      time.Time
    Params                       map[string]string
}

// Index struct.
type Index struct {
    Title, Category, Layout, Url string
    Params map[string]string
    Pages  []Page
}

// Page Methods
func (p *Page) HTML() template.HTML {
    return view.HTML(p.Layout, p)
}

// Index Methods
func (index *Index) HTML() template.HTML {
    return view.HTML(index.Layout, index)
}

// Sort Index
type IndexByDate []Page

func (a IndexByDate) Len() int           { return len(a) }
func (a IndexByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a IndexByDate) Less(i, j int) bool { return a[i].Date.Unix() > a[j].Date.Unix() }

func (index *Index) Sort() {
    sort.Sort(IndexByDate(index.Pages))
}

// Basic website information
var websiteConfig struct {
    Domain       string
    PortNumber   int
    PublicDir    string
    TemplateDir  string
    DateFormat   string
    Ext          string
    Encode       string
    NotFound     string
    IsProduction bool
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

// Recreate cached file when markdown was changed
func recreateCachedFileIfChanged(requestFilePath string, markdownPath string, cachedPath string) {
    markdownStat, err1 := os.Stat(markdownPath)
    cachedStat, err2   := os.Stat(cachedPath)

    if err1 == nil && err2 == nil {
        if cachedStat.ModTime() != markdownStat.ModTime() {
            markdown.Parse(requestFilePath, markdownPath)
        }
    }
}

// Look for an markdown file
func responseMarkdownFile(w http.ResponseWriter, r *http.Request) error {
    requestFilePath := paths.Request(r.RequestURI)
    cachedPath      := paths.Cache(requestFilePath)
    markdownPath    := paths.Markdown(requestFilePath)

    // if markdown not exist
    if _, err := os.Stat(markdownPath); err != nil {
        return err
    }

    cachedStat, errStat := os.Stat(cachedPath)

    // Development mode always parse markdown file.
    // if not exist cached file, parse markdown content and create HTML file synchronous.
    if websiteConfig.IsProduction == false || errStat != nil {
        markdown.Parse(requestFilePath, markdownPath)
        cachedStat, _ = os.Stat(cachedPath)
    }

    server.Content(w, r, cachedPath, cachedStat.ModTime())

    go recreateCachedFileIfChanged(requestFilePath, markdownPath, cachedPath)

    return nil
}

func readDirListAndAppend(dir string) []string {
  var files []string

  dirlist, _ := ioutil.ReadDir(dir)
  for _, fi := range dirlist {
    f := filepath.Join(dir, fi.Name())
    ext := filepath.Ext(f)

    if ext == ".html" || ext == ".md" {
      files = append(files, f)
    } else {
      // recursively
      files = append(files, readDirListAndAppend(f)...)
    }
  }

  return files
}

// Look for index config file
func responseIndexFile(w http.ResponseWriter, r *http.Request) error {
    // check if is index path
    files := readDirListAndAppend( paths.IndexPath(r.RequestURI) )

    index := Index{
        Layout: "category",
    }

    for _, markdownPath := range files {
        page, _ := markdown.PageInfo(markdownPath)
        page.Url = paths.Page(markdownPath)
        index.Pages = append(index.Pages, page)
    }

    index.Sort()

    w.Write([]byte( index.HTML() ))

    return nil
}

