package presdo

import (
  "fmt"
  "time"
  "net/http"
)

// Log all request
func LogRequest(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("%s %s [%s] %s %s\n", "->", time.Now(), r.RemoteAddr, r.Method, r.RequestURI)
}

// Log all response
func LogResponse(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("%s %s [%s] %s %s\n", "<-", time.Now(), r.RemoteAddr, r.Method, r.RequestURI)
}

func LogAction(message string) {
  fmt.Printf("   %s %s\n", time.Now(), message)
}