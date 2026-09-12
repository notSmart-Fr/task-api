package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSignUpReturnsSupabaseResponseWrapper(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/v1/signup" {
			t.Fatalf("request path = %q, want /auth/v1/signup", r.URL.Path)
		}
		if actual := r.Header.Get("apikey"); actual != "test-api-key" {
			t.Fatalf("apikey = %q, want test-api-key", actual)
		}

		var requestBody SignUpRequest
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Fatalf("decode signup request: %v", err)
		}
		if requestBody.Email != "user@example.com" || requestBody.Password != "password123" {
			t.Fatalf("unexpected signup request: %#v", requestBody)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"user":{"id":"user-id","email":"user@example.com"},"session":null}`))
	}))
	defer server.Close()

	endpoint := NewSignUpEndpoint(server.URL, "test-api-key", server.Client())
	request := httptest.NewRequest(http.MethodPost, "/auth/signup", strings.NewReader(`{"email":"user@example.com","password":"password123"}`))
	recorder := httptest.NewRecorder()

	endpoint.Handle(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}

	var responseBody struct {
		User struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
		Session json.RawMessage `json:"session"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&responseBody); err != nil {
		t.Fatalf("decode signup response: %v", err)
	}
	if responseBody.User.ID != "user-id" || responseBody.User.Email != "user@example.com" {
		t.Fatalf("unexpected user response: %#v", responseBody.User)
	}
	if string(responseBody.Session) != "null" {
		t.Fatalf("session = %s, want null", responseBody.Session)
	}
}
