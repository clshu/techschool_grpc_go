package service

import "sync"

// RatingStore is a store for storing laptop ratings.
type RatingStore interface {
	Add(laptopID string, score float64) (*Rating, error) // Add adds a new rating for a laptop
}

// Rating is a laptop rating.
type Rating struct {
	Count uint32
	Sum float64
}

// InMemoryRatingStore is an in-memory store for storing laptop ratings.
type InMemoryRatingStore struct {
	mutex   sync.RWMutex
	ratings map[string]*Rating
}

// NewInMemoryRatingStore creates a new InMemoryRatingStore.
func NewInMemoryRatingStore() *InMemoryRatingStore {
	return &InMemoryRatingStore{
		ratings: make(map[string]*Rating),
	}
}

// Add adds a new rating for a laptop.
func (store *InMemoryRatingStore) Add(laptopID string, score float64) (*Rating, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	rating := store.ratings[laptopID]
	if rating == nil {
		 rating = &Rating{
			Count: 1,
			Sum: score,
		 }	
	} else {
		rating.Count++
		rating.Sum += score
	} 
	
	store.ratings[laptopID] = rating

	return rating, nil
}

