package main
import (
    "testing"
    "fmt"
    "io"
    "net/http"
    "net/url"
)


func TestHttpProxy(t *testing.T) {
    // username := "devtest"
    // password := "123123aaaA"
    // entry := fmt.Sprintf("http://%s:%s@dc.jp-pr.oxylabs.io:12000", username, password)
    entry := "http://127.0.0.1:4780"
    // fmt.Println(entry)
    proxyURL, err := url.Parse(entry)
    if err != nil {
        t.Error("Failed to parse proxy URL:", err)
        return
    }
    client := &http.Client{
        Transport: &http.Transport{
            Proxy: http.ProxyURL(proxyURL),
        },
    }
    resp, err := client.Get("https://www.google.com")
    if err != nil {
        t.Error("Failed to send HTTP request:", err)
        return
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        t.Error("Failed to read response body:", err)
        return
    }
    fmt.Println(string(body))
}