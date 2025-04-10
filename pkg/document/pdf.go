package document

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// DocumentType represents the type of document
type DocumentType int

const (
	TypeUnknown DocumentType = iota
	TypePoliceReport
	TypeWitnessStatement
	TypeForensicReport
	TypeCourtFiling
	TypeMedicalReport
	TypePersonalIdentification
	TypeEvidence
	TypeTranscript
)

// Metadata contains document metadata
type Metadata struct {
	Author       string
	Subject      string
	Keywords     []string
	CreationDate time.Time
	ModifiedDate time.Time
	CustomFields map[string]string
}

// Document represents a document in the system
type Document struct {
	ID          string
	Title       string
	Type        DocumentType
	CaseID      string
	FilePath    string
	ContentType string
	FileSize    int64
	Hash        string
	Content     string
	Metadata    Metadata
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// DocumentProcessor defines the interface for document processors
type DocumentProcessor interface {
	Process(filePath string) (*Document, error)
}

// ImportDocument imports a document from a file
func ImportDocument(filePath, destDir string, processor DocumentProcessor) (*Document, error) {
	// Validate file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}

	// Process document using the provided processor
	doc, err := processor.Process(filePath)
	if err != nil {
		return nil, err
	}

	// Set creation time
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = doc.CreatedAt

	// Generate ID if not set
	if doc.ID == "" {
		doc.ID = fmt.Sprintf("DOC-%d", time.Now().UnixNano())
	}

	return doc, nil
}

// GetDocumentTypeString returns a string representation of DocumentType
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
	case TypeMedicalReport:
		return "Medical Report"
	case TypePersonalIdentification:
		return "Personal Identification"
	case TypeEvidence:
		return "Evidence Document"
	case TypeTranscript:
		return "Transcript"
	default:
		return "Unknown"
	}
}

// PDFProcessor implements DocumentProcessor for PDF files
type PDFProcessor struct {
	// Configuration
	PdfToTextPath string // Path to pdftotext executable
	UseOCR        bool   // Whether to use OCR for image-based PDFs
	TempDir       string // Directory for temporary files
}

// NewPDFProcessor creates a new PDFProcessor
func NewPDFProcessor(pdfToTextPath string, useOCR bool, tempDir string) *PDFProcessor {
	if pdfToTextPath == "" {
		// Try to find pdftotext in PATH
		path, err := exec.LookPath("pdftotext")
		if err == nil {
			pdfToTextPath = path
		}
	}

	return &PDFProcessor{
		PdfToTextPath: pdfToTextPath,
		UseOCR:        useOCR,
		TempDir:       tempDir,
	}
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
