package session_storage

import "kode-education/internal/storage"

type LocalSessionStorage struct {
	sessions map[string]uint
}

func New() *LocalSessionStorage {
	return &LocalSessionStorage{
		sessions: map[string]uint{},
	}
}

func (l *LocalSessionStorage) Save(token string, idUser uint) {
	l.sessions[token] = idUser
}

func (l *LocalSessionStorage) GetByToken(token string) (uint, error) {
	if val, ok := l.sessions[token]; ok {
		return val, nil
	}

	return 0, storage.ErrorNotFound
}
