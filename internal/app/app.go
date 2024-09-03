package app

import (
	"kode-education/internal/app/rest"
	auth_service "kode-education/internal/service/auth"
	note_service "kode-education/internal/service/note"
	note_storage "kode-education/internal/storage/local/note"
	session_storage "kode-education/internal/storage/local/session"
	user_storage "kode-education/internal/storage/local/user"
	auth_server "kode-education/internal/transport/rest/auth"
	note_server "kode-education/internal/transport/rest/note"
	"kode-education/pkg/speller"
	"kode-education/pkg/text_validator"
	"log/slog"
	"net/http"
)

type App struct {
	RestApp rest.App
}

func New(
	log *slog.Logger,
	port uint,
) *App {
	userProvider := user_storage.New()
	userProvider.Save("example", "daaad6e5604e8e17bd9f108d91e26afe6281dac8fda0091040a7a6d7bd9b43b5")
	userProvider.Save("example2", "daaad6e5604e8e17bd9f108d91e26afe6281dac8fda0091040a7a6d7bd9b43b5")
	noteProvider := note_storage.New()
	sessionProvider := session_storage.New()

	client := http.Client{}
	speller := speller.New(client)

	wordValidator := text_validator.New(*speller)

	authService := auth_service.New(log, userProvider, sessionProvider)
	noteService := note_service.New(log, wordValidator, noteProvider)

	authServer := auth_server.New(*authService)
	noteServer := note_server.New(*noteService, *authService)

	restApp := rest.New(
		log,
		authServer,
		noteServer,
		port,
	)

	return &App{
		RestApp: *restApp,
	}
}
