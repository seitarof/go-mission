package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
	"github.com/mileusna/useragent"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Getting UserAgent Information
	ua := ua.Parse(r.UserAgent())
	// Getting OS name

	ctx := context.WithValue(r.Context(), "OS", ua.OS)

	accessTimeBeforeHandler := time.Now()

	switch r.Method {
	case http.MethodGet:
		var prevId, size string
		if prevId = r.URL.Query().Get("prev_id"); prevId == "" {
			prevId = "0"
		}

		if size = r.URL.Query().Get("size"); size == "" {
			size = "5"
		}

		prevIdInt, err := strconv.Atoi(prevId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		sizeInt, err := strconv.Atoi(size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		todos, err := h.svc.ReadTODO(ctx, int64(prevIdInt), int64(sizeInt))
		if err != nil {
			log.Println(err)
			return
		}

		readTodoResponse := &model.ReadTODOResponse {
			TODOs: todos,
		}

		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(readTodoResponse)
		if err != nil {
			log.Println(err)
			return
		}
	case http.MethodPost:
		var createTodoRequest model.CreateTODORequest
		err := json.NewDecoder(r.Body).Decode(&createTodoRequest) // Decoding Request-Body
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		if createTodoRequest.Subject == "" {
			w.WriteHeader(http.StatusBadRequest) // 400: Bad Request
			return
		}

		todo, err := h.svc.CreateTODO(ctx, createTodoRequest.Subject, createTodoRequest.Description)
		if err != nil {
			log.Println(err)
			return
		}

		createTodoResponse := &model.CreateTODOResponse{
			TODO: todo,
		}

		err = json.NewEncoder(w).Encode(createTodoResponse) // Encoding CreateTODOResponse
		if err != nil {
			log.Println(err)
			return
		}
	case http.MethodPut:
		var updateTodoRequest model.UpdateTODORequest
		err := json.NewDecoder(r.Body).Decode(&updateTodoRequest) // Decoding Request-Body
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		if updateTodoRequest.ID == 0 || updateTodoRequest.Subject == "" {
			w.WriteHeader(http.StatusBadRequest) // 400: Bad Request
			return
		}
		todo, err := h.svc.UpdateTODO(ctx, updateTodoRequest.ID, updateTodoRequest.Subject, updateTodoRequest.Description)
		if err != nil {
			log.Println(err)
			return
		}

		updateTodoResponse := &model.UpdateTODOResponse{
			TODO: todo,
		}

		err = json.NewEncoder(w).Encode(updateTodoResponse) // Encoding UpdateTODOResponse
		if err != nil {
			log.Println(err)
			return
		}
	case http.MethodDelete:
		var deleteTodoRequest model.DeleteTODORequest
		err := json.NewDecoder(r.Body).Decode(&deleteTodoRequest) // Decoding Request-Body
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500: Internal Serve Error
			log.Println(http.StatusText(http.StatusInternalServerError))
			return
		}

		if len(deleteTodoRequest.IDs) == 0 {
			w.WriteHeader(http.StatusBadRequest) // 400: Bad Request
			return
		}

		err = h.svc.DeleteTODO(ctx, deleteTodoRequest.IDs)
		if err != nil {
			switch err.(type) {
			case *model.ErrNotFound:
				w.WriteHeader(http.StatusNotFound) // 404: Not Found
			default:
				w.WriteHeader(http.StatusBadGateway) // 502: Bad Gateway
			}
			return
		}

		deleteTodoResponse := &model.DeleteTODOResponse{}
		err = json.NewEncoder(w).Encode(deleteTodoResponse) // Encoding DeleteTODOResponse
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(http.StatusText(http.StatusInternalServerError))
			log.Println(err)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

	accessTimeAfterHandler := time.Now()
	accessTimeDiff := accessTimeAfterHandler.Sub(accessTimeBeforeHandler).Microseconds()
	accessLog := middleware.NewAccessLog(accessTimeBeforeHandler, accessTimeDiff, r.URL.Path, ctx.Value("OS").(string))
	accessLog.PrintAccessLogJson()
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	_, _ = h.svc.CreateTODO(ctx, "", "")
	return &model.CreateTODOResponse{}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
