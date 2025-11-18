package api

import (
    "encoding/json"
    "log"
    "net/http"
    "strings"

    "github.com/username/url-lengthener/backend/internal/shortener"
    "github.com/username/url-lengthener/backend/internal/storage"
)

type Handler struct {
    store storage.Store
    s     *shortener.Shortener
}

func NewHandler(store storage.Store, s *shortener.Shortener) *Handler {
    return &Handler{store: store, s: s}
}

type lengthenRequest struct {
    URL string `json:"url"`
}

type lengthenResponse struct {
    Lengthened string `json:"lengthened"`
}

// set simple CORS headers for local dev
func setCORS(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// Lengthen accepts JSON { url: "https://..." } and returns a much longer URL
func (h *Handler) Lengthen(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodOptions {
        setCORS(w)
        w.WriteHeader(http.StatusOK)
        return
    }

    setCORS(w)
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req lengthenRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }
    req.URL = strings.TrimSpace(req.URL)
    if req.URL == "" {
        http.Error(w, "url required", http.StatusBadRequest)
        return
    }

    slug, err := h.s.Generate()
    if err != nil {
        log.Printf("generate slug: %v", err)
        http.Error(w, "internal", http.StatusInternalServerError)
        return
    }

    h.store.Save(slug, req.URL)

    scheme := "http"
    if r.TLS != nil {
        scheme = "https"
    }
    host := r.Host
    full := scheme + "://" + host + "/" + slug

    resp := lengthenResponse{Lengthened: full}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

// Redirect looks up the path as slug and redirects to target if found.
// If request path is root or path is /api/* it returns 404.
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
    setCORS(w)
    p := strings.TrimPrefix(r.URL.Path, "/")
    if p == "" || strings.HasPrefix(p, "api/") {
        http.NotFound(w, r)
        return
    }

    target, ok := h.store.Get(p)
    if !ok {
        http.NotFound(w, r)
        return
    }
    http.Redirect(w, r, target, http.StatusFound)
}
