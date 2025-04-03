package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LaulauChau/note-nest/internal/application/use_cases"
	"github.com/LaulauChau/note-nest/internal/domain/entities"
)

type NoteController struct {
	noteUseCase *use_cases.NoteUseCase
}

func NewNoteController(noteUseCase *use_cases.NoteUseCase) *NoteController {
	return &NoteController{
		noteUseCase: noteUseCase,
	}
}

type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Label   string `json:"label"`
}

type NoteResponse struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	IsArchived bool   `json:"is_archived"`
	Label      string `json:"label"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func (c *NoteController) CreateNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from context (added by the auth middleware)
	user, ok := r.Context().Value(UserContextKey).(*entities.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body
	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Create the note
	note, err := c.noteUseCase.CreateNote(ctx, user.ID, req.Title, req.Content, req.Label)
	if err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	// Return the note
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(NoteResponse{
		ID:         note.ID,
		Title:      note.Title,
		Content:    note.Content,
		IsArchived: note.IsArchived,
		Label:      note.Label,
		CreatedAt:  note.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  note.UpdatedAt.Format(time.RFC3339),
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
