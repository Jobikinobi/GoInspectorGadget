package casefile

import (
	"fmt"
	"time"
)

// Case status constants
const (
	StatusOpen     = "OPEN"
	StatusClosed   = "CLOSED"
	StatusSuspended = "SUSPENDED"
	StatusCold     = "COLD"
)

// Priority levels
const (
	PriorityHigh   = "HIGH"
	PriorityMedium = "MEDIUM"
	PriorityLow    = "LOW"
)

// Case represents a case file
type Case struct {
	ID          string
	CaseNumber  string
	Title       string
	Description string
	CaseType    string
	Status      string
	Priority    string
	AssignedTo  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ClosedAt    time.Time
}

// CaseRepository defines the interface for case data operations
type CaseRepository interface {
	Save(c *Case) error
	Find(id string) (*Case, error)
	FindByCaseNumber(caseNumber string) (*Case, error)
	Search(query string) ([]*Case, error)
	List(limit, offset int) ([]*Case, error)
	Update(c *Case) error
	Delete(id string) error
}

// CaseService provides case management functionality
type CaseService struct {
	repo CaseRepository
}

// NewCaseService creates a new case service with the given repository
func NewCaseService(repo CaseRepository) *CaseService {
	return &CaseService{
		repo: repo,
	}
}

// CreateCase creates a new case
func (s *CaseService) CreateCase(c *Case) error {
	// Set default values
	if c.ID == "" {
		c.ID = generateID()
	}
	if c.Status == "" {
		c.Status = StatusOpen
	}
	if c.Priority == "" {
		c.Priority = PriorityMedium
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	c.UpdatedAt = c.CreatedAt

	return s.repo.Save(c)
}

// GetCase retrieves a case by ID
func (s *CaseService) GetCase(id string) (*Case, error) {
	return s.repo.Find(id)
}

// UpdateCase updates a case
func (s *CaseService) UpdateCase(c *Case) error {
	c.UpdatedAt = time.Now()
	return s.repo.Update(c)
}

// CloseCase closes a case
func (s *CaseService) CloseCase(id string) error {
	c, err := s.repo.Find(id)
	if err != nil {
		return err
	}

	c.Status = StatusClosed
	c.UpdatedAt = time.Now()
	c.ClosedAt = c.UpdatedAt

	return s.repo.Update(c)
}

// DeleteCase deletes a case
func (s *CaseService) DeleteCase(id string) error {
	return s.repo.Delete(id)
}

// ListCases lists cases with pagination
func (s *CaseService) ListCases(limit, offset int) ([]*Case, error) {
	return s.repo.List(limit, offset)
}

// Helper function to generate a simple ID
func generateID() string {
	return fmt.Sprintf("CASE-%d", time.Now().UnixNano())
}