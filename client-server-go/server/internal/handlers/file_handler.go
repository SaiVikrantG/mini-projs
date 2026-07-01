package handlers

import (
	"os"
	"path/filepath"

	apperrors "github.com/SaiVikrantG/server/internal/errors"
	"github.com/SaiVikrantG/server/internal/models"
	"github.com/SaiVikrantG/server/internal/response"
)

type FileHandler struct {
	RootDir string
}

func (h *FileHandler) ServeHTTP(w response.ResponseWriter, r *models.Request) {
	path := filepath.Join(h.RootDir, filepath.Clean(r.Path))

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			httpErr := apperrors.NewNotFound("file not found")
			w.Write(httpErr.StatusCode, []byte(httpErr.Message))
			return
		}
		httpErr := apperrors.NewInternalServerError("failed to read file")
		w.Write(httpErr.StatusCode, []byte(httpErr.Message))
		return
	}

	w.Header()["Content-Type"] = contentType(path)
	w.Write(200, data)
}

func contentType(path string) string {
	switch filepath.Ext(path) {
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	default:
		return "text/plain"
	}
}
