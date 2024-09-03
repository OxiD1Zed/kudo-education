package note_storage

import "kode-education/internal/domain"

type LocalNoteStorage struct {
	nextIndex uint
	notes     []domain.Note
}

func New() *LocalNoteStorage {
	return &LocalNoteStorage{
		nextIndex: 0,
		notes:     []domain.Note{},
	}
}

func (l *LocalNoteStorage) GetByIdUser(idUser uint) []domain.Note {
	var temp []domain.Note

	for i := 0; i < len(l.notes); i++ {
		if l.notes[i].IdUser == idUser {
			temp = append(temp, l.notes[i])
		}
	}

	return temp
}

func (l *LocalNoteStorage) Save(idUser uint, title, body string) uint {
	note := domain.Note{
		Id:     l.nextIndex,
		IdUser: idUser,
		Title:  title,
		Body:   body,
	}
	l.notes = append(l.notes, note)
	l.nextIndex++

	return note.Id
}
