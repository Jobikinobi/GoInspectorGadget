package document

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// DocumentType defines the type of document
type DocumentType int

const (
	// Document types relevant to police investigation
	TypeUnknown DocumentType = iota
	TypePoliceReport
	TypeWitnessStatement
	TypeForensicReport
	TypeCourtFiling
	TypeEvidence
	TypeMedicalReport
	TypeTranscript
	TypePersonalIdentification
	TypeBackground
	TypeNote
	TypeForensicAnalysis
	TypeCourtDocument
	TypeMedicalRecord
	TypeEvidenceItem
	TypeTranscriptRecord
)

// Document represents any document in the investigation system
type Document struct {
	ID             string
	Title          string
	Type           DocumentType
	FilePath       string
	ContentType    string
	FileSize       int64
	CreatedAt      time.Time
	ModifiedAt     time.Time
	Content        string      // Plain text content (if available)
	Metadata       Metadata    // Document metadata
	CaseID         string      // ID of the case this document belongs to
	Tags           []string    // User-defined tags
	Redactions     []Redaction // Any redacted portions
	Annotations    []Annotation
	IsConfidential bool
}

// Metadata contains document metadata
type Metadata struct {
	Author       string
	CreationDate time.Time
	Subject      string
	Keywords     []string
	Source       string
	CustomFields map[string]string
}

// Redaction represents a redacted portion of a document
type Redaction struct {
	StartPos    int
	EndPos      int
	Reason      string
	RedactedBy  string
	RedactedAt  time.Time
	IsTemporary bool
}

// Annotation represents a note or comment on a document
type Annotation struct {
	ID        string
	UserID    string
	Text      string
	CreatedAt time.Time
	Position  int // Position in the document where annotation is attached
	IsPrivate bool
}

// DocumentProcessor defines the interface for processing different document types
type DocumentProcessor interface {
	Process(filePath string) (*Document, error)
}

// DocumentRepository defines the interface for document storage
type DocumentRepository interface {
	Save(doc *Document) error
	Find(id string) (*Document, error)
	FindByCase(caseID string) ([]*Document, error)
	FindByType(docType DocumentType) ([]*Document, error)
	Search(query string) ([]*Document, error)
	Delete(id string) error
	Update(doc *Document) error
}

// ImportDocument imports a document into the system
func ImportDocument(filePath string, targetDir string, processor DocumentProcessor) (*Document, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	// Process document to extract metadata and content
	doc, err := processor.Process(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to process document: %w", err)
	}

	// Generate a unique ID for the document if not set
	if doc.ID == "" {
		doc.ID = generateID()
	}

	// Copy file to target directory
	if targetDir != "" {
		destPath := filepath.Join(targetDir, filepath.Base(filePath))
		if err := copyFile(filePath, destPath); err != nil {
			return nil, fmt.Errorf("failed to copy file: %w", err)
		}
		doc.FilePath = destPath
	} else {
		doc.FilePath = filePath
	}

	// Set timestamps
	now := time.Now()
	doc.CreatedAt = now
	doc.ModifiedAt = now

	return doc, nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy content
	_, err = io.Copy(destFile, sourceFile)
	return err
}

// generateID generates a unique ID for a document
func generateID() string {
	// Simple implementation - would use UUID in production
	return fmt.Sprintf("DOC-%d", time.Now().UnixNano())
}

// GetDocumentTypeString returns a string representation of a document type
func GetDocumentTypeString(docType DocumentType) string {
	switch docType {
	case TypePoliceReport:
		return "Police Report"
	case TypeWitnessStatement:
		return "Witness Statement"
	case TypeForensicReport:
		return "Forensic Report"
	case TypeCourtFiling:
		return "Court Filing"
	case TypeEvidence:
		return "Evidence Record"
	case TypeMedicalReport:
		return "Medical Report"
	case TypeTranscript:
		return "Transcript"
	case TypePersonalIdentification:
		return "Personal Identification"
	case TypeBackground:
		return "Background Information"
	case TypeNote:
		return "Investigator Note"
	case TypeForensicAnalysis:
		return "Forensic Analysis"
	case TypeCourtDocument:
		return "Court Document"
	case TypeMedicalRecord:
		return "Medical Record"
	case TypeEvidenceItem:
		return "Evidence Item"
	case TypeTranscriptRecord:
		return "Transcript Record"
	default:
		return "Unknown Document"
	}
}
