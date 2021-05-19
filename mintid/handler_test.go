package function

/*
func TestHandle(t *testing.T) {
	type response struct {
		StatusCode int
		Body       []byte
	}

	type args struct {
		cpr      string
		password string
	}
	cases := []struct {
		name     string
		path     string
		method   string
		input    args
		expected response
	}{
		{
			name:     "Forsiden",
			path:     "/",
			method:   http.MethodGet,
			expected: response{http.StatusOK, []byte("Du er på forsiden")},
		},
	}

	for _, c := range cases {
		// Create a request to pass to our handler.
		input := strings.NewReader("")
		req, err := http.NewRequest(c.method, c.path, input)

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Handle)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		assert.Nil(t, err)
		assert.Equal(t, c.expected.StatusCode, rr.Code)
		assert.Equal(t, string(c.expected.Body), rr.Body.String())
	}

}
*/
