package resources

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/httprate"
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
	router.With((httprate.Limit(
		10,             // requests
		60*time.Second, // per duration
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			render.Status(r, http.StatusInternalServerError)
			render.Render(w, r, payloads.ErTooManyRequests)
		}),
	))).Post("/", rs.Verify) // POST /users

	router.Route("/config", func(r chi.Router) {
		r.Use(GlobalResource{}.CacheHandler())
		r.Get("/providers", rs.GetAllowedProviders) // GET
		r.Get("/mx-servers", rs.GetAllowedServers)  // GET

	})

	return router
}

func (rs InspectorResource) Verify(w http.ResponseWriter, r *http.Request) {
	data := &payloads.InspectorRequest{}

	logger.Debug("attempting to verify", logger.WithFields{
		"email": data,
	})

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	res, err := service.Instance().Inspector.Verify(data.Email, r.Context())

	i := inspectmx.Inspector{
		Email:    data.Email,
		Provider: res,
		Valid:    false,
	}
	if err != nil {
		render.Status(r, http.StatusExpectationFailed)
		i.Status = err.Error()
	}
	i.Valid = (i.Status == "")
	render.Render(w, r, payloads.NewInspectorResponse(&i))
}

func (rs InspectorResource) GetAllowedProviders(w http.ResponseWriter, r *http.Request) {
	res := service.Instance().Inspector.GetAllowedProviders(r.Context())

	render.Render(w, r, payloads.NewListResponse(res))
}

func (rs InspectorResource) GetAllowedServers(w http.ResponseWriter, r *http.Request) {
	res := service.Instance().Inspector.GetAllowedServers(r.Context())

	render.Render(w, r, payloads.NewListResponse(res))
}
