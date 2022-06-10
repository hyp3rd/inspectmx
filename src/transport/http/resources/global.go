package resources

import (
	"net/http"
	"time"

	"github.com/go-chi/stampede"
)

type GlobalResource struct{}

// Prevents cache stampede https://en.wikipedia.org/wiki/Cache_stampede
// by only running a single data fetch operation per expired / missing key regardless of number of requests to that key.
func (rs GlobalResource) CacheHandler() func(next http.Handler) http.Handler {
	cached := stampede.Handler(512, 5*time.Second)
	return cached
}
