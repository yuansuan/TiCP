package iam_client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestValidUserIDMiddleware(t *testing.T) {
	// Create a new Gin router and add the middleware function
	router := gin.New()
	router.Use(ValidUserIDMiddleware("123"))

	// Define test routes that require the middleware function
	router.GET("/admin/test", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})
	router.GET("/system/test", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})
	router.GET("/test", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	// Create a new HTTP request with a valid user ID for the "/admin" URL
	req1, err := http.NewRequest("GET", "/admin/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req1.Header.Set(userIDKey, "123")

	// Create a new HTTP request with an invalid user ID for the "/admin" URL
	req2, err := http.NewRequest("GET", "/admin/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req2.Header.Set(userIDKey, "456")

	// Create a new HTTP request with a valid user ID for the "/system" URL
	req3, err := http.NewRequest("GET", "/system/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req3.Header.Set(userIDKey, "123")

	// Create a new HTTP request with an invalid user ID for the "/system" URL
	req4, err := http.NewRequest("GET", "/system/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req4.Header.Set(userIDKey, "456")

	// Create a new HTTP request with a valid user ID for a URL that does not start with "/admin" or "/system"
	req5, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req5.Header.Set(userIDKey, "123")

	// Test the middleware function with the valid user ID for the "/admin" URL
	resp1 := httptest.NewRecorder()
	router.ServeHTTP(resp1, req1)
	if resp1.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, resp1.Code)
	}

	// Test the middleware function with the invalid user ID for the "/admin" URL
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)
	if resp2.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d but got %d", http.StatusUnauthorized, resp2.Code)
	}

	// Test the middleware function with the valid user ID for the "/system" URL
	resp3 := httptest.NewRecorder()
	router.ServeHTTP(resp3, req3)
	if resp3.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, resp3.Code)
	}

	// Test the middleware function with the invalid user ID for the "/system" URL
	resp4 := httptest.NewRecorder()
	router.ServeHTTP(resp4, req4)
	if resp4.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d but got %d", http.StatusUnauthorized, resp4.Code)
	}

	// Test the middleware function with the valid user ID for a URL that does not start with "/admin" or "/system"
	resp5 := httptest.NewRecorder()
	router.ServeHTTP(resp5, req5)
	if resp5.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, resp5.Code)
	}
}
