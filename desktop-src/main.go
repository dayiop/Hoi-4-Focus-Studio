package main

import (
	"bytes"
	"context"
	"embed"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

//go:embed web/*
var appFiles embed.FS

const (
	appName = "HOI4 Focus Studio"
	address = "127.0.0.1:38429"
	origin  = "http://" + address
	magic   = "HOI4-Focus-Studio-Desktop-2026"
)

var heartbeat atomic.Int64

const desktopHeartbeat = `<script>
(function(){
  function ping(){fetch('/_desktop_ping',{method:'POST',cache:'no-store'}).catch(function(){});}
  ping();setInterval(ping,8000);
  document.addEventListener('visibilitychange',function(){if(!document.hidden)ping();});
})();
</script>`

func baseDataDir() string {
	if p := os.Getenv("LOCALAPPDATA"); p != "" {
		return filepath.Join(p, "HOI4FocusStudio")
	}
	if p, err := os.UserConfigDir(); err == nil {
		return filepath.Join(p, "HOI4FocusStudio")
	}
	return filepath.Join(os.TempDir(), "HOI4FocusStudio")
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		if checkAlreadyRunning() {
			if e := launchBrowser(); e != nil {
				showError(e.Error())
			}
			return
		}
		showError("Не удалось запустить локальный сервер редактора. Порт 38429 занят другой программой.\n\n" + err.Error())
		return
	}
	defer listener.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/_desktop_health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Focus-Studio", magic)
		_, _ = io.WriteString(w, magic)
	})
	mux.HandleFunc("/_desktop_ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		heartbeat.Store(time.Now().Unix())
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/", serveApp)
	srv := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != address {
			http.Error(w, "Invalid host", http.StatusForbidden)
			return
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		mux.ServeHTTP(w, r)
	}), ReadHeaderTimeout: 5 * time.Second}
	go func() { _ = srv.Serve(listener) }()
	heartbeat.Store(time.Now().Unix())
	if os.Getenv("FOCUS_STUDIO_HEADLESS_TEST") != "1" {
		if e := launchBrowser(); e != nil {
			showError(e.Error())
			return
		}
	}
	// Exit when there have been no active clients for 90 seconds.
	// Data lives in the persistent Edge profile, not in this process.
	ticker := time.NewTicker(12 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if time.Since(time.Unix(heartbeat.Load(), 0)) > 90*time.Second {
			ctx, cancel := timeoutContext()
			_ = srv.Shutdown(ctx)
			cancel()
			return
		}
	}
}

// isolated for testing and avoids external dependencies
func timeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 2*time.Second)
}

func serveApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/")
	if name == "" {
		name = "index.html"
	}
	allowed := map[string]bool{"index.html": true, "sw.js": true, "manifest.webmanifest": true, "icon.svg": true, "preview.png": true}
	if !allowed[name] {
		http.NotFound(w, r)
		return
	}
	data, err := appFiles.ReadFile("web/" + name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	ext := filepath.Ext(name)
	typ := mime.TypeByExtension(ext)
	switch ext {
	case ".webmanifest":
		typ = "application/manifest+json"
	case ".svg":
		typ = "image/svg+xml"
	case ".js":
		typ = "application/javascript"
	case ".html":
		typ = "text/html; charset=utf-8"
	}
	if typ != "" {
		w.Header().Set("Content-Type", typ)
	}
	w.Header().Set("Cache-Control", "no-store")
	if name == "index.html" {
		data = bytes.Replace(data, []byte("</body>"), []byte(desktopHeartbeat+"\n</body>"), 1)
	}
	http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(data))
}

func checkAlreadyRunning() bool {
	c := http.Client{Timeout: 1250 * time.Millisecond}
	resp, err := c.Get(origin + "/_desktop_health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(resp.Body, 128))
	return resp.StatusCode == 200 && string(data) == magic && resp.Header.Get("X-Focus-Studio") == magic
}
