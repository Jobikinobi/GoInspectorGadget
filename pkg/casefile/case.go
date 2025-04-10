package casefile

import (
	"fmt"
	"time"
)

// Status represents the current status of a case
type Status string

const (
	StatusOpen       Status = "OPEN"
	StatusClosed     Status = "CLOSED"
	StatusSuspended  Status = "SUSPENDED"
	StatusInactive   Status = "INACTIVE"
	StatusReferred   Status = "REFERRED"
	StatusProsecuted Status = "PROSECUTED"
)

// Priority level of a case
type Priority int

const (
	PriorityLow Priority = iota + 1
	PriorityMedium
	PriorityHigh
	PriorityCritical
)

// Case represents a police investigation case file
type Case struct {
	ID               string
	CaseNumber       string // Official case number
	Title            string
	Description      string
	Status           Status
	Priority         Priority
	CaseType         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	AssignedTo       []string // IDs of investigators assigned to the case
	LeadInvestigator string
	Jurisdiction     string
	Location         string
	IncidentDate     time.Time
	ReportDate       time.Time
	Documents        []Document
}

// Document represents a document in a case file
type Document struct {
	ID          string
	Title       string
	Type        string
	Description string
	FilePath    string
	FileSize    int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CaseRepository defines the interface for case storage
type CaseRepository interface {
	Save(c *Case) error
	Find(id string) (*Case, error)
	FindByCaseNumber(caseNumber string) (*Case, error)
	Search(query string) ([]*Case, error)
	List(limit, offset int) ([]*Case, error)
	Update(c *Case) error
	Delete(id string) error
}

// CaseService provides business logic for case management
type CaseService struct {
	repo CaseRepository
}

// NewCaseService creates a new case service
func NewCaseService(repo CaseRepository) *CaseService {
	return &CaseService{repo: repo}
}

// CreateCase creates a new case
func (s *CaseService) CreateCase(c *Case) error {
	if c.ID == "" {
		c.ID = generateID()
	}

	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	if c.Status == "" {
		c.Status = StatusOpen
	}

	return s.repo.Save(c)
}

// GetCase retrieves a case by ID
func (s *CaseService) GetCase(id string) (*Case, error) {
	return s.repo.Find(id)
}

// UpdateCase updates an existing case
func (s *CaseService) UpdateCase(c *Case) error {
	c.UpdatedAt = time.Now()
	return s.repo.Update(c)
}

// CloseCase closes a case
func (s *CaseService) CloseCase(id string, reason string) error {
	c, err := s.repo.Find(id)
	if err != nil {
		return err
	}

	c.Status = StatusClosed
	c.UpdatedAt = time.Now()

	return s.repo.Update(c)
}

// generateID generates a unique ID
func generateID() string {
	return fmt.Sprintf("CF-%d", time.Now().UnixNano())
}
