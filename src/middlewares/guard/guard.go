package guard

import (
	"sync"
	"time"

	san "github.com/go-sanitize/sanitize"
	"github.com/microcosm-cc/bluemonday"
	"io.hyperd/inspectmx/logger"
)

var (
	once          sync.Once
	guardInstance Guard
)

type Guard struct {
	UGC       *bluemonday.Policy
	Strict    *bluemonday.Policy
	Sanitizer *san.Sanitizer
}

// Guard is provides a wrapper to sanitize the user's input.
func Instance() Guard {
	once.Do(func() {
		var err error
		guardInstance = Guard{
			UGC:    bluemonday.UGCPolicy(),
			Strict: bluemonday.StrictPolicy(),
		}

		guardInstance.UGC.AddTargetBlankToFullyQualifiedLinks(true)
		guardInstance.UGC.AddSpaceWhenStrippingTag(true)

		guardInstance.Strict.AddSpaceWhenStrippingTag(false)

		guardInstance.Sanitizer, err = san.New(san.OptionDateFormat{
			Input: []string{
				time.RFC3339,
				time.RFC3339Nano,
			},
			KeepFormat: false,
			Output:     time.RFC1123,
		})

		if err != nil {
			logger.Fatal("failed to configure go-sanitize", logger.WithFields{
				"error": err,
			})
		}
	})

	return guardInstance
}
