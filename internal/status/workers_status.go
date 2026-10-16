package status

import (
	"errors"
	"sync"
	"time"

	"github.com/sebx771/monitor-bot/internal/worker"
)

type Status string

const(
	 Operational Status = "operational"
	 Error Status = "error"
	 Paused Status= "paused"
	 Degraded    Status = "degraded"
)

type MonitoredService struct {
	Target  string `json:"target"` // Eg: "mysql-13396a35", "mc-server", "/v1/pay"
	Scope   string `json:"scope,omitempty"`  // Eg: "terranova1244-6338" (optional)
	State   string `json:"state"`   // Eg: "RUNNING", "POWEROFF"
	Message string `json:"message,omitempty"` // eg: "api supabase returned 401 "
}

type WorkerStatus struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Status    Status             `json:"status"`
	LastRun   time.Time          `json:"last_run"`
	Services  []MonitoredService `json:"services"`
	Message   string             `json:"message,omitempty"`
}

func NewWorkerStatus(id, name string, s Status, services []MonitoredService, message string) (WorkerStatus,error){
   switch s {
	case Operational, Degraded, Error, Paused:
		// valid states , continue
	default:
		return WorkerStatus{}, errors.New("invalid worker state: " + string(s))
	}
	return WorkerStatus{
		ID:       id,
		Name:     name,
		Status:   s,
		LastRun:  time.Now().UTC(),
		Services: services,
		Message:  message,
	},nil
}

// status storage
type Store struct {
	mu       sync.RWMutex
	workers  map[string]WorkerStatus
	registry map[string]*worker.Worker
}

func NewStore() *Store {
	return &Store{
		workers:  make(map[string]WorkerStatus),
		registry: make(map[string]*worker.Worker),
	}
}

func (s *Store) RegisterWorker(id string, w *worker.Worker) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.registry[id] = w
}

func (s *Store) GetWorker(id string) (*worker.Worker, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w, ok := s.registry[id]
	return w, ok
}

func (s *Store) Set(ws WorkerStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.workers[ws.ID] = ws
}

func (s *Store) GetAll() []WorkerStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]WorkerStatus, 0, len(s.workers))
	for _, w := range s.workers {
		result = append(result, w)
	}
	return result
}