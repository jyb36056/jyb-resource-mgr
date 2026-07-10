package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"jyb-resource-mgr/internal/service"
)

type HTTPServer struct {
	svc *service.UserService
}

// NewHTTPServer 依赖注入
func NewHTTPServer(svc *service.UserService) *HTTPServer {
	return &HTTPServer{svc: svc}
}

// RegisterRoutes 将http url和handler注册到ServeMux中
func (h *HTTPServer) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /users/{id}", h.GetUser)
	mux.HandleFunc("POST /users", h.CreateUser)
	mux.HandleFunc("PUT /users/{id}", h.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", h.DeleteUser)
	mux.HandleFunc("GET /users", h.ListUsers)
}

// 将相似的代码抽离出来，简化CRUD代码
// writeJSON 写入json响应
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError 写入错误响应
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// GetUser 获取用户
func (h *HTTPServer) GetUser(w http.ResponseWriter, r *http.Request) {
	// 从url参数中获取id
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	// 调用业务逻辑svc.GetUser方法
	user, err := h.svc.GetUser(r.Context(), int32(id))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// CreateUser 创建用户
func (h *HTTPServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	// 将请求体的json数据，绑定到req结构体对应的字段上
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// 调用业务逻辑svc.CreateUser方法
	user, err := h.svc.CreateUser(r.Context(), req.Name, req.Email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

// UpdateUser 更新用户
func (h *HTTPServer) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// 从url参数中获取id
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	// 将请求体的json数据，绑定到req结构体对应的字段上
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// 调用业务逻辑svc.UpdateUser方法
	user, err := h.svc.UpdateUser(r.Context(), int32(id), req.Name, req.Email)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// DeleteUser 删除用户
func (h *HTTPServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// 从url参数中获取id
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	// 调用业务逻辑svc.DeleteUser方法
	if err := h.svc.DeleteUser(r.Context(), int32(id)); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListUsers 获取用户列表
func (h *HTTPServer) ListUsers(w http.ResponseWriter, r *http.Request) {
	// 调用业务逻辑svc.ListUsers方法
	users, err := h.svc.ListUsers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, users)
}
