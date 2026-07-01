package handlers

import (
	"github.com/SaiVikrantG/server/internal/models"
	"github.com/SaiVikrantG/server/internal/response"
)

type Handler interface {
	ServeHTTP(w response.ResponseWriter, r *models.Request)
}
