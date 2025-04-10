package casemanagement

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

// Case represents a police investigation case
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
	Victims          []Person
	Suspects         []Person
	Witnesses        []Person
	EvidenceIDs      []string // IDs of evidence items
	DocumentIDs      []string // IDs of documents
	InterviewIDs     []string // IDs of interviews
	Timeline         []Event  // Timeline of events
	Notes            []Note   // Investigator notes
	Tags             []string // Tags for categorization
	RelatedCases     []string // IDs of related cases
}

// Person represents an individual involved in a case
type Person struct {
	ID              string
	FullName        string
	DateOfBirth     time.Time
	Gender          string
	Address         string
	PhoneNumbers    []string
	EmailAddresses  []string
	Role            string // e.g., "Victim", "Suspect", "Witness"
	Description     string
	Notes           string
	Relationship    string // Relationship to case/other persons
	InterviewIDs    []string
	DocumentIDs     []string
	IsProtected     bool // Special privacy protection
	IsCooperative   bool
	HasPriorHistory bool
	PriorCases      []string // IDs of prior cases involving this person
}

// Event represents a single event in the case timeline
type Event struct {
	ID           string
	Timestamp    time.Time
	Description  string
	Location     string
	Participants []string // Person IDs
	DocumentIDs  []string
	EvidenceIDs  []string
	CreatedBy    string // User who created this event
	CreatedAt    time.Time
}

// Note represents an investigator's note
type Note struct {
	ID        string
	Title     string
	Content   string
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
	Tags      []string
	IsPrivate bool
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

	// Add a note about closure
	c.Notes = append(c.Notes, Note{
		ID:        generateID(),
		Title:     "Case Closure",
		Content:   fmt.Sprintf("Case closed. Reason: %s", reason),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	return s.repo.Update(c)
}

// AddPerson adds a person to a case
func (s *CaseService) AddPerson(caseID string, p Person) error {
	c, err := s.repo.Find(caseID)
	if err != nil {
		return err
	}

	if p.ID == "" {
		p.ID = generateID()
	}

	switch p.Role {
	case "Victim":
		c.Victims = append(c.Victims, p)
	case "Suspect":
		c.Suspects = append(c.Suspects, p)
	case "Witness":
		c.Witnesses = append(c.Witnesses, p)
	default:
		return fmt.Errorf("invalid person role: %s", p.Role)
	}

	c.UpdatedAt = time.Now()
	return s.repo.Update(c)
}

// AddEvent adds an event to a case timeline
func (s *CaseService) AddEvent(caseID string, event Event) error {
	c, err := s.repo.Find(caseID)
	if err != nil {
		return err
	}

	if event.ID == "" {
		event.ID = generateID()
	}

	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}

	c.Timeline = append(c.Timeline, event)
	c.UpdatedAt = time.Now()

	return s.repo.Update(c)
}

// AddNote adds a note to a case
func (s *CaseService) AddNote(caseID string, note Note) error {
	c, err := s.repo.Find(caseID)
	if err != nil {
		return err
	}

	if note.ID == "" {
		note.ID = generateID()
	}

	now := time.Now()
	note.CreatedAt = now
	note.UpdatedAt = now

	c.Notes = append(c.Notes, note)
	c.UpdatedAt = now

	return s.repo.Update(c)
}

// SearchCases searches for cases
func (s *CaseService) SearchCases(query string) ([]*Case, error) {
	return s.repo.Search(query)
}

// generateID generates a unique ID
func generateID() string {
	return fmt.Sprintf("CASE-%d", time.Now().UnixNano())
}
