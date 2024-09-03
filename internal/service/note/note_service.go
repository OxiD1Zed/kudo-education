package note_service

import (
	"fmt"
	"kode-education/internal/domain"
	"log/slog"
)

type NoteService struct {
	log           *slog.Logger
	textValidator TextValidator
	noteProvider  NoteProvider
}

type TextValidator interface {
	ValidateTextMultiLang(text string) (string, error)
}

type NoteProvider interface {
	Save(idUser uint, title, body string) uint
	GetByIdUser(idUser uint) []domain.Note
}

func New(
	log *slog.Logger,
	wordValidator TextValidator,
	noteProvider NoteProvider,
) *NoteService {
	return &NoteService{
		log:           log,
		textValidator: wordValidator,
		noteProvider:  noteProvider,
	}
}

func (s *NoteService) CreateNote(idUser uint, title, body string) (uint, error) {
	const op = "NoteService.CreateNote"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("creating a note")

	log.Info("the title is being validated")

	var err error
	title, err = s.textValidator.ValidateTextMultiLang(title)
	if err != nil {
		log.Warn("failed to validate the title", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("the body is being validated")

	body, err = s.textValidator.ValidateTextMultiLang(body)
	if err != nil {
		log.Warn("failed to validate the body", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("validation was successful")

	id := s.noteProvider.Save(idUser, title, body)
	log.Info("the note was saved successfully")

	return id, nil
}

func (s *NoteService) GetAllNotes(idUser uint) []domain.Note {
	const op = "NoteService.GetAllNotes"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("notes are being received")

	return s.noteProvider.GetByIdUser(idUser)
}
