package function

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type response struct {
	StatusCode int
	Body       []byte
}

func TestHandleReturnsCorrectResponse(t *testing.T) {
	expected := response{
		StatusCode: http.StatusOK,
		Body:       []byte("Body: helloo"),
	}

	// Create a request to pass to our handler.
	input := strings.NewReader("helloo")
	req, err := http.NewRequest("GET", "/", input)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handle)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	assert.Nil(t, err)
	assert.Equal(t, rr.Code, expected.StatusCode)
	assert.Equal(t, rr.Body.String(), string(expected.Body))
}
