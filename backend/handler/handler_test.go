package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rohitsh16/go-blog/backend/handler"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRegisterRoutes(t *testing.T) {
	router := gin.New()
	svc := handler.Service{
		GinFramework: router,
	}

	handler.RegisterRoutes(svc)

	routes := router.Routes()
	expectedRoutes := map[string]string{
		"POST:/v1/post":            "CreatePostHandler",
		"GET:/":                    "HomeHandler",
		"GET:/create":              "CreatePostHandler",
		"POST:/create":             "CreatePostHandler",
		"GET:/writer":              "WriterDashboardHandler",
		"GET:/writer/edit/:id":     "WriterEditHandler",
		"POST:/writer/edit/:id":    "WriterEditHandler",
		"GET:/admin":               "AdminHandler",
		"GET:/post/:identifier":    "PostDetailHandler",
	}

	registered := make(map[string]bool)
	for _, r := range routes {
		key := r.Method + ":" + r.Path
		registered[key] = true
	}

	for expectedKey := range expectedRoutes {
		if !registered[expectedKey] {
			t.Errorf("expected route %s to be registered", expectedKey)
		}
	}
}

func TestResolveTemplatePath(t *testing.T) {
	templates := []string{
		"index_user.html",
		"post.html",
		"create.html",
		"writer_dashboard.html",
		"writer_edit.html",
		"404.html",
	}

	for _, tmplName := range templates {
		path := handler.ResolveTemplatePath(tmplName)
		if path == "" {
			t.Errorf("expected non-empty path for template %s", tmplName)
		}

		tmpl, err := handler.ParseTemplate(tmplName)
		if err != nil {
			t.Errorf("failed to parse template %s at %s: %v", tmplName, path, err)
		}
		if tmpl == nil {
			t.Errorf("expected parsed template object for %s", tmplName)
		}
	}
}

func TestCreatePostHandler_JSONValidation(t *testing.T) {
	router := gin.New()
	svc := handler.Service{
		GinFramework: router,
	}
	router.POST("/v1/post", svc.CreatePostHandler)

	// Missing title and content
	invalidJSON := `{"author_name": "Test"}`
	req, _ := http.NewRequest(http.MethodPost, "/v1/post", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected HTTP 400 Bad Request for missing required fields, got %d", w.Code)
	}
}

func TestCreatePostHandler_FormValidation(t *testing.T) {
	router := gin.New()
	svc := handler.Service{
		GinFramework: router,
	}
	router.POST("/create", svc.CreatePostHandler)

	// Empty title and content form POST
	req, _ := http.NewRequest(http.MethodPost, "/create", bytes.NewBufferString("title=&content="))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected HTTP 400 Bad Request for empty form fields, got %d", w.Code)
	}
}

func TestPostDetailHandler_MissingIdentifier(t *testing.T) {
	router := gin.New()
	svc := handler.Service{
		GinFramework: router,
	}
	router.GET("/post/:identifier", svc.PostDetailHandler)

	// Sending request with whitespace only
	req, _ := http.NewRequest(http.MethodGet, "/post/%20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected HTTP 400 Bad Request for blank identifier, got %d", w.Code)
	}
}
