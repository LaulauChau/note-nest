package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LaulauChau/note-nest/internal/application/use_cases"
	"github.com/LaulauChau/note-nest/internal/domain/entities"
)

type LabelController struct {
	labelUseCase *use_cases.LabelUseCase
}

func NewLabelController(labelUseCase *use_cases.LabelUseCase) *LabelController {
	return &LabelController{
		labelUseCase: labelUseCase,
	}
}

type CreateLabelRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type LabelResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (c *LabelController) CreateLabel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from context (added by the auth middleware)
	user, ok := r.Context().Value(UserContextKey).(*entities.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body
	var req CreateLabelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Validate color (should be a valid hex color or other format)
	if req.Color == "" {
		// Default color if not provided
		req.Color = "#3498db" // Blue
	}

	// Create the label
	label, err := c.labelUseCase.CreateLabel(ctx, user.ID, req.Name, req.Color)
	if err != nil {
		if err.Error() == "label with this name already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create label", http.StatusInternalServerError)
		return
	}

	// Return the label
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(LabelResponse{
		ID:        label.ID,
		Name:      label.Name,
		Color:     label.Color,
		CreatedAt: label.CreatedAt.Format(time.RFC3339),
		UpdatedAt: label.UpdatedAt.Format(time.RFC3339),
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
