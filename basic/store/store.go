package store

import (
	pb "grpc-log/login"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]time.Time
}

func New() Store {
	return Store{
		mu:   sync.RWMutex{},
		data: make(map[string]time.Time),
	}
}

func (s *Store) List() *pb.LoginList {
	results := &pb.LoginList{}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, v := range s.data {
		results.Logins = append(results.Logins, &pb.Login{Username: k, LastLogin: timestamppb.New(v)})
	}
	return results
}

func (s *Store) Save(key string, ts time.Time) {
	s.mu.Lock()
	s.data[key] = ts
	s.mu.Unlock()
}
