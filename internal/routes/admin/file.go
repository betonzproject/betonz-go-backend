package admin

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/BetOnz-Company/betonz-go/internal/acl"
	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/auth"
	"github.com/BetOnz-Company/betonz-go/internal/utils/fileutils"

	"github.com/go-chi/chi/v5"
)

func GetFile(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ViewFiles) != nil {
			return
		}

		filename := chi.URLParam(r, "filename")

		if !isValidFilename(filename) {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}

		workDir, _ := os.Getwd()
		fileDir := http.Dir(filepath.Join(workDir, "uploads", filename[:2], filename[2:4], filename[4:]))

		if !fileutils.FileExists(string(fileDir)) {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}

		http.ServeFile(w, r, string(fileDir))
	}
}

func isValidFilename(filename string) bool {
	// Regular expression to match a SHA-256 hash string
	// It should consist of exactly 64 characters, each being a hexadecimal digit
	validFilename := regexp.MustCompile(`^[0-9a-fA-F]{64}$`)
	return validFilename.MatchString(filename)
}
