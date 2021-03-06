package presdo

import (
  "strings"
  "html/template"
  "github.com/russross/blackfriday"
  "io/ioutil"
  "time"
)

type MarkdownStruct struct { }

var markdown MarkdownStruct = MarkdownStruct{}

func (md *MarkdownStruct) PageInfo(markdownPath string) Page {
    markdownContent, _ := ioutil.ReadFile(markdownPath)

    var pageInfo Page = Page{
        Layout: "post",
    }

    content := extractHeader(markdownContent, &pageInfo)

    pageInfo.Content = render([]byte(content))

    return pageInfo
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
                    case "indexSort":
                    page.IndexSort = strings.TrimSpace(value)
                    case "title":
                    page.Title = strings.TrimSpace(value)
                    case "layout":
                    page.Layout = strings.TrimSpace(value)
                    case "date":
                    page.Date, _ = time.Parse(websiteConfig.DateFormat, value)
                    default:
                    page.Params[key] = strings.TrimSpace(value)
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
