package http

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"

	"io.hyperd/inspectmx/transport/http/resources"
	"io.hyperd/inspectmx/types"
)

func Service() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RedirectSlashes)
	router.Use(middleware.CleanPath)
	router.Use(middleware.Timeout(20 * time.Second))
	router.Use(render.SetContentType(render.ContentTypeJSON))
	// router.Use(middleware.AllowContentType(types.MimeType))
	// router.Use(resources.SessionResource{}.SessionCtx)

	// CORS implementation
	// see: https://developer.github.com/v3/#cross-origin-resource-sharing
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://ideas.ing.net"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("orange ideas"))
	})

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	router.Route("/api", func(router chi.Router) {

		router.Route("/v1", func(router chi.Router) {
			// API version
			router.Use(apiVersionCtx("v1")) // inject the API version in the http request context

			// mount the V1 routes
			router.Mount("/imx", resources.InspectorResource{}.Routes())
		})
	})

	return router
}

// apiVersionCtx injects in the routes the API version, enabling the possibility of different behaviors based on the version number.
// Here to follow an examlpe:
/*
	apiVersion := r.Context().Value(types.CtxApiVersionKey{}).(string)
	switch apiVersion {
	case "v1":
		idea <- v1.NewIdeaResponse(idea)
	case "v2":
		idea <- v2.NewIdeaResponse(idea)
	default:
		idea <- v3.NewIdeaResponse(idea)
	}
*/
func apiVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), types.CtxApiVersionKey, version))
			next.ServeHTTP(w, r)
		})
	}
}
