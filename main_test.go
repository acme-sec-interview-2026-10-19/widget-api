package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const testSigningKey = "test-signing-key"

func testRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	cfg := &Config{}
	cfg.Service.Name = "widget-api"
	cfg.AWS.Region = "us-east-1"
	cfg.Auth.JWTSigningKey = testSigningKey
	return newRouter(cfg)
}

func get(r http.Handler, path, bearer string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func signToken(t *testing.T, key, audience string) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "ops@acmewidgets.invalid",
		"aud": audience,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(key))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return signed
}

func TestHello(t *testing.T) {
	w := get(testRouter(), "/hello", "")

	if w.Code != http.StatusOK || w.Body.String() != "hello, world" {
		t.Fatalf("GET /hello = %d %q, want 200 %q", w.Code, w.Body.String(), "hello, world")
	}
}

func TestAdminStatusRequiresToken(t *testing.T) {
	cases := []struct {
		name   string
		bearer string
		want   int
	}{
		{"no token", "", http.StatusUnauthorized},
		{"garbage token", "not-a-jwt", http.StatusUnauthorized},
		{"wrong key", signToken(t, "some-other-key", "widget-api"), http.StatusUnauthorized},
		{"wrong audience", signToken(t, testSigningKey, "billing-api"), http.StatusUnauthorized},
		{"valid token", signToken(t, testSigningKey, "widget-api"), http.StatusOK},
	}

	r := testRouter()
	for _, tc := range cases {
		if got := get(r, "/admin/status", tc.bearer).Code; got != tc.want {
			t.Errorf("%s: GET /admin/status = %d, want %d", tc.name, got, tc.want)
		}
	}
}
