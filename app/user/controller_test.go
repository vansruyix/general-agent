package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	repo := newFakeRepo()
	repo.users["1"] = &User{ID: "1", Username: "alice", Name: "Alice", RoleID: 1, Status: 1, CreateTime: 1000}
	svc := NewService(repo)
	ctrl := NewController(svc)

	r := gin.New()
	rg := r.Group("/api/v1/general-agent")
	rg.GET("/user/:id", ctrl.GetByID)
	rg.GET("/user", ctrl.List)
	rg.POST("/user", ctrl.Create)
	rg.PUT("/user/:id", ctrl.Update)
	rg.DELETE("/user/:id", ctrl.Delete)
	return r
}

func TestController_GetByID(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/general-agent/user/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp struct {
		Code int      `json:"code"`
		Msg  string   `json:"msg"`
		Data UserResp `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Username != "alice" {
		t.Errorf("expected username=alice, got %s", resp.Data.Username)
	}
}

func TestController_GetByID_NotFound(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/general-agent/user/999", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 10001 {
		t.Errorf("expected code=10001, got %d", resp.Code)
	}
}

func TestController_List(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/general-agent/user?pageNum=1&pageSize=10", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp struct {
		Code int      `json:"code"`
		Data PageResp `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Total != 1 {
		t.Errorf("expected total=1, got %d", resp.Data.Total)
	}
}

func TestController_Create(t *testing.T) {
	r := setupRouter()
	body := `{"username":"bob","password":"pass","name":"Bob"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/general-agent/user", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp struct {
		Code int      `json:"code"`
		Data UserResp `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Username != "bob" {
		t.Errorf("expected username=bob, got %s", resp.Data.Username)
	}
}

func TestController_Update(t *testing.T) {
	r := setupRouter()
	body := `{"name":"Alice Updated"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/general-agent/user/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp struct {
		Code int      `json:"code"`
		Data UserResp `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Name != "Alice Updated" {
		t.Errorf("expected Name='Alice Updated', got %s", resp.Data.Name)
	}
}

func TestController_Delete(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/general-agent/user/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp struct {
		Code int `json:"code"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 0 {
		t.Errorf("expected code=0, got %d", resp.Code)
	}
}