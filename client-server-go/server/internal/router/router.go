package router

import (
	"fmt"

	apperrors "github.com/SaiVikrantG/server/internal/errors"
	"github.com/SaiVikrantG/server/internal/handlers"
	"github.com/SaiVikrantG/server/internal/models"
	"github.com/SaiVikrantG/server/internal/response"
)

type Router struct {
	routes map[string]handlers.Handler
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]handlers.Handler),
	}
}

func (r *Router) Handle(method, path string, h handlers.Handler) {
	key := fmt.Sprintf("%s %s", method, path)
	r.routes[key] = h
}

func (r *Router) ServeHTTP(w response.ResponseWriter, req *models.Request) {
	key := fmt.Sprintf("%s %s", req.HTTPMethod, req.Path)
	h, ok := r.routes[key]
	if !ok {
		httpErr := apperrors.NewNotFound("no route matched")
		w.Write(httpErr.StatusCode, []byte(httpErr.Message))
		return
	}
	h.ServeHTTP(w, req)
}
