package user_storage

import (
	"kode-education/internal/domain"
	"kode-education/internal/storage"
)

type LocalUserStorage struct {
	users     []domain.User
	nextIndex uint
}

func New() *LocalUserStorage {
	return &LocalUserStorage{
		users:     []domain.User{},
		nextIndex: 0,
	}
}

func (l *LocalUserStorage) Save(username, password string) {
	user := domain.User{
		Id:       l.nextIndex,
		Username: username,
		Password: password,
	}
	l.users = append(l.users, user)
	l.nextIndex++
}

func (l *LocalUserStorage) GetByUsername(username string) (domain.User, error) {
	for i := 0; i < len(l.users); i++ {
		if l.users[i].Username == username {
			return l.users[i], nil
		}
	}

	return domain.User{}, storage.ErrorNotFound
}
