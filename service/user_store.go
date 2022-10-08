package service

import "sync"

// UserStore is a store for storing users.
type UserStore interface {
	Save(user *User) error
	Find(username string) (*User, error)	
}

// InMemoryUserStore is an in-memory store for storing users.
type InMemoryUserStore struct {
	mutex sync.RWMutex
	users map[string]*User
}

// NewInMemoryUserStore creates a new InMemoryUserStore.
func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]*User),
	}
}

// Save saves a user to the store.
func (store *InMemoryUserStore) Save(user *User) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.users[user.Username] != nil {
		return ErrAlreadyExists
	}
		
	store.users[user.Username] = user.Clone()

	return nil
}

// Find finds a user by username.
func (store *InMemoryUserStore) Find(username string) (*User, error)	{
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	user := store.users[username]
	if user == nil {
		return nil, nil
	}

	return user.Clone(), nil	
}