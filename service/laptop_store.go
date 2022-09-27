package service

import (
	"errors"
	"fmt"
	"learngrpc/pcbook/pb"
	"sync"

	"github.com/jinzhu/copier"
)

// ErrAlreadyExists is returned when we try to create a laptop with an ID that already exists.
var ErrAlreadyExists = errors.New("laptop already exists")

// LaptopStore is a store for laptop
type LaptopStore interface {
	Save(laptop *pb.Laptop) error
	Find(id string) (*pb.Laptop, error)
}

// InMemoryLaptopStore is an in-memory store for laptop
type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data map[string]*pb.Laptop
}

// NewInMemoryLaptopStore creates a new InMemoryLaptopStore
func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

// Save saves the laptop to the store
func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}
	// deep copy
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return fmt.Errorf("cannot copy laptop data: %w", err)
	}

	store.data[laptop.Id] = other
	return nil
}

// Find finds a laptop by ID
func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	laptop := store.data[id]
	if laptop == nil {
		return nil, fmt.Errorf("laptop not found: %v", id)
	}
	// deep copy
	result := &pb.Laptop{}
	err := copier.Copy(result, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}
	return result, nil
}

