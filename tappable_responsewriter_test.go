package tapresp

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeader(t *testing.T) {

	t.Run("headers are available in both Response and ResponseWriter", func(t *testing.T) {
		tappable := New(httptest.NewRecorder())

		hKey := "content-type"
		hVal := "application/json"

		tappable.Header().Set(hKey, hVal)

		actual := tappable.Response().Header.Get(hKey)
		if actual != hVal {
			t.Fatalf("expected %v but got %v", actual, hVal)
		}
	})

}

func TestStatusCode(t *testing.T) {

	t.Run("status code should be 0 if unset", func(t *testing.T) {
		tappable := New(httptest.NewRecorder())

		if tappable.StatusCode() != 0 {
			t.Fatalf("status code should be 0 if unset, got %v", tappable.StatusCode())
		}
	})

	t.Run("StatusCode() should return correct status code", func(t *testing.T) {
		code := http.StatusBadRequest

		tappable := TappableResponseWriter{
			status: code,
		}

		if tappable.StatusCode() != http.StatusBadRequest {
			t.Fatalf("expected %v if unset, got %v", code, tappable.StatusCode())
		}
	})
}

func TestWriteHeader(t *testing.T) {
	t.Run("status code is available to underlying ResponseWriter", func(t *testing.T) {
		rw := httptest.NewRecorder()
		tappable := New(rw)

		expected := http.StatusOK
		tappable.WriteHeader(expected)

		if tappable.StatusCode() != expected {
			t.Fatalf("expected %d, got %v", expected, tappable.StatusCode())
		}

		if rw.Code != expected {
			t.Fatalf("status code not set in underlying ResponseWriter")
		}
	})
}

func TestBody(t *testing.T) {
	t.Run("Body() should return the response body stored in ResponseWriter", func(t *testing.T) {
		expected := []byte{'s', 'u', 'c', 'c', 'e', 's', 's'}

		tappable := TappableResponseWriter{
			body: *bytes.NewBuffer(expected),
		}

		if !bytes.Equal(tappable.Body(), expected) {
			t.Fatalf("body %v does not match expected body %v", tappable.Body(), expected)
		}
	})
}

func TestWrite(t *testing.T) {
	t.Run("value written with Wrtie() is available to be tapped", func(t *testing.T) {
		rw := httptest.NewRecorder()
		tappable := New(rw)

		expected := []byte{'s', 'u', 'c', 'c', 'e', 's', 's'}

		n, err := tappable.Write(expected)
		if err != nil {
			t.Fatalf("failed to write data: %v", err)
		}

		if n != len(expected) {
			t.Fatalf("bytes written %d does not equal expected %v", n, len(expected))
		}

		// Verify Tappable Body has been populated
		if !bytes.Equal(tappable.Body(), expected) {
			t.Fatalf("body %v does not match expected body %v", tappable.Body(), expected)
		}

		// Verify Underlying ResponseWriter body has been populated
		if !bytes.Equal(rw.Body.Bytes(), expected) {
			t.Fatalf("body %v does not match expected body %v", rw.Body.Bytes(), expected)
		}

	})
}

func TestResponse(t *testing.T) {
	t.Run("response should contain all values held in underlying ResponseWriter", func(t *testing.T) {
		expectedBody := []byte{'s', 'u', 'c', 'c', 'e', 's', 's'}
		expectedStatus := http.StatusOK
		expectedHeaders := map[string]string{
			"foo":  "bar",
			"test": "success",
		}

		tappable := TappableResponseWriter{
			w:      httptest.NewRecorder(),
			status: expectedStatus,
			body:   *bytes.NewBuffer(expectedBody),
		}

		for key, val := range expectedHeaders {
			tappable.Header().Set(key, val)
		}

		resp := tappable.Response()

		if resp.StatusCode != expectedStatus {
			t.Fatalf("expected status code %v, got %v", expectedStatus, resp.StatusCode)
		}

		defer resp.Body.Close()
		actualBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response body %v", err)
		}

		if !bytes.Equal(actualBody, expectedBody) {
			t.Fatalf("body %v does not match expected body %v", tappable.Body(), expectedBody)
		}

		// Validate all headers are present in Response
		if len(expectedHeaders) != len(resp.Header) {
			t.Fatalf("response headers not the same size as expected")
		}

		for key, val := range expectedHeaders {
			if resp.Header.Get(key) != val {
				t.Fatalf("invalid header value at %v", key)
			}
		}
	})

	t.Run("response headers should be a copy of those in ResponseWriter", func(t *testing.T) {
		tappable := New(httptest.NewRecorder())

		hKey := "content-type"
		hVal := "application/json"

		tappable.Header().Set(hKey, hVal)

		tappable.Response().Header.Set(hKey, "foobar")

		actual := tappable.Header().Get(hKey)
		if actual == "foobar" {
			t.Fatalf("header was modified to %v", actual)
		}
	})
}
