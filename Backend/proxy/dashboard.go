package proxy

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type dashboardServer struct {
	root string
}

func newDashboardServer() *dashboardServer {
	for _, dir := range dashboardCandidates() {
		index := filepath.Join(dir, "index.html")
		if info, err := os.Stat(index); err == nil && !info.IsDir() {
			return &dashboardServer{root: dir}
		}
	}
	return nil
}

func dashboardCandidates() []string {
	candidates := []string{}
	if dir := strings.TrimSpace(os.Getenv("DASHBOARD_DIR")); dir != "" {
		candidates = append(candidates, dir)
	}
	if exe, err := os.Executable(); err == nil {
		base := filepath.Dir(exe)
		candidates = append(candidates, filepath.Join(base, "dashboard"))
		candidates = append(candidates, filepath.Join(base, "..", "Dashboard", "out"))
	}
	candidates = append(candidates,
		"dashboard",
		filepath.Join("Backend", "dashboard"),
		filepath.Join("Dashboard", "out"),
		filepath.Join("..", "Dashboard", "out"),
	)
	return candidates
}

func (s *dashboardServer) Serve(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return false
	}
	target := filepath.Clean(string(filepath.Separator) + strings.TrimPrefix(r.URL.Path, "/"))
	target = strings.TrimPrefix(target, string(filepath.Separator))
	if target == "." || strings.HasPrefix(target, "..") {
		target = ""
	}
	file := filepath.Join(s.root, target)
	if info, err := os.Stat(file); err == nil && !info.IsDir() {
		http.ServeFile(w, r, file)
		return true
	}
	if info, err := os.Stat(file); err == nil && info.IsDir() {
		index := filepath.Join(file, "index.html")
		if indexInfo, err := os.Stat(index); err == nil && !indexInfo.IsDir() {
			http.ServeFile(w, r, index)
			return true
		}
	}
	htmlFile := filepath.Join(s.root, target+".html")
	if info, err := os.Stat(htmlFile); err == nil && !info.IsDir() {
		http.ServeFile(w, r, htmlFile)
		return true
	}
	index := filepath.Join(s.root, "index.html")
	if info, err := os.Stat(index); err == nil && !info.IsDir() {
		http.ServeFile(w, r, index)
		return true
	}
	return false
}
