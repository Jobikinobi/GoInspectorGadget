package document

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

//-----------------------------------------------------------------------------
// PDF Document Processor
//-----------------------------------------------------------------------------

// PDFProcessor implements DocumentProcessor for PDF files
type PDFProcessor struct {
	// Configuration
	PdfToTextPath string // Path to pdftotext executable
	UseOCR        bool   // Whether to use OCR for image-based PDFs
	TempDir       string // Directory for temporary files
}

// Process processes a PDF file and returns a Document
func (p *PDFProcessor) Process(filePath string) (*Document, error) {
	// Check if file exists and is PDF
	if !strings.HasSuffix(strings.ToLower(filePath), ".pdf") {
		return nil, fmt.Errorf("not a PDF file: %s", filePath)
	}

	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Extract text
	text, err := p.ExtractText(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text: %w", err)
	}

	// Extract metadata
	metadata, err := p.ExtractMetadata(filePath)
	if err != nil {
		// Don't fail completely on metadata extraction failure
		fmt.Printf("Warning: Failed to extract metadata from %s: %v\n", filePath, err)
	}

	// Create document
	doc := &Document{
		Title:       filepath.Base(filePath),
		Type:        TypeUnknown, // Will need to be determined based on content
		FilePath:    filePath,
		ContentType: "application/pdf",
		FileSize:    fileInfo.Size(),
		Content:     text,
		Metadata:    metadata,
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
	}

	// Try to infer document type from content
	doc.Type = inferDocumentType(doc)

	return doc, nil
}

// ExtractText extracts text from a PDF file
func (p *PDFProcessor) ExtractText(filePath string) (string, error) {
	if p.PdfToTextPath == "" {
		return "", fmt.Errorf("pdftotext not found, please install poppler-utils")
	}

	// Create temporary file for output
	outputFile := filepath.Join(p.TempDir, fmt.Sprintf("%d.txt", time.Now().UnixNano()))
	defer os.Remove(outputFile)

	// Run pdftotext command
	cmd := exec.Command(p.PdfToTextPath, "-layout", filePath, outputFile)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pdftotext failed: %w: %s", err, stderr.String())
	}

	// Read the output file
	content, err := os.ReadFile(outputFile)
	if err != nil {
		return "", fmt.Errorf("failed to read extracted text: %w", err)
	}

	return string(content), nil
}

// ExtractMetadata extracts metadata from a PDF file
func (p *PDFProcessor) ExtractMetadata(filePath string) (Metadata, error) {
	metadata := Metadata{
		CustomFields: make(map[string]string),
	}

	// Use pdfinfo to extract metadata
	pdfinfoPath, err := exec.LookPath("pdfinfo")
	if err != nil {
		return metadata, fmt.Errorf("pdfinfo not found, please install poppler-utils")
	}

	cmd := exec.Command(pdfinfoPath, filePath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return metadata, fmt.Errorf("pdfinfo failed: %w", err)
	}

	// Parse the output
	output := stdout.String()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "Title":
			metadata.Subject = value
		case "Author":
			metadata.Author = value
		case "Keywords":
			metadata.Keywords = strings.Split(value, ", ")
		case "CreationDate":
			// Parse date in format: Thu Jan 1 00:00:00 2000
			t, err := time.Parse("Mon Jan 2 15:04:05 2006", value)
			if err == nil {
				metadata.CreationDate = t
			}
		default:
			// Store all other metadata in custom fields
			metadata.CustomFields[key] = value
		}
	}

	return metadata, nil
}

// inferDocumentType attempts to determine the document type based on content
func inferDocumentType(doc *Document) DocumentType {
	text := strings.ToLower(doc.Content)

	// Police report indicators
	if containsAny(text, []string{
		"police report",
		"incident report",
		"offense report",
		"case report",
		"officer narrative",
		"reporting officer",
	}) {
		return TypePoliceReport
	}

	// Witness statement indicators
	if containsAny(text, []string{
		"witness statement",
		"statement of witness",
		"i, the undersigned",
		"do hereby state",
		"to the best of my recollection",
		"i witnessed",
		"i observed",
	}) {
		return TypeWitnessStatement
	}

	// Forensic report indicators
	if containsAny(text, []string{
		"forensic report",
		"laboratory report",
		"examination results",
		"dna analysis",
		"ballistics report",
		"toxicology report",
		"fingerprint analysis",
	}) {
		return TypeForensicReport
	}

	// Court filing indicators
	if containsAny(text, []string{
		"court of",
		"state vs",
		"plaintiff",
		"defendant",
		"motion to",
		"hereby ordered",
		"judge",
		"docket",
		"hearing",
		"trial",
	}) {
		return TypeCourtFiling
	}

	// Medical report indicators
	if containsAny(text, []string{
		"medical report",
		"patient name",
		"diagnosis",
		"treatment",
		"physician",
		"hospital",
		"medical record",
		"symptoms",
	}) {
		return TypeMedicalReport
	}

	// Default to unknown
	return TypeUnknown
}

// containsAny checks if the text contains any of the keywords
func containsAny(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}
