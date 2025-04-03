package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

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

func (c *NoteController) GetNoteByID(w http.ResponseWriter, r *http.Request) {
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

	// Get the note
	note, err := c.noteUseCase.GetNoteByID(ctx, noteID, user.ID)
	if err != nil {
		if err.Error() == "note not found" {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get note", http.StatusInternalServerError)
		return
	}

	// Return the note
	w.Header().Set("Content-Type", "application/json")
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

func (c *NoteController) GetActiveNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from context (added by the auth middleware)
	user, ok := r.Context().Value(UserContextKey).(*entities.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the notes
	notes, err := c.noteUseCase.GetActiveNotes(ctx, user.ID)
	if err != nil {
		http.Error(w, "Failed to get notes", http.StatusInternalServerError)
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

func (c *NoteController) GetArchivedNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from context (added by the auth middleware)
	user, ok := r.Context().Value(UserContextKey).(*entities.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the notes
	notes, err := c.noteUseCase.GetArchivedNotes(ctx, user.ID)
	if err != nil {
		http.Error(w, "Failed to get archived notes", http.StatusInternalServerError)
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

type UpdateNoteRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	IsArchived bool   `json:"is_archived"`
	Label      string `json:"label"`
}

func (c *NoteController) UpdateNote(w http.ResponseWriter, r *http.Request) {
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

	// Parse the request body
	var req UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Update the note
	note, err := c.noteUseCase.UpdateNote(ctx, noteID, user.ID, req.Title, req.Content, req.Label, req.IsArchived)
	if err != nil {
		if err.Error() == "note not found" {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update note", http.StatusInternalServerError)
		return
	}

	// Return the updated note
	w.Header().Set("Content-Type", "application/json")
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

func (c *NoteController) DeleteNote(w http.ResponseWriter, r *http.Request) {
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

	// Delete the note
	err := c.noteUseCase.DeleteNote(ctx, noteID, user.ID)
	if err != nil {
		if err.Error() == "note not found" {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	// Return success with no content
	w.WriteHeader(http.StatusNoContent)
}
