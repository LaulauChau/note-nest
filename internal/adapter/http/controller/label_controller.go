package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

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

func (c *LabelController) GetLabelByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from context (added by the auth middleware)
	user, ok := r.Context().Value(UserContextKey).(*entities.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get label ID from URL parameter
	labelID := chi.URLParam(r, "labelID")
	if labelID == "" {
		http.Error(w, "Label ID is required", http.StatusBadRequest)
		return
	}

	// Get the label
	label, err := c.labelUseCase.GetLabelByID(ctx, labelID, user.ID)
	if err != nil {
		if err.Error() == "label not found" {
			http.Error(w, "Label not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get label", http.StatusInternalServerError)
		return
	}

	// Return the label
	w.Header().Set("Content-Type", "application/json")
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

func (c *LabelController) GetLabels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from context (added by the auth middleware)
	user, ok := r.Context().Value(UserContextKey).(*entities.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get labels for the user
	labels, err := c.labelUseCase.GetLabelsByUser(ctx, user.ID)
	if err != nil {
		http.Error(w, "Failed to get labels", http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := make([]LabelResponse, len(labels))
	for i, label := range labels {
		response[i] = LabelResponse{
			ID:        label.ID,
			Name:      label.Name,
			Color:     label.Color,
			CreatedAt: label.CreatedAt.Format(time.RFC3339),
			UpdatedAt: label.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Return the labels
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (c *LabelController) GetNoteLabels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from context (added by the auth middleware)
	user, ok := r.Context().Value(UserContextKey).(*entities.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get note ID from URL parameter
	noteID := chi.URLParam(r, "noteID")
	if noteID == "" {
		http.Error(w, "Note ID is required", http.StatusBadRequest)
		return
	}

	// Get labels for the note
	labels, err := c.labelUseCase.GetLabelsForNote(ctx, noteID, user.ID)
	if err != nil {
		if err.Error() == "note not found" {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get labels for note", http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := make([]LabelResponse, len(labels))
	for i, label := range labels {
		response[i] = LabelResponse{
			ID:        label.ID,
			Name:      label.Name,
			Color:     label.Color,
			CreatedAt: label.CreatedAt.Format(time.RFC3339),
			UpdatedAt: label.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Return the labels
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (c *LabelController) GetNotesForLabel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from context (added by the auth middleware)
	user, ok := r.Context().Value(UserContextKey).(*entities.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get label ID from URL parameter
	labelID := chi.URLParam(r, "labelID")
	if labelID == "" {
		http.Error(w, "Label ID is required", http.StatusBadRequest)
		return
	}

	// Get notes for the label
	notes, err := c.labelUseCase.GetNotesForLabel(ctx, labelID, user.ID)
	if err != nil {
		if err.Error() == "label not found" {
			http.Error(w, "Label not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get notes for label", http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := make([]NoteResponse, len(notes))
	for i, note := range notes {
		response[i] = NoteResponse{
			ID:         note.ID,
			Title:      note.Title,
			Content:    note.Content,
			IsArchived: note.IsArchived,
			Label:      note.Label,
			CreatedAt:  note.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  note.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Return the notes
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
