package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jth/claude/GoCode/pkg/casefile"
	"github.com/jth/claude/GoCode/pkg/casemanagement"
	"github.com/jth/claude/GoCode/pkg/correspondence"
	"github.com/jth/claude/GoCode/pkg/document"
	"github.com/jth/claude/GoCode/pkg/evidence"
	"github.com/jth/claude/GoCode/pkg/interview"
)

// Simple in-memory repositories for demonstration
type inMemoryRepo struct {
	caseRecords    map[string]*casemanagement.Case
	casefiles      map[string]*casefile.Case
	documents      map[string]*document.Document
	evidence       map[string]*evidence.Evidence
	interviews     map[string]*interview.Interview
	transcripts    map[string]*interview.Transcript
	correspondence map[string]*correspondence.Correspondence
	templates      map[string]*correspondence.Template
}

// CLI application state
type InvestigatorApp struct {
	// Current state
	workingDir    string
	currentCaseID string

	// Services
	caseService           *casemanagement.CaseService
	casefileService       *casefile.CaseService
	documentService       document.DocumentProcessor
	evidenceService       *evidence.EvidenceService
	interviewService      *interview.InterviewService
	correspondenceService *correspondence.CorrespondenceService

	// Repositories
	repo *inMemoryRepo
}

func NewInvestigatorApp(workingDir string) *InvestigatorApp {
	repo := &inMemoryRepo{
		caseRecords:    make(map[string]*casemanagement.Case),
		casefiles:      make(map[string]*casefile.Case),
		documents:      make(map[string]*document.Document),
		evidence:       make(map[string]*evidence.Evidence),
		interviews:     make(map[string]*interview.Interview),
		transcripts:    make(map[string]*interview.Transcript),
		correspondence: make(map[string]*correspondence.Correspondence),
		templates:      make(map[string]*correspondence.Template),
	}

	app := &InvestigatorApp{
		workingDir: workingDir,
		repo:       repo,
	}

	// Initialize services
	app.initializeServices()

	return app
}

func (app *InvestigatorApp) initializeServices() {
	// Initialize case repository implementation
	caseRepo := &inMemoryCaseRepo{cases: app.repo.caseRecords}
	app.caseService = casemanagement.NewCaseService(caseRepo)

	// Initialize casefile repository implementation
	casefileRepo := &inMemoryCasefileRepo{cases: app.repo.casefiles}
	app.casefileService = casefile.NewCaseService(casefileRepo)

	// Initialize document service
	tempDir := filepath.Join(app.workingDir, "temp")
	os.MkdirAll(tempDir, 0755)
	// Create a new PDF processor for document handling
	pdfProcessor := &document.PDFProcessor{
		PdfToTextPath: "",
		UseOCR:        false,
		TempDir:       tempDir,
	}
	app.documentService = pdfProcessor

	// Initialize evidence repository implementation
	evidenceRepo := &inMemoryEvidenceRepo{evidence: app.repo.evidence}
	app.evidenceService = evidence.NewEvidenceService(evidenceRepo)

	// Initialize interview repository implementations
	interviewRepo := &inMemoryInterviewRepo{interviews: app.repo.interviews}
	transcriptRepo := &inMemoryTranscriptRepo{transcripts: app.repo.transcripts}

	// Create a simple speech recognizer (would be replaced with real implementation)
	recognizer := &dummySpeechRecognizer{}

	app.interviewService = interview.NewInterviewService(interviewRepo, transcriptRepo, recognizer)

	// Initialize correspondence repository implementations
	correspondenceRepo := &inMemoryCorrespondenceRepo{correspondence: app.repo.correspondence}
	templateRepo := &inMemoryTemplateRepo{templates: app.repo.templates}

	app.correspondenceService = correspondence.NewCorrespondenceService(correspondenceRepo, templateRepo)

	// Initialize default templates in the repository
	for _, t := range correspondence.GetDefaultTemplates() {
		app.repo.templates[t.ID] = t
	}
}

func main() {
	// Set up command line flags
	// Case subcommands
	caseCreateCmd := flag.NewFlagSet("case create", flag.ExitOnError)
	caseOpenCmd := flag.NewFlagSet("case open", flag.ExitOnError)
	caseListCmd := flag.NewFlagSet("case list", flag.ExitOnError)

	// Case create flags
	caseTitle := caseCreateCmd.String("title", "", "Case title")
	caseDesc := caseCreateCmd.String("desc", "", "Case description")
	caseType := caseCreateCmd.String("type", "", "Case type")

	// Document subcommands
	docImportCmd := flag.NewFlagSet("doc import", flag.ExitOnError)

	// Document import flags
	docPath := docImportCmd.String("path", "", "Path to document file")
	docCase := docImportCmd.String("case", "", "Case ID to associate document with")

	// Evidence subcommands
	evidenceAddCmd := flag.NewFlagSet("evidence add", flag.ExitOnError)
	evidenceListCmd := flag.NewFlagSet("evidence list", flag.ExitOnError)

	// Evidence add flags
	evidenceDesc := evidenceAddCmd.String("desc", "", "Evidence description")
	evidenceType := evidenceAddCmd.String("type", "PHYSICAL", "Evidence type (PHYSICAL, DIGITAL, etc.)")
	evidenceCase := evidenceAddCmd.String("case", "", "Case ID to associate evidence with")

	// Interview subcommands
	interviewAddCmd := flag.NewFlagSet("interview add", flag.ExitOnError)
	interviewTranscribeCmd := flag.NewFlagSet("interview transcribe", flag.ExitOnError)

	// Interview add flags
	interviewTitle := interviewAddCmd.String("title", "", "Interview title")
	interviewCase := interviewAddCmd.String("case", "", "Case ID to associate interview with")
	interviewType := interviewAddCmd.String("type", "WITNESS", "Interview type (WITNESS, SUSPECT, etc.)")

	// Interview transcribe flags
	interviewID := interviewTranscribeCmd.String("id", "", "Interview ID to transcribe")

	// Correspondence subcommands
	corrCreateCmd := flag.NewFlagSet("correspondence create", flag.ExitOnError)
	corrListCmd := flag.NewFlagSet("correspondence list", flag.ExitOnError)
	corrSendCmd := flag.NewFlagSet("correspondence send", flag.ExitOnError)
	corrTemplateListCmd := flag.NewFlagSet("correspondence templates", flag.ExitOnError)

	// Correspondence create flags
	corrType := corrCreateCmd.String("type", "", "Correspondence type (EMAIL, LETTER, etc.)")
	corrSubject := corrCreateCmd.String("subject", "", "Subject")
	corrBody := corrCreateCmd.String("body", "", "Body content")
	corrRecipient := corrCreateCmd.String("recipient", "", "Recipient name")
	corrCase := corrCreateCmd.String("case", "", "Case ID")
	corrTemplate := corrCreateCmd.String("template", "", "Template ID to use")

	// Correspondence send flags
	corrID := corrSendCmd.String("id", "", "Correspondence ID to send")

	// Create the application with working directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	appDir := filepath.Join(homeDir, "investigator-simulator")
	os.MkdirAll(appDir, 0755)

	app := NewInvestigatorApp(appDir)

	// Check command line arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Process commands
	switch os.Args[1] {
	case "case":
		if len(os.Args) < 3 {
			fmt.Println("Missing case subcommand")
			os.Exit(1)
		}

		switch os.Args[2] {
		case "create":
			caseCreateCmd.Parse(os.Args[3:])
			app.handleCaseCreate(*caseTitle, *caseDesc, *caseType)

		case "open":
			caseOpenCmd.Parse(os.Args[3:])
			if caseOpenCmd.NArg() < 1 {
				fmt.Println("Missing case ID")
				os.Exit(1)
			}
			app.handleCaseOpen(caseOpenCmd.Arg(0))

		case "list":
			caseListCmd.Parse(os.Args[3:])
			app.handleCaseList()

		default:
			fmt.Printf("Unknown case subcommand: %s\n", os.Args[2])
			os.Exit(1)
		}

	case "doc":
		if len(os.Args) < 3 {
			fmt.Println("Missing document subcommand")
			os.Exit(1)
		}

		switch os.Args[2] {
		case "import":
			docImportCmd.Parse(os.Args[3:])
			app.handleDocImport(*docPath, *docCase)

		default:
			fmt.Printf("Unknown document subcommand: %s\n", os.Args[2])
			os.Exit(1)
		}

	case "evidence":
		if len(os.Args) < 3 {
			fmt.Println("Missing evidence subcommand")
			os.Exit(1)
		}

		switch os.Args[2] {
		case "add":
			evidenceAddCmd.Parse(os.Args[3:])
			app.handleEvidenceAdd(*evidenceDesc, *evidenceType, *evidenceCase)

		case "list":
			evidenceListCmd.Parse(os.Args[3:])
			if evidenceListCmd.NArg() > 0 {
				app.handleEvidenceList(evidenceListCmd.Arg(0))
			} else {
				app.handleEvidenceList("")
			}

		default:
			fmt.Printf("Unknown evidence subcommand: %s\n", os.Args[2])
			os.Exit(1)
		}

	case "interview":
		if len(os.Args) < 3 {
			fmt.Println("Missing interview subcommand")
			os.Exit(1)
		}

		switch os.Args[2] {
		case "add":
			interviewAddCmd.Parse(os.Args[3:])
			app.handleInterviewAdd(*interviewTitle, *interviewType, *interviewCase)

		case "transcribe":
			interviewTranscribeCmd.Parse(os.Args[3:])
			app.handleInterviewTranscribe(*interviewID)

		default:
			fmt.Printf("Unknown interview subcommand: %s\n", os.Args[2])
			os.Exit(1)
		}

	case "correspondence":
		if len(os.Args) < 3 {
			fmt.Println("Missing correspondence subcommand")
			os.Exit(1)
		}

		switch os.Args[2] {
		case "create":
			corrCreateCmd.Parse(os.Args[3:])
			app.handleCorrespondenceCreate(*corrType, *corrSubject, *corrBody, *corrRecipient, *corrCase, *corrTemplate)

		case "list":
			corrListCmd.Parse(os.Args[3:])
			if corrListCmd.NArg() > 0 {
				app.handleCorrespondenceList(corrListCmd.Arg(0))
			} else {
				app.handleCorrespondenceList("")
			}

		case "send":
			corrSendCmd.Parse(os.Args[3:])
			app.handleCorrespondenceSend(*corrID)

		case "templates":
			corrTemplateListCmd.Parse(os.Args[3:])
			app.handleCorrespondenceTemplateList()

		default:
			fmt.Printf("Unknown correspondence subcommand: %s\n", os.Args[2])
			os.Exit(1)
		}

	case "help":
		printUsage()

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Police Investigator Simulator")
	fmt.Println("Usage:")
	fmt.Println("  investigator case create --title \"Title\" --desc \"Description\" --type \"Homicide\"")
	fmt.Println("  investigator case open <case-id>")
	fmt.Println("  investigator case list")
	fmt.Println("  investigator doc import --path \"path/to/file.pdf\" --case <case-id>")
	fmt.Println("  investigator evidence add --desc \"Description\" --type \"PHYSICAL\" --case <case-id>")
	fmt.Println("  investigator evidence list [case-id]")
	fmt.Println("  investigator interview add --title \"Interview\" --type \"WITNESS\" --case <case-id>")
	fmt.Println("  investigator interview transcribe --id <interview-id>")
	fmt.Println("  investigator correspondence create --type \"EMAIL\" --subject \"Subject\" --recipient \"Name\" --case <case-id>")
	fmt.Println("  investigator correspondence create --template <template-id> --recipient \"Name\" --case <case-id>")
	fmt.Println("  investigator correspondence list [case-id]")
	fmt.Println("  investigator correspondence send --id <correspondence-id>")
	fmt.Println("  investigator correspondence templates")
}

// Command handlers
func (app *InvestigatorApp) handleCaseCreate(title, description, caseType string) {
	if title == "" {
		fmt.Println("Error: Case title is required")
		os.Exit(1)
	}

	c := &casemanagement.Case{
		Title:       title,
		Description: description,
		CaseType:    caseType,
		Priority:    casemanagement.PriorityMedium,
		Status:      casemanagement.StatusOpen,
	}

	err := app.caseService.CreateCase(c)
	if err != nil {
		fmt.Printf("Error creating case: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Case created successfully with ID: %s\n", c.ID)
	app.currentCaseID = c.ID
}

func (app *InvestigatorApp) handleCaseOpen(caseID string) {
	c, err := app.caseService.GetCase(caseID)
	if err != nil {
		fmt.Printf("Error opening case: %v\n", err)
		os.Exit(1)
	}

	app.currentCaseID = caseID
	fmt.Printf("Opened case: %s - %s\n", c.ID, c.Title)
	fmt.Printf("Status: %s, Type: %s\n", c.Status, c.CaseType)
	fmt.Printf("Created: %s\n", c.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Description: %s\n", c.Description)
}

func (app *InvestigatorApp) handleCaseList() {
	// In a real implementation, this would query the repository
	// Here we just list all cases in the in-memory repo
	if len(app.repo.caseRecords) == 0 {
		fmt.Println("No cases found")
		return
	}

	fmt.Println("\nCase List:")
	fmt.Println("-------------------------------------------------")
	fmt.Println("ID\t\tTitle\t\tStatus\tDate")
	fmt.Println("-------------------------------------------------")

	for _, c := range app.repo.caseRecords {
		fmt.Printf("%s\t%s\t%s\t%s\n",
			c.ID,
			c.Title,
			c.Status,
			c.CreatedAt.Format("2006-01-02"))
	}
}

func (app *InvestigatorApp) handleDocImport(path, caseID string) {
	if path == "" {
		fmt.Println("Error: Document path is required")
		os.Exit(1)
	}

	if caseID == "" {
		if app.currentCaseID == "" {
			fmt.Println("Error: No case specified and no case is currently open")
			os.Exit(1)
		}
		caseID = app.currentCaseID
	}

	// Ensure the case exists
	_, err := app.caseService.GetCase(caseID)
	if err != nil {
		fmt.Printf("Error: Case not found: %v\n", err)
		os.Exit(1)
	}

	// Process the document
	doc, err := document.ImportDocument(path, filepath.Join(app.workingDir, "documents"), app.documentService)
	if err != nil {
		fmt.Printf("Error importing document: %v\n", err)
		os.Exit(1)
	}

	// Set case ID
	doc.CaseID = caseID

	// In a real implementation, save to repository
	app.repo.documents[doc.ID] = doc

	fmt.Printf("Document imported successfully. ID: %s, Type: %s\n",
		doc.ID, document.GetDocumentTypeString(doc.Type))
	fmt.Printf("Content preview: %s\n", preview(doc.Content, 150))
}

func (app *InvestigatorApp) handleEvidenceAdd(description, evidenceType, caseID string) {
	if description == "" {
		fmt.Println("Error: Evidence description is required")
		os.Exit(1)
	}

	if caseID == "" {
		if app.currentCaseID == "" {
			fmt.Println("Error: No case specified and no case is currently open")
			os.Exit(1)
		}
		caseID = app.currentCaseID
	}

	// Ensure the case exists
	_, err := app.caseService.GetCase(caseID)
	if err != nil {
		fmt.Printf("Error: Case not found: %v\n", err)
		os.Exit(1)
	}

	// Create evidence
	e := &evidence.Evidence{
		Description:    description,
		CaseID:         caseID,
		Type:           evidence.EvidenceType(strings.ToUpper(evidenceType)),
		Status:         evidence.StatusCollected,
		CollectedBy:    "Current User", // Would come from auth system
		CollectionDate: time.Now(),
		Location: evidence.Location{
			Description: "Not specified",
		},
		StorageLocation: "Evidence Locker",
	}

	err = app.evidenceService.CreateEvidence(e)
	if err != nil {
		fmt.Printf("Error adding evidence: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Evidence added successfully. ID: %s\n", e.ID)
}

func (app *InvestigatorApp) handleEvidenceList(caseID string) {
	if caseID == "" {
		if app.currentCaseID == "" {
			fmt.Println("Error: No case specified and no case is currently open")
			os.Exit(1)
		}
		caseID = app.currentCaseID
	}

	// Filter evidence by case ID
	var items []*evidence.Evidence
	for _, e := range app.repo.evidence {
		if e.CaseID == caseID {
			items = append(items, e)
		}
	}

	if len(items) == 0 {
		fmt.Printf("No evidence found for case: %s\n", caseID)
		return
	}

	fmt.Printf("\nEvidence for Case %s:\n", caseID)
	fmt.Println("-------------------------------------------------")
	fmt.Println("ID\t\tType\t\tStatus\tDescription")
	fmt.Println("-------------------------------------------------")

	for _, e := range items {
		fmt.Printf("%s\t%s\t%s\t%s\n",
			e.ID,
			e.Type,
			e.Status,
			e.Description)
	}
}

func (app *InvestigatorApp) handleInterviewAdd(title, interviewType, caseID string) {
	if title == "" {
		fmt.Println("Error: Interview title is required")
		os.Exit(1)
	}

	if caseID == "" {
		if app.currentCaseID == "" {
			fmt.Println("Error: No case specified and no case is currently open")
			os.Exit(1)
		}
		caseID = app.currentCaseID
	}

	// Ensure the case exists
	_, err := app.caseService.GetCase(caseID)
	if err != nil {
		fmt.Printf("Error: Case not found: %v\n", err)
		os.Exit(1)
	}

	// Create interview
	i := &interview.Interview{
		Title:         title,
		CaseID:        caseID,
		InterviewType: interview.InterviewType(strings.ToUpper(interviewType)),
		InterviewerID: "Current User", // Would come from auth system
		Date:          time.Now(),
		MediaType:     "Audio",
		Status:        "Scheduled",
	}

	err = app.interviewService.CreateInterview(i)
	if err != nil {
		fmt.Printf("Error adding interview: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Interview added successfully. ID: %s\n", i.ID)
}

func (app *InvestigatorApp) handleInterviewTranscribe(interviewID string) {
	if interviewID == "" {
		fmt.Println("Error: Interview ID is required")
		os.Exit(1)
	}

	// Get the interview
	i, err := app.interviewService.GetInterview(interviewID)
	if err != nil {
		fmt.Printf("Error: Interview not found: %v\n", err)
		os.Exit(1)
	}

	// Check if there's a recording path
	if i.RecordingPath == "" {
		fmt.Println("This interview has no recording to transcribe")
		return
	}

	// Create transcription options
	options := interview.SpeechRecognitionOptions{
		Language:     "en-US",
		MultiSpeaker: true,
	}

	// Transcribe
	transcript, err := app.interviewService.TranscribeInterview(interviewID, options)
	if err != nil {
		fmt.Printf("Error transcribing interview: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Interview transcribed successfully. Transcript ID: %s\n", transcript.ID)
	fmt.Printf("Transcript preview: %s\n", preview(transcript.Content, 150))
}

// New correspondence handlers
func (app *InvestigatorApp) handleCorrespondenceCreate(corrType, subject, body, recipient, caseID, templateID string) {
	if caseID == "" {
		if app.currentCaseID == "" {
			fmt.Println("Error: No case specified and no case is currently open")
			os.Exit(1)
		}
		caseID = app.currentCaseID
	}

	// Ensure the case exists
	_, err := app.caseService.GetCase(caseID)
	if err != nil {
		fmt.Printf("Error: Case not found: %v\n", err)
		os.Exit(1)
	}

	// Create simple sender (current user)
	sender := correspondence.Person{
		Name:        "Current User",
		IsOfficer:   true,
		Department:  "Police Department",
		BadgeNumber: "12345",
	}

	// Create recipient
	recipientPerson := correspondence.Person{
		Name: recipient,
	}

	var c *correspondence.Correspondence

	// If using a template
	if templateID != "" {
		// Create from template
		c, err = app.correspondenceService.CreateFromTemplate(
			templateID,
			caseID,
			sender,
			[]correspondence.Person{recipientPerson},
		)
		if err != nil {
			fmt.Printf("Error creating correspondence from template: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Create manually
		if corrType == "" {
			fmt.Println("Error: Correspondence type is required when not using a template")
			os.Exit(1)
		}

		if subject == "" {
			fmt.Println("Error: Subject is required when not using a template")
			os.Exit(1)
		}

		// Create correspondence
		c = &correspondence.Correspondence{
			CaseID:             caseID,
			CorrespondenceType: correspondence.CorrespondenceType(strings.ToUpper(corrType)),
			Subject:            subject,
			Body:               body,
			Sender:             sender,
			Recipients:         []correspondence.Person{recipientPerson},
			Direction:          "OUTGOING",
			Priority:           correspondence.PriorityNormal,
			Status:             correspondence.StatusDraft,
		}

		err = app.correspondenceService.CreateCorrespondence(c)
		if err != nil {
			fmt.Printf("Error creating correspondence: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Correspondence created successfully. ID: %s\n", c.ID)
	fmt.Printf("Status: %s, Type: %s\n", c.Status, c.CorrespondenceType)
	fmt.Printf("Subject: %s\n", c.Subject)
}

func (app *InvestigatorApp) handleCorrespondenceList(caseID string) {
	if caseID == "" {
		if app.currentCaseID == "" {
			fmt.Println("Error: No case specified and no case is currently open")
			os.Exit(1)
		}
		caseID = app.currentCaseID
	}

	// Filter correspondence by case ID
	var items []*correspondence.Correspondence
	for _, c := range app.repo.correspondence {
		if c.CaseID == caseID {
			items = append(items, c)
		}
	}

	if len(items) == 0 {
		fmt.Printf("No correspondence found for case: %s\n", caseID)
		return
	}

	fmt.Printf("\nCorrespondence for Case %s:\n", caseID)
	fmt.Println("-------------------------------------------------")
	fmt.Println("ID\t\tType\t\tStatus\tSubject")
	fmt.Println("-------------------------------------------------")

	for _, c := range items {
		fmt.Printf("%s\t%s\t%s\t%s\n",
			c.ID,
			c.CorrespondenceType,
			c.Status,
			c.Subject)
	}
}

func (app *InvestigatorApp) handleCorrespondenceSend(id string) {
	if id == "" {
		fmt.Println("Error: Correspondence ID is required")
		os.Exit(1)
	}

	// Send the correspondence
	err := app.correspondenceService.SendCorrespondence(id, time.Now())
	if err != nil {
		fmt.Printf("Error sending correspondence: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Correspondence sent successfully. ID: %s\n", id)
}

func (app *InvestigatorApp) handleCorrespondenceTemplateList() {
	// Get all templates
	templates := make([]*correspondence.Template, 0, len(app.repo.templates))
	for _, t := range app.repo.templates {
		templates = append(templates, t)
	}

	if len(templates) == 0 {
		fmt.Println("No templates found")
		return
	}

	fmt.Println("\nAvailable Correspondence Templates:")
	fmt.Println("-------------------------------------------------")
	fmt.Println("ID\t\tType\t\tName")
	fmt.Println("-------------------------------------------------")

	for _, t := range templates {
		fmt.Printf("%s\t%s\t%s\n",
			t.ID,
			t.Type,
			t.Name)
	}
}

// Helper functions
func preview(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// In-memory repositories implementation
type inMemoryCaseRepo struct {
	cases map[string]*casemanagement.Case
}

func (r *inMemoryCaseRepo) Save(c *casemanagement.Case) error {
	r.cases[c.ID] = c
	return nil
}

func (r *inMemoryCaseRepo) Find(id string) (*casemanagement.Case, error) {
	if c, ok := r.cases[id]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("case not found: %s", id)
}

func (r *inMemoryCaseRepo) FindByCaseNumber(caseNumber string) (*casemanagement.Case, error) {
	for _, c := range r.cases {
		if c.CaseNumber == caseNumber {
			return c, nil
		}
	}
	return nil, fmt.Errorf("case not found with number: %s", caseNumber)
}

func (r *inMemoryCaseRepo) Search(query string) ([]*casemanagement.Case, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *inMemoryCaseRepo) List(limit, offset int) ([]*casemanagement.Case, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *inMemoryCaseRepo) Update(c *casemanagement.Case) error {
	r.cases[c.ID] = c
	return nil
}

func (r *inMemoryCaseRepo) Delete(id string) error {
	delete(r.cases, id)
	return nil
}

type inMemoryCasefileRepo struct {
	cases map[string]*casefile.Case
}

func (r *inMemoryCasefileRepo) Save(c *casefile.Case) error {
	r.cases[c.ID] = c
	return nil
}

func (r *inMemoryCasefileRepo) Find(id string) (*casefile.Case, error) {
	if c, ok := r.cases[id]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("case not found: %s", id)
}

func (r *inMemoryCasefileRepo) FindByCaseNumber(caseNumber string) (*casefile.Case, error) {
	for _, c := range r.cases {
		if c.CaseNumber == caseNumber {
			return c, nil
		}
	}
	return nil, fmt.Errorf("case not found with number: %s", caseNumber)
}

func (r *inMemoryCasefileRepo) Search(query string) ([]*casefile.Case, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *inMemoryCasefileRepo) List(limit, offset int) ([]*casefile.Case, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *inMemoryCasefileRepo) Update(c *casefile.Case) error {
	r.cases[c.ID] = c
	return nil
}

func (r *inMemoryCasefileRepo) Delete(id string) error {
	delete(r.cases, id)
	return nil
}

type inMemoryEvidenceRepo struct {
	evidence map[string]*evidence.Evidence
}

func (r *inMemoryEvidenceRepo) Save(e *evidence.Evidence) error {
	r.evidence[e.ID] = e
	return nil
}

func (r *inMemoryEvidenceRepo) Find(id string) (*evidence.Evidence, error) {
	if e, ok := r.evidence[id]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("evidence not found: %s", id)
}

func (r *inMemoryEvidenceRepo) FindByCase(caseID string) ([]*evidence.Evidence, error) {
	var result []*evidence.Evidence
	for _, e := range r.evidence {
		if e.CaseID == caseID {
			result = append(result, e)
		}
	}
	return result, nil
}

func (r *inMemoryEvidenceRepo) Search(query string) ([]*evidence.Evidence, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *inMemoryEvidenceRepo) Update(e *evidence.Evidence) error {
	r.evidence[e.ID] = e
	return nil
}

func (r *inMemoryEvidenceRepo) Delete(id string) error {
	delete(r.evidence, id)
	return nil
}

type inMemoryInterviewRepo struct {
	interviews map[string]*interview.Interview
}

func (r *inMemoryInterviewRepo) Save(i *interview.Interview) error {
	r.interviews[i.ID] = i
	return nil
}

func (r *inMemoryInterviewRepo) Find(id string) (*interview.Interview, error) {
	if i, ok := r.interviews[id]; ok {
		return i, nil
	}
	return nil, fmt.Errorf("interview not found: %s", id)
}

func (r *inMemoryInterviewRepo) FindByCase(caseID string) ([]*interview.Interview, error) {
	var result []*interview.Interview
	for _, i := range r.interviews {
		if i.CaseID == caseID {
			result = append(result, i)
		}
	}
	return result, nil
}

func (r *inMemoryInterviewRepo) Search(query string) ([]*interview.Interview, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *inMemoryInterviewRepo) Update(i *interview.Interview) error {
	r.interviews[i.ID] = i
	return nil
}

func (r *inMemoryInterviewRepo) Delete(id string) error {
	delete(r.interviews, id)
	return nil
}

type inMemoryTranscriptRepo struct {
	transcripts map[string]*interview.Transcript
}

func (r *inMemoryTranscriptRepo) Save(t *interview.Transcript) error {
	r.transcripts[t.ID] = t
	return nil
}

func (r *inMemoryTranscriptRepo) Find(id string) (*interview.Transcript, error) {
	if t, ok := r.transcripts[id]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("transcript not found: %s", id)
}

func (r *inMemoryTranscriptRepo) FindByInterview(interviewID string) (*interview.Transcript, error) {
	for _, t := range r.transcripts {
		if t.InterviewID == interviewID {
			return t, nil
		}
	}
	return nil, fmt.Errorf("transcript not found for interview: %s", interviewID)
}

func (r *inMemoryTranscriptRepo) Update(t *interview.Transcript) error {
	r.transcripts[t.ID] = t
	return nil
}

func (r *inMemoryTranscriptRepo) Delete(id string) error {
	delete(r.transcripts, id)
	return nil
}

// Dummy speech recognizer for demonstration
type dummySpeechRecognizer struct{}

func (r *dummySpeechRecognizer) Initialize() error {
	return nil
}

func (r *dummySpeechRecognizer) Transcribe(audioPath string, options interview.SpeechRecognitionOptions) (*interview.Transcript, error) {
	// Create a dummy transcript
	return &interview.Transcript{
		Content:     "This is a simulated transcript. In a real implementation, this would contain the actual transcribed text from the audio file.",
		Language:    options.Language,
		IsAutomated: true,
		Segments: []interview.Segment{
			{
				SpeakerRole: "Interviewer",
				StartTime:   0,
				EndTime:     30 * time.Second,
				Text:        "Can you describe what you witnessed on the night of the incident?",
				Confidence:  0.95,
			},
			{
				SpeakerRole: "Interviewee",
				StartTime:   31 * time.Second,
				EndTime:     90 * time.Second,
				Text:        "I was walking my dog when I heard a loud noise coming from the alley. When I looked, I saw someone running away.",
				Confidence:  0.87,
			},
		},
	}, nil
}

func (r *dummySpeechRecognizer) Close() error {
	return nil
}

// Add new repository implementations
type inMemoryCorrespondenceRepo struct {
	correspondence map[string]*correspondence.Correspondence
}

func (r *inMemoryCorrespondenceRepo) Save(c *correspondence.Correspondence) error {
	r.correspondence[c.ID] = c
	return nil
}

func (r *inMemoryCorrespondenceRepo) Find(id string) (*correspondence.Correspondence, error) {
	if c, ok := r.correspondence[id]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("correspondence not found: %s", id)
}

func (r *inMemoryCorrespondenceRepo) FindByCase(caseID string) ([]*correspondence.Correspondence, error) {
	var result []*correspondence.Correspondence
	for _, c := range r.correspondence {
		if c.CaseID == caseID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *inMemoryCorrespondenceRepo) FindByType(correspondenceType correspondence.CorrespondenceType) ([]*correspondence.Correspondence, error) {
	var result []*correspondence.Correspondence
	for _, c := range r.correspondence {
		if c.CorrespondenceType == correspondenceType {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *inMemoryCorrespondenceRepo) FindByStatus(status correspondence.Status) ([]*correspondence.Correspondence, error) {
	var result []*correspondence.Correspondence
	for _, c := range r.correspondence {
		if c.Status == status {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *inMemoryCorrespondenceRepo) FindByReference(refNumber string) (*correspondence.Correspondence, error) {
	for _, c := range r.correspondence {
		if c.ReferenceNumber == refNumber {
			return c, nil
		}
	}
	return nil, fmt.Errorf("correspondence not found with reference number: %s", refNumber)
}

func (r *inMemoryCorrespondenceRepo) Search(query string) ([]*correspondence.Correspondence, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *inMemoryCorrespondenceRepo) Update(c *correspondence.Correspondence) error {
	r.correspondence[c.ID] = c
	return nil
}

func (r *inMemoryCorrespondenceRepo) Delete(id string) error {
	delete(r.correspondence, id)
	return nil
}

type inMemoryTemplateRepo struct {
	templates map[string]*correspondence.Template
}

func (r *inMemoryTemplateRepo) Save(t *correspondence.Template) error {
	r.templates[t.ID] = t
	return nil
}

func (r *inMemoryTemplateRepo) Find(id string) (*correspondence.Template, error) {
	if t, ok := r.templates[id]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("template not found: %s", id)
}

func (r *inMemoryTemplateRepo) FindByName(name string) (*correspondence.Template, error) {
	for _, t := range r.templates {
		if t.Name == name {
			return t, nil
		}
	}
	return nil, fmt.Errorf("template not found with name: %s", name)
}

func (r *inMemoryTemplateRepo) FindByType(correspondenceType correspondence.CorrespondenceType) ([]*correspondence.Template, error) {
	var result []*correspondence.Template
	for _, t := range r.templates {
		if t.Type == correspondenceType {
			result = append(result, t)
		}
	}
	return result, nil
}

func (r *inMemoryTemplateRepo) FindByDepartment(department string) ([]*correspondence.Template, error) {
	var result []*correspondence.Template
	for _, t := range r.templates {
		if t.Department == department || t.Department == "Any" {
			result = append(result, t)
		}
	}
	return result, nil
}

func (r *inMemoryTemplateRepo) Update(t *correspondence.Template) error {
	r.templates[t.ID] = t
	return nil
}

func (r *inMemoryTemplateRepo) Delete(id string) error {
	delete(r.templates, id)
	return nil
}
