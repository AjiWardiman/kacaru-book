package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"tourism-booking/internal/middleware"
	"tourism-booking/internal/util"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuthMiddleware(t *testing.T) {
	jwtSecret := "test-secret-key-999"
	validToken, _ := util.GenerateToken(jwtSecret, 100, "admin", 1)

	r := gin.New()
	r.Use(middleware.Auth(jwtSecret))
	r.GET("/protected", func(c *gin.Context) {
		userID, _ := c.Get(middleware.ContextUserID)
		role, _ := c.Get(middleware.ContextRole)
		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
			"role":    role,
		})
	})

	t.Run("valid bearer token header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, resp.Code)
		}
	})

	t.Run("valid token in query param", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected?token="+validToken, nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, resp.Code)
		}
	})

	t.Run("missing token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, resp.Code)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid-token-string")
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, resp.Code)
		}
	})
}

func TestRequireRoleMiddleware(t *testing.T) {
	jwtSecret := "test-secret-key-999"
	adminToken, _ := util.GenerateToken(jwtSecret, 1, "admin", 1)
	userToken, _ := util.GenerateToken(jwtSecret, 2, "user", 1)

	r := gin.New()
	r.Use(middleware.Auth(jwtSecret))
	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.RequireRole("admin"))
	adminGroup.GET("/dashboard", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "welcome admin"})
	})

	t.Run("access allowed role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.Code)
		}
	})

	t.Run("access forbidden role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusForbidden {
			t.Errorf("expected status 403, got %d", resp.Code)
		}
	})
}
