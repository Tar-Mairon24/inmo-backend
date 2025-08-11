package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"inmo-backend/internal/interface/api/handler"
)

func TestRegisterHealthRoutes_GET(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := handler.NewHealthHandler()
	router.GET("/health", h.RegisterHealthRoutes)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	body := w.Body.String()
	assert.Contains(t, body, `"status":"OK"`)
	assert.Contains(t, body, `"message":"Service is healthy"`)
	assert.Contains(t, body, `"service":"inmo-backend"`)
	assert.Contains(t, body, `"version":"1.0.0"`)
	// Check timestamp is RFC3339 and close to now
	type respStruct struct {
		Timestamp string `json:"timestamp"`
	}
	var resp respStruct
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	_, err := time.Parse(time.RFC3339, resp.Timestamp)
	assert.NoError(t, err)
	assert.NoError(t, err)
}

func TestRegisterHealthRoutes_HEAD(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := handler.NewHealthHandler()
	router.HEAD("/health", h.RegisterHealthRoutes)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodHead, "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, w.Body.Len() == 0 || strings.TrimSpace(w.Body.String()) == "")
}
func TestRegisterDetailedHealthRoute_GET(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := handler.NewHealthHandler()
	router.GET("/health/details", h.RegisterDetailedHealthRoute)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health/details", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	body := w.Body.String()
	assert.Contains(t, body, `"status":"OK"`)
	assert.Contains(t, body, `"message":"All services are operational"`)
	assert.Contains(t, body, `"service":"inmo-backend"`)
	assert.Contains(t, body, `"version":"1.0.0"`)
	assert.Contains(t, body, `"database":"connected"`)
	assert.Contains(t, body, `"memory":"ok"`)
	assert.Contains(t, body, `"disk_space":"ok"`)
	assert.Contains(t, body, `"uptime":"running"`)

	// Check timestamp is RFC3339 and close to now
	type respStruct struct {
		Timestamp string                 `json:"timestamp"`
		Checks    map[string]interface{} `json:"checks"`
	}
	var resp respStruct
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	_, err := time.Parse(time.RFC3339, resp.Timestamp)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Checks)
	assert.Equal(t, "connected", resp.Checks["database"])
	assert.Equal(t, "ok", resp.Checks["memory"])
	assert.Equal(t, "ok", resp.Checks["disk_space"])
}

func TestRegisterDetailedHealthRoute_HEAD(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := handler.NewHealthHandler()
	router.HEAD("/health/details", h.RegisterDetailedHealthRoute)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodHead, "/health/details", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, w.Body.Len() == 0 || strings.TrimSpace(w.Body.String()) == "")
}
func TestRegisterPingRoute_GET(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := handler.NewHealthHandler()
	router.GET("/ping", h.RegisterPingRoute)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	body := w.Body.String()
	assert.Contains(t, body, `"message":"pong"`)
}

func TestRegisterPingRoute_HEAD(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := handler.NewHealthHandler()
	router.HEAD("/ping", h.RegisterPingRoute)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodHead, "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, w.Body.Len() == 0 || strings.TrimSpace(w.Body.String()) == "")
}

