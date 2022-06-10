package resources

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"io.hyperd/inspectmx"
	"io.hyperd/inspectmx/logger"
	"io.hyperd/inspectmx/service"
	"io.hyperd/inspectmx/transport/http/payloads"
)

type InspectorResource struct{}

// Routes defines the routes for /API/{version}/mx/**
func (rs InspectorResource) Routes() http.Handler {
	router := chi.NewRouter()
	router.Post("/", rs.Verify) // POST /users

	return router
}

func (rs InspectorResource) Verify(w http.ResponseWriter, r *http.Request) {
	data := &payloads.InspectorRequest{}

	logger.Debug("attempting to verify", logger.WithFields{
		"request": data,
	})

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	res, err := service.Instance().Inspector.Verify(data.Email, r.Context())

	i := inspectmx.Inspector{
		Email:    data.Email,
		Provider: res,
	}
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, payloads.NewInspectorResponse(&i))
}
