package service

import (
	"context"
	"errors"
	"fmt"
	"learngrpc/pcbook/pb"
	"log"
	"sync"

	"github.com/jinzhu/copier"
)

// ErrAlreadyExists is returned when we try to create a laptop with an ID that already exists.
var ErrAlreadyExists = errors.New("laptop already exists")

// LaptopStore is a store for laptop
type LaptopStore interface {
	Save(laptop *pb.Laptop) error
	Find(id string) (*pb.Laptop, error)
	Search(ctx context.Context, filter *pb.LaptopFilter, found func(ldaptop *pb.Laptop) error) error
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
	other, err := deepCopy(laptop)
	if err != nil {
		return err
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
		return nil, nil
	}
	// deep copy
	return deepCopy(laptop)
}
// Search searches for laptops that match the filter criteria
func (store *InMemoryLaptopStore) Search(ctx context.Context, filter *pb.LaptopFilter, found func(laptop *pb.Laptop) error) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	for _, laptop := range store.data {
		// time.Sleep(1 * time.Second)
		// log.Print("checking laptop id: ", laptop.Id)
		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			log.Print("context is cancelled")
			return errors.New("context is cancelled")
		}
		if isQualified(filter, laptop) {
			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}

			err = found(other)
			if err != nil {
				return err
			}
		}			
	}
	return nil
}


func isQualified(filter *pb.LaptopFilter, laptop *pb.Laptop) bool {
	if laptop.GetPriceUsd() > filter.MaxPriceUsd {
		return false
	}
	if laptop.GetCpu().NumCores < filter.MinCpuCores {
		return false
	}
	if laptop.GetCpu().GetMinGhz() < filter.MinCpuGhz {
		return false
	}
	if toBits(laptop.GetMemory()) < toBits(filter.MinMemory) {
		return false
	}
	return true
}

func toBits(memory *pb.Memory) uint64 {
	value := memory.GetValue()

	switch(memory.GetUnit()) {
		case pb.Memory_BIT:
			return value
			case pb.Memory_BYTE:
				return value << 3 // 8 = 2^3
			case pb.Memory_KILLOBYTE:
				return value << 13 // 8 * 1024 = 2^13
			case pb.Memory_MEGABYTE:
				return value << 23 // 8 * 1024 * 1024 = 2^23
			case pb.Memory_GIGABYTE:
				return value << 33 // 8 * 1024 * 1024 * 1024 = 2^33	
			case pb.Memory_TERABYTE:
					return value << 43 // 8 * 1024 * 1024 * 1024 * 1024 = 2^43
			default:
				return 0
			}
}

func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	result := &pb.Laptop{}
	err := copier.Copy(result, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}
	return result, nil
}	