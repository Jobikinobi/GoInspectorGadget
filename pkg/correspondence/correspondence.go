package correspondence

import (
	"fmt"
	"time"
)

// CorrespondenceType defines the type of correspondence
type CorrespondenceType string

const (
	TypeEmail        CorrespondenceType = "EMAIL"
	TypeLetter       CorrespondenceType = "LETTER"
	TypeFax          CorrespondenceType = "FAX"
	TypeMemo         CorrespondenceType = "MEMO"
	TypeNotice       CorrespondenceType = "NOTICE"
	TypeSubpoena     CorrespondenceType = "SUBPOENA"
	TypeWarrant      CorrespondenceType = "WARRANT"
	TypeEvidenceReq  CorrespondenceType = "EVIDENCE_REQUEST"
	TypeReport       CorrespondenceType = "REPORT"
	TypeCourtFiling  CorrespondenceType = "COURT_FILING"
	TypeSMS          CorrespondenceType = "SMS"
	TypePressRelease CorrespondenceType = "PRESS_RELEASE"
)

// Status defines the status of correspondence
type Status string

const (
	StatusDraft     Status = "DRAFT"
	StatusSent      Status = "SENT"
	StatusDelivered Status = "DELIVERED"
	StatusRead      Status = "READ"
	StatusResponded Status = "RESPONDED"
	StatusFailed    Status = "FAILED"
	StatusCancelled Status = "CANCELLED"
	StatusPending   Status = "PENDING_APPROVAL"
	StatusApproved  Status = "APPROVED"
	StatusRejected  Status = "REJECTED"
)

// Priority defines the priority level of correspondence
type Priority string

const (
	PriorityLow    Priority = "LOW"
	PriorityNormal Priority = "NORMAL"
	PriorityHigh   Priority = "HIGH"
	PriorityUrgent Priority = "URGENT"
)

// Correspondence represents a communication record
type Correspondence struct {
	ID                 string
	CaseID             string
	CorrespondenceType CorrespondenceType
	Subject            string
	Body               string
	Sender             Person
	Recipients         []Person
	SentAt             time.Time
	ReceivedAt         time.Time
	Direction          string // "OUTGOING" or "INCOMING"
	ReferenceNumber    string
	Priority           Priority
	Status             Status
	Attachments        []Attachment
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// Person represents an individual involved in correspondence
type Person struct {
	ID           string
	Name         string
	Title        string
	Organization string
	Department   string
	Address      string
	Phone        string
	Email        string
	IsOfficer    bool
	BadgeNumber  string
}

// Attachment represents a file attached to correspondence
type Attachment struct {
	ID          string
	Name        string
	ContentType string
	FilePath    string
	Size        int64
	CreatedAt   time.Time
}

// Template represents a pre-defined correspondence template
type Template struct {
	ID           string
	Name         string
	Type         CorrespondenceType
	Subject      string
	Body         string
	TemplateVars []string // Variables that can be replaced in the template
	Department   string   // Department that owns the template
	IsApproved   bool
	ApprovedBy   string
	ApprovedAt   time.Time
	CreatedAt    time.Time
	UpdatedBy    string
	UpdatedAt    time.Time
}

// CorrespondenceRepository defines the interface for correspondence storage
type CorrespondenceRepository interface {
	Save(c *Correspondence) error
	Find(id string) (*Correspondence, error)
	FindByCase(caseID string) ([]*Correspondence, error)
	FindByType(correspondenceType CorrespondenceType) ([]*Correspondence, error)
	FindByStatus(status Status) ([]*Correspondence, error)
	FindByReference(refNumber string) (*Correspondence, error)
	Search(query string) ([]*Correspondence, error)
	Update(c *Correspondence) error
	Delete(id string) error
}

// TemplateRepository defines the interface for template storage
type TemplateRepository interface {
	Save(t *Template) error
	Find(id string) (*Template, error)
	FindByName(name string) (*Template, error)
	FindByType(correspondenceType CorrespondenceType) ([]*Template, error)
	FindByDepartment(department string) ([]*Template, error)
	Update(t *Template) error
	Delete(id string) error
}

// CorrespondenceService provides business logic for correspondence management
type CorrespondenceService struct {
	correspondenceRepo CorrespondenceRepository
	templateRepo       TemplateRepository
}

// NewCorrespondenceService creates a new correspondence service
func NewCorrespondenceService(correspondenceRepo CorrespondenceRepository, templateRepo TemplateRepository) *CorrespondenceService {
	return &CorrespondenceService{
		correspondenceRepo: correspondenceRepo,
		templateRepo:       templateRepo,
	}
}

// CreateCorrespondence creates a new correspondence
func (s *CorrespondenceService) CreateCorrespondence(c *Correspondence) error {
	if c.ID == "" {
		c.ID = generateID("CORR")
	}

	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	return s.correspondenceRepo.Save(c)
}

// GetCorrespondence retrieves a correspondence by ID
func (s *CorrespondenceService) GetCorrespondence(id string) (*Correspondence, error) {
	return s.correspondenceRepo.Find(id)
}

// CreateFromTemplate creates a correspondence from a template
func (s *CorrespondenceService) CreateFromTemplate(
	templateID string,
	caseID string,
	sender Person,
	recipients []Person,
) (*Correspondence, error) {
	template, err := s.templateRepo.Find(templateID)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	corr := &Correspondence{
		ID:                 generateID("CORR"),
		CaseID:             caseID,
		CorrespondenceType: template.Type,
		Subject:            template.Subject,
		Body:               template.Body,
		Sender:             sender,
		Recipients:         recipients,
		Direction:          "OUTGOING",
		Priority:           PriorityNormal,
		Status:             StatusDraft,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	return corr, s.correspondenceRepo.Save(corr)
}

// SendCorrespondence marks a correspondence as sent
func (s *CorrespondenceService) SendCorrespondence(id string, sentAt time.Time) error {
	corr, err := s.correspondenceRepo.Find(id)
	if err != nil {
		return err
	}

	if corr.Status != StatusDraft && corr.Status != StatusPending {
		return fmt.Errorf("cannot send correspondence with status %s", corr.Status)
	}

	corr.Status = StatusSent
	corr.SentAt = sentAt
	corr.UpdatedAt = time.Now()

	return s.correspondenceRepo.Update(corr)
}

// generateID generates a unique ID with a prefix
func generateID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}
