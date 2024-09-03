package note_server

import (
	"encoding/json"
	auth_service "kode-education/internal/service/auth"
	note_service "kode-education/internal/service/note"
	"kode-education/internal/transport/rest"
	"mime"
	"net/http"
)

type NoteServer interface {
	CreateNote(w http.ResponseWriter, req *http.Request)
	GetAllNotes(w http.ResponseWriter, req *http.Request)
}

type serverAPI struct {
	noteService note_service.NoteService
	authService auth_service.AuthService
}

func New(
	noteService note_service.NoteService,
	authService auth_service.AuthService,
) NoteServer {
	return &serverAPI{
		noteService: noteService,
		authService: authService,
	}
}

func (s *serverAPI) CreateNote(w http.ResponseWriter, req *http.Request) {
	type RequastNote struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	type ResponseId struct {
		Id uint `json:"id"`
	}

	contentType := req.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediaType != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	idUser := req.Context().Value("idUser").(uint)

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var reqNote RequastNote
	if err := dec.Decode(&reqNote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	id, err := s.noteService.CreateNote(idUser, reqNote.Title, reqNote.Body)
	if err != nil {
		http.Error(w, "failed to create a note", http.StatusInternalServerError)
		return
	}

	rest.RenderJSON(w, ResponseId{Id: id})
}

func (s *serverAPI) GetAllNotes(w http.ResponseWriter, req *http.Request) {
	idUser := req.Context().Value("idUser").(uint)

	rest.RenderJSON(w, s.noteService.GetAllNotes(idUser))
}
