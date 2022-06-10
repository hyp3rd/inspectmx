package payloads

import (
	"errors"
	"net/http"

	"io.hyperd/inspectmx"
	"io.hyperd/inspectmx/logger"
	"io.hyperd/inspectmx/middlewares/validator"
)

// InspectorRequest is the request payload for Inspector data model.
//
// NOTE: It's good practice to have well defined request and response payloads
// so you can manage the specific inputs and outputs for clients, and also gives
// you the opportunity to transform data on input or output, for example
// on request, we'd like to protect certain fields and on output perhaps
// we'd like to include a computed field based on other values that aren't
// in the data model. Also, check out this awesome blog post on struct composition:
// http://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
type InspectorRequest struct {
	*inspectmx.Inspector
}

func (rs *InspectorRequest) Bind(r *http.Request) error {
	// If no Inspector fields are sent in the request return an
	// error to avoid a nil pointer dereference.
	if rs.Inspector == nil {
		return errors.New("missing required Inspector fields")
	}

	// rs.Email = guard.Instance().Strict.Sanitize(rs.Email)

	err := validator.Validate().Struct(rs.Inspector)
	if err != nil {
		logger.Error("failed validating the user input", logger.WithFields{
			"error": err,
		})

		err = validator.ReportError(err)

		return err
	}

	return nil
}

// InspectorResponse is the response payload for the Request data model.
// See NOTE above in InspectorRequest as well.
//
// In the InspectorResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type InspectorResponse struct {
	*inspectmx.Inspector
	// We add an additional field to the response here.
}

func NewInspectorResponse(v *inspectmx.Inspector) *InspectorResponse {
	resp := &InspectorResponse{Inspector: v}

	return resp
}

func (rd *InspectorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire

	return nil
}
