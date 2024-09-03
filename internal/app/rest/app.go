package rest

import (
	"fmt"
	auth_server "kode-education/internal/transport/rest/auth"
	note_server "kode-education/internal/transport/rest/note"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	log        *slog.Logger
	authServer auth_server.AuthServer
	noteServer note_server.NoteServer
	port       uint
}

func New(
	log *slog.Logger,
	authHandler auth_server.AuthServer,
	noteHandler note_server.NoteServer,
	port uint,
) *App {
	return &App{
		log:        log,
		authServer: authHandler,
		noteServer: noteHandler,
		port:       port,
	}
}

func (a *App) Run() {
	const op = "rest.App.Run"

	log := a.log.With(slog.String("op", op))

	log.Info("the rest server is running")

	err := http.ListenAndServe(fmt.Sprint("0.0.0.0:", a.port), a.setupRouter())
	if err != nil {
		log.Warn("the server could not be started", slog.String("error", err.Error()))
	}
	log.Info("the server is stopped")
}

func (a *App) setupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/auth", a.authServer.Authenticate).Methods("GET")
	router.HandleFunc("/note", a.authServer.AuthMiddleware(a.noteServer.GetAllNotes)).Methods("GET")
	router.HandleFunc("/note", a.authServer.AuthMiddleware(a.noteServer.CreateNote)).Methods("POST")

	return router
}
