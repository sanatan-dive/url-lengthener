package main

import (
    "log"
    "net/http"
    "os"
    "strconv"

    "github.com/username/url-lengthener/backend/internal/api"
    "github.com/username/url-lengthener/backend/internal/storage"
    "github.com/username/url-lengthener/backend/internal/shortener"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // length for generated slugs (number of characters)
    defaultLen := 256
    if v := os.Getenv("DEFAULT_SLUG_LEN"); v != "" {
        if n, err := strconv.Atoi(v); err == nil {
            defaultLen = n
        }
    }

    store := storage.NewInMemoryStore()
    s := shortener.New(defaultLen)
    handler := api.NewHandler(store, s)

    mux := http.NewServeMux()
    mux.HandleFunc("/api/lengthen", handler.Lengthen)
    mux.HandleFunc("/", handler.Redirect)

    log.Printf("starting server on :%s (default slug len=%d)", port, defaultLen)
    if err := http.ListenAndServe("0.0.0.0:"+port, mux); err != nil {
        log.Fatal(err)
    }
}
