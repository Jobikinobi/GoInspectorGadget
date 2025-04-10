package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// EvidenceType represents types of evidence
type EvidenceType string

const (
	TypePhysical   EvidenceType = "PHYSICAL"
	TypeDigital    EvidenceType = "DIGITAL"
	TypeDocument   EvidenceType = "DOCUMENT"
	TypeBiological EvidenceType = "BIOLOGICAL"
	TypeWeapon     EvidenceType = "WEAPON"
	TypeOther      EvidenceType = "OTHER"
)

// EvidenceStatus represents the status of evidence
type EvidenceStatus string

const (
	StatusCollected   EvidenceStatus = "COLLECTED"
	StatusProcessing  EvidenceStatus = "PROCESSING"
	StatusAnalyzed    EvidenceStatus = "ANALYZED"
	StatusInStorage   EvidenceStatus = "IN_STORAGE"
	StatusTransferred EvidenceStatus = "TRANSFERRED"
	StatusReleased    EvidenceStatus = "RELEASED"
	StatusDestroyed   EvidenceStatus = "DESTROYED"
)

// Evidence represents an item of evidence in a case
type Evidence struct {
	ID                string
	CaseID            string
	EvidenceNumber    string // Official evidence number
	Description       string
	Type              EvidenceType
	Status            EvidenceStatus
	CollectedBy       string // ID of the person who collected the evidence
	CollectionDate    time.Time
	CollectionMethod  string
	CollectionNotes   string
	CollectionDetails *CollectionDetails
	Location          Location // Where the evidence was found/collected
	StorageLocation   string   // Current storage location
	ChainOfCustody    []CustodyEvent
	Tags              []string
	RelatedEvidence   []string // IDs of related evidence items
	ImagePaths        []string // Paths to images of the evidence
	FileHash          string   // For digital evidence, hash of the file
	IsConfidential    bool
	Notes             string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// CollectionDetails contains specific details about evidence collection
type CollectionDetails struct {
	Weather            string
	LightingCondition  string
	Temperature        string
	Environment        string
	Packaging          string
	ContainedWithin    string // If evidence was contained within something
	PreservationMethod string
}

// Location represents a physical location where evidence was found
type Location struct {
	Description    string
	Address        string
	Room           string
	GPS            string // GPS coordinates
	MapReference   string
	LocationNotes  string
	PhotoReference string // Reference to photos of the location
}

// CustodyEvent represents a single event in the chain of custody
type CustodyEvent struct {
	ID                 string
	EvidenceID         string
	Timestamp          time.Time
	Action             string // e.g., "COLLECTED", "TRANSFERRED", "ANALYZED", "STORED"
	FromPerson         string // ID of person transferring
	ToPerson           string // ID of person receiving
	FromLocation       string
	ToLocation         string
	Reason             string
	Notes              string
	DocumentID         string // ID of associated document (if any)
	AuthorizedBy       string // ID of authorizing person
	TransportMethod    string
	VerificationMethod string // How the evidence was verified during transfer
}

// DigitalEvidence contains additional fields for digital evidence
type DigitalEvidence struct {
	Evidence
	FileType         string
	FilePath         string
	FileSize         int64
	CreationDate     time.Time
	ModifiedDate     time.Time
	DeviceSource     string // Device the evidence came from
	OriginalHash     string // Original hash of the file
	WorkingCopyHash  string
	ExtractionMethod string
	Encrypted        bool
	Decrypted        bool
	Password         string // This would be securely stored
	Metadata         map[string]string
}

// BiologicalEvidence contains additional fields for biological evidence
type BiologicalEvidence struct {
	Evidence
	BiologicalType     string // e.g., "BLOOD", "DNA", "TISSUE"
	SampleID           string
	PreservationMethod string
	AnalysisResults    string
	ContainerType      string
	StorageConditions  string
	ExpirationDate     time.Time
}

// EvidenceRepository defines the interface for evidence storage
type EvidenceRepository interface {
	Save(e *Evidence) error
	Find(id string) (*Evidence, error)
	FindByCase(caseID string) ([]*Evidence, error)
	Search(query string) ([]*Evidence, error)
	Update(e *Evidence) error
	Delete(id string) error
}

// EvidenceService provides business logic for evidence management
type EvidenceService struct {
	repo EvidenceRepository
}

// NewEvidenceService creates a new evidence service
func NewEvidenceService(repo EvidenceRepository) *EvidenceService {
	return &EvidenceService{
		repo: repo,
	}
}

// CreateEvidence creates a new evidence item
func (s *EvidenceService) CreateEvidence(e *Evidence) error {
	if e.ID == "" {
		e.ID = generateID("EV")
	}

	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now

	if e.Status == "" {
		e.Status = StatusCollected
	}

	// If this is digital evidence, calculate file hash
	if e.Type == TypeDigital && e.FileHash == "" {
		// Get type assertion to access specific fields
		if de, ok := interface{}(e).(DigitalEvidence); ok {
			hash, err := calculateFileHash(de.FilePath)
			if err != nil {
				return fmt.Errorf("failed to calculate file hash: %w", err)
			}
			e.FileHash = hash
		}
	}

	// Initialize chain of custody with collection event
	if len(e.ChainOfCustody) == 0 {
		e.ChainOfCustody = []CustodyEvent{
			{
				ID:           generateID("CE"),
				EvidenceID:   e.ID,
				Timestamp:    e.CollectionDate,
				Action:       "COLLECTED",
				FromPerson:   "",
				ToPerson:     e.CollectedBy,
				FromLocation: fmt.Sprintf("%s, %s", e.Location.Description, e.Location.Address),
				ToLocation:   e.StorageLocation,
				Reason:       "Initial collection",
				Notes:        e.CollectionNotes,
			},
		}
	}

	return s.repo.Save(e)
}

// GetEvidence retrieves an evidence item by ID
func (s *EvidenceService) GetEvidence(id string) (*Evidence, error) {
	return s.repo.Find(id)
}

// UpdateEvidence updates an existing evidence item
func (s *EvidenceService) UpdateEvidence(e *Evidence) error {
	e.UpdatedAt = time.Now()
	return s.repo.Update(e)
}

// TransferCustody records a transfer in the chain of custody
func (s *EvidenceService) TransferCustody(
	evidenceID, fromPerson, toPerson, fromLocation, toLocation, reason, notes string,
) error {
	evidence, err := s.repo.Find(evidenceID)
	if err != nil {
		return fmt.Errorf("failed to find evidence: %w", err)
	}

	// Create custody event
	event := CustodyEvent{
		ID:           generateID("CE"),
		EvidenceID:   evidenceID,
		Timestamp:    time.Now(),
		Action:       "TRANSFERRED",
		FromPerson:   fromPerson,
		ToPerson:     toPerson,
		FromLocation: fromLocation,
		ToLocation:   toLocation,
		Reason:       reason,
		Notes:        notes,
	}

	// Add to chain of custody
	evidence.ChainOfCustody = append(evidence.ChainOfCustody, event)

	// Update storage location
	evidence.StorageLocation = toLocation
	evidence.UpdatedAt = time.Now()

	return s.repo.Update(evidence)
}

// SearchEvidence searches for evidence
func (s *EvidenceService) SearchEvidence(query string) ([]*Evidence, error) {
	return s.repo.Search(query)
}

// VerifyIntegrity checks if a digital evidence file is intact
func (s *EvidenceService) VerifyIntegrity(evidenceID string) (bool, error) {
	evidence, err := s.repo.Find(evidenceID)
	if err != nil {
		return false, fmt.Errorf("failed to find evidence: %w", err)
	}

	if evidence.Type != TypeDigital {
		return false, fmt.Errorf("integrity verification is only applicable to digital evidence")
	}

	// Get digital evidence specific fields
	if de, ok := interface{}(evidence).(DigitalEvidence); ok {
		// Calculate current hash
		currentHash, err := calculateFileHash(de.FilePath)
		if err != nil {
			return false, fmt.Errorf("failed to calculate current file hash: %w", err)
		}

		// Compare with stored hash
		return currentHash == evidence.FileHash, nil
	}

	return false, fmt.Errorf("failed to convert to digital evidence type")
}

// CreateDigitalEvidence creates a new digital evidence item with file validation
func (s *EvidenceService) CreateDigitalEvidence(e *DigitalEvidence) error {
	// Get file info
	fileInfo, err := os.Stat(e.FilePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Set file size and type
	e.FileSize = fileInfo.Size()
	e.FileType = filepath.Ext(e.FilePath)

	// Calculate hash
	hash, err := calculateFileHash(e.FilePath)
	if err != nil {
		return fmt.Errorf("failed to calculate file hash: %w", err)
	}
	e.FileHash = hash
	e.OriginalHash = hash

	// Set as digital evidence type
	e.Type = TypeDigital

	// Create evidence
	return s.CreateEvidence(&e.Evidence)
}

// calculateFileHash calculates the SHA-256 hash of a file
func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// generateID generates a unique ID with a prefix
func generateID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}
