package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"tz/internal/domain"
	"tz/internal/service"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewRouter(service *service.TaskService) *chi.Mux {
	r := chi.NewRouter()
	h := &TaskHandler{service: service}

	r.Post("/tasks", h.CreateTask)
	r.Get("/tasks/{id}", h.GetTaskStatus)
	r.Get("/tasks/{id}/result", h.GetTaskResult)

	return r
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req domain.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request"})
		return
	}

	taskID, err := h.service.CreateTask(r.Context(), &req)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, map[string]string{"task_id": taskID})
}

func (h *TaskHandler) GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	status, err := h.service.GetTaskStatus(r.Context(), taskID)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "task not found"})
		return
	}

	render.JSON(w, r, map[string]string{"status": string(status)})
}

func (h *TaskHandler) GetTaskResult(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	result, err := h.service.GetTaskResult(r.Context(), taskID)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "task not found or not completed"})
		return
	}

	render.JSON(w, r, result)
}
