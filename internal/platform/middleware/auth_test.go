package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBearerToken(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
		valid  bool
	}{
		{name: "standard", header: "Bearer token", want: "token", valid: true},
		{name: "extra whitespace", header: "  Bearer   token  ", want: "token", valid: true},
		{name: "case insensitive scheme", header: "bearer token", want: "token", valid: true},
		{name: "missing token", header: "Bearer", valid: false},
		{name: "duplicate scheme", header: "Bearer Bearer token", valid: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, valid := bearerToken(test.header)
			if got != test.want || valid != test.valid {
				t.Fatalf("bearerToken(%q) = (%q, %t), want (%q, %t)", test.header, got, valid, test.want, test.valid)
			}
		})
	}
}

func TestGetUserValidatesTokenWithSupabase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/v1/user" {
			t.Fatalf("path = %q, want /auth/v1/user", r.URL.Path)
		}
		if authorization := r.Header.Get("Authorization"); authorization != "Bearer access-token" {
			t.Fatalf("Authorization = %q, want bearer token", authorization)
		}
		if apiKey := r.Header.Get("apikey"); apiKey != "project-key" {
			t.Fatalf("apikey = %q, want project-key", apiKey)
		}
		_, _ = w.Write([]byte(`{"id":"user-id","email":"user@example.com"}`))
	}))
	defer server.Close()

	middleware := &AuthMiddleware{supabaseURL: server.URL, apiKey: "project-key"}
	request := httptest.NewRequest(http.MethodGet, "/protected/profile", nil)
	user, err := middleware.getUser(request, "access-token")
	if err != nil {
		t.Fatalf("getUser() error = %v", err)
	}
	if user.ID != "user-id" || user.Email != "user@example.com" {
		t.Fatalf("getUser() = %#v, want returned Supabase user", user)
	}
}
