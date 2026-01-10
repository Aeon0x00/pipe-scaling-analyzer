package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"pipe-scaling-analyzer/internal/domain"
	"pipe-scaling-analyzer/internal/usecase"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type HTTPHandlers struct {
	us     *usecase.UseCase
	logger *zap.SugaredLogger
}

func NewHTTPHandlers (us *usecase.UseCase,logger *zap.SugaredLogger) *HTTPHandlers {
	return &HTTPHandlers{
		us:     us,
		logger: logger,
	}
}



func (h *HTTPHandlers) HandleCreatePipe(w http.ResponseWriter, r *http.Request) {
	var dto domain.CreatePipeDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		errDTO := domain.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	fmt.Println("create pipe dto:", dto)


	pipe, err := h.us.CreatePipe(r.Context(), dto)
	if err != nil {
		errDTO := domain.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, domain.ErrInvalidPipeData) {
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(pipe)
}
func (h *HTTPHandlers) HandleListPipes(w http.ResponseWriter, r *http.Request) {
	pipes, err := h.us.ListPipes(r.Context())
	if err != nil {
		errDTO := domain.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(pipes)
}

func (h *HTTPHandlers) HandleGetPipe(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pipe, err := h.us.GetPipe(r.Context(), id)
	if err != nil {
		errDTO := domain.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, domain.ErrPipeNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}


	w.Header().Set("Content-Type", "application/json")
	
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(pipe)
}
func (h *HTTPHandlers) HandleDeletePipe(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.us.DeletePipe(r.Context(), id); err != nil {
		errDTO := domain.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, domain.ErrPipeNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func (h *HTTPHandlers) HandleCalculateScaling(w http.ResponseWriter, r *http.Request) {
	var reqDTO domain.CalcRequestDTO
	
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		errDTO := domain.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
		
	}

	res, err := h.us.CalculateScaling(r.Context(), reqDTO)
	if err != nil {
		errDTO := domain.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, domain.ErrPipeNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else if errors.Is(err, domain.ErrEmptyPoints) {
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	
	_ = json.NewEncoder(w).Encode(res)
}

func (h *HTTPHandlers) HandleGetPipeStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	status, err := h.us.GetPipeStatus(r.Context(), id)
	if err != nil {
		errDTO := domain.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		
		if errors.Is(err, domain.ErrPipeNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(status)
}
func (h *HTTPHandlers) HandleGetPipeHistory(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	history, err := h.us.GetPipeHistory(r.Context(), id)
	if err != nil {
		errDTO := domain.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, domain.ErrPipeNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(history)
}

func (h *HTTPHandlers) HandleListAlerts(w http.ResponseWriter, r *http.Request) {
	alerts, err := h.us.ListAlerts(r.Context())
	if err != nil {
		errDTO := domain.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(alerts)
}
func (h *HTTPHandlers) HandleFixAlert(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.us.FixAlert(r.Context(), id); err != nil {
		errDTO := domain.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, domain.ErrAlertNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}


	w.WriteHeader(http.StatusOK)
	
	_, _ = w.Write([]byte(`{"status":"fixed"}`))
}
