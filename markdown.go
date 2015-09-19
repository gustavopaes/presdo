package presdo

import (
  "strings"
  "html/template"
  "github.com/russross/blackfriday"
  "io/ioutil"
  "log"
  "path/filepath"
  "os"
  "time"
)

type MarkdownStruct struct { }

var markdown MarkdownStruct = MarkdownStruct{}

func (md *MarkdownStruct) Parse(requestFilePath string, markdownPath string) template.HTML {
    page, content := markdown.PageInfo(markdownPath)
    page.Content = render([]byte(content))

    pageHtml := page.HTML()

    save(requestFilePath, markdownPath, pageHtml)

    return pageHtml
}

func (md *MarkdownStruct) PageInfo(markdownPath string) (Page, string) {
    markdownContent, _ := ioutil.ReadFile(markdownPath)

    var pageContent Page = Page{
        Layout: "post",
    }

    content := extractHeader(markdownContent, &pageContent)

    return pageContent, content
}

// configure markdown render options
// See blackfriday markdown source for details
func render(content []byte) template.HTML {
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

func save(requestFilePath string, markdownPath string, parsedMarkdown template.HTML) {
    cachedPath := paths.Cache(requestFilePath)
    cachedFolder := filepath.Dir(cachedPath)

    // check if full path file exists
    if _, err := os.Stat(cachedFolder); os.IsNotExist(err) {
        // create full path
        os.MkdirAll(cachedFolder, 0755)
    }

    if err := ioutil.WriteFile(cachedPath, []byte(parsedMarkdown), 0664); err != nil {
        log.Println("Error on create markdown cache file: ", err)
    } else {
        if markdownStat, err := os.Stat(markdownPath); err == nil {
            os.Chtimes(cachedPath, markdownStat.ModTime(), markdownStat.ModTime())
        }
    }
}

func extractHeader(markdownContent []byte, page *Page) string {
    page.Params = make(map[string]string)

    var lines = strings.Split(string(markdownContent), "\n")
    var found = 0

    for i, line := range lines {
        line = strings.TrimSpace(line)

        if found == 1 {
            // parse line for param
            colonIndex := strings.Index(line, ":")
            if colonIndex > 0 {
                key := strings.TrimSpace(line[:colonIndex])
                value := strings.TrimSpace(line[colonIndex+1:])
                value = strings.Trim(value, "\"") //remove quotes
                switch key {
                    case "index":
                    page.Index = false

                    if value == "true" {
                        page.Index = true
                    }
                    case "title":
                    page.Title = value
                    case "layout":
                    page.Layout = value
                    case "date":
                    page.Date, _ = time.Parse(websiteConfig.DateFormat, value)
                    default:
                    page.Params[key] = value
                }
            }
        } else if found >= 2 {
            // params over
            lines = lines[i:]
            break
        }

        if line == "---" {
            found += 1
        }
    }

    // convert markdown content
    //content := strings.Join(lines, "\n")
    return strings.Join(lines, "\n")
}
