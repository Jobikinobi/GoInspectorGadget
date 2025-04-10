# GoInspectorGadget Developer Guide

This guide is intended for developers who want to extend, modify, or contribute to the GoInspectorGadget system.

## Architecture Overview

GoInspectorGadget follows a modular architecture with well-defined interfaces:

```
┌───────────────────┐      ┌────────────────────┐
│ Command Interface │      │ Document Processor  │
│ (investigator)    │      │ (docprocessor)     │
└─────────┬─────────┘      └──────────┬─────────┘
          │                           │
          ▼                           ▼
┌───────────────────────────────────────────────┐
│                Core Packages                  │
│ ┌─────────┐ ┌──────────┐ ┌────────────────┐  │
│ │casefile │ │interview │ │correspondence  │  │
│ └─────────┘ └──────────┘ └────────────────┘  │
│ ┌─────────┐ ┌──────────┐ ┌────────────────┐  │
│ │evidence │ │document  │ │speech          │  │
│ └─────────┘ └──────────┘ └────────────────┘  │
└───────────────────────────────────────────────┘
```

### Key Design Principles

1. **Interface-Based Design**: Core functionality is defined using interfaces, allowing for different implementations.
2. **Repository Pattern**: Data access is abstracted through repository interfaces.
3. **Service Layer**: Business logic is encapsulated in services that depend on repositories.
4. **Command-Line Interface**: User interaction is through command-line tools.

## Development Environment Setup

### Prerequisites

- Go 1.23.0 or higher
- Git
- Your favorite IDE or text editor (VS Code recommended with Go extension)
- FFmpeg (for audio processing development)

### Setting Up

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR-USERNAME/GoInspectorGadget.git
   cd GoInspectorGadget
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Build the project to verify your setup:
   ```bash
   go build ./...
   ```

## Project Structure

- `cmd/`: Application entry points
  - `investigator/`: Main investigation management tool
  - `docprocessor/`: Document and audio processing tool
  
- `pkg/`: Core packages and functionality
  - `casefile/`: Case file management
  - `casemanagement/`: Case tracking and workflow
  - `document/`: Document processing and analysis
  - `evidence/`: Evidence tracking and chain of custody
  - `interview/`: Interview management and transcription
  - `correspondence/`: Communication templates and tracking
  - `speech/`: Speech recognition and transcription

## Key Interfaces

Understanding these interfaces is essential for extending the system:

### Repository Interfaces

Each domain has repository interfaces for data access:

```go
// Example: Interview Repository
type InterviewRepository interface {
    Save(i *Interview) error
    Find(id string) (*Interview, error)
    FindByCase(caseID string) ([]*Interview, error)
    Search(query string) ([]*Interview, error)
    Update(i *Interview) error
    Delete(id string) error
}
```

### Service Interfaces

Services implement business logic and depend on repositories:

```go
// Example: Speech Recognizer
type SpeechRecognizer interface {
    Initialize() error
    Transcribe(audioPath string, options SpeechRecognitionOptions) (*Transcript, error)
    Close() error
}
```

## Extending the System

### Adding a New Repository Implementation

To add a new storage backend (e.g., SQL database):

1. Create a new file for your implementation:
   ```go
   // pkg/interview/sql_repository.go
   package interview
   
   import "database/sql"
   
   type sqlInterviewRepository struct {
       db *sql.DB
   }
   
   func NewSQLInterviewRepository(db *sql.DB) InterviewRepository {
       return &sqlInterviewRepository{db: db}
   }
   
   func (r *sqlInterviewRepository) Save(i *Interview) error {
       // Implementation here
   }
   
   // Implement other methods...
   ```

2. Update the application to use your new repository:
   ```go
   // In cmd/investigator/main.go
   db, err := sql.Open("mysql", "user:password@/dbname")
   if err != nil {
       // Handle error
   }
   interviewRepo := interview.NewSQLInterviewRepository(db)
   ```

### Adding a New Speech Recognition Engine

To integrate a different speech recognition technology:

1. Create a new file for your implementation:
   ```go
   // pkg/speech/deepgramrecognizer.go
   package speech
   
   type DeepgramRecognizer struct {
       apiKey string
       // Other fields...
   }
   
   func NewDeepgramRecognizer(apiKey string) SpeechRecognizer {
       return &DeepgramRecognizer{apiKey: apiKey}
   }
   
   func (r *DeepgramRecognizer) Initialize() error {
       // Implementation here
   }
   
   // Implement other methods...
   ```

2. Update the factory function to include your new recognizer:
   ```go
   // In pkg/speech/recognizer.go
   func RecognizerFactory(accentType AccentType, options RecognizerOptions) (SpeechRecognizer, error) {
       if options.Provider == "deepgram" {
           return NewDeepgramRecognizer(options.APIKey), nil
       }
       // Existing code...
   }
   ```

### Adding a New Command

To add a new command to the investigator tool:

1. Define a new flag set in main.go:
   ```go
   // In cmd/investigator/main.go
   reportGenCmd := flag.NewFlagSet("report generate", flag.ExitOnError)
   reportType := reportGenCmd.String("type", "summary", "Report type (summary, detailed)")
   reportCase := reportGenCmd.String("case", "", "Case ID to generate report for")
   ```

2. Add a handler function:
   ```go
   func (app *InvestigatorApp) handleReportGenerate(reportType, caseID string) {
       // Implementation here
   }
   ```

3. Add the command to the main command switch:
   ```go
   case "report":
       if len(os.Args) < 3 {
           fmt.Println("Missing report subcommand")
           os.Exit(1)
       }

       switch os.Args[2] {
       case "generate":
           reportGenCmd.Parse(os.Args[3:])
           app.handleReportGenerate(*reportType, *reportCase)
       default:
           fmt.Printf("Unknown report subcommand: %s\n", os.Args[2])
           os.Exit(1)
       }
   ```

## Testing

### Unit Testing

Unit tests should be written for all packages. Tests should be in the same package as the code they test, with a `_test.go` suffix:

```go
// pkg/interview/interview_test.go
package interview

import "testing"

func TestCreateInterview(t *testing.T) {
    // Test implementation
}
```

Run tests with:
```bash
go test ./...
```

### Integration Testing

Integration tests should verify interactions between components:

```go
// pkg/interview/integration_test.go
package interview_test  // Note: different package name

import (
    "testing"
    "github.com/jth/claude/GoInspectorGadget/pkg/interview"
)

func TestInterviewTranscription(t *testing.T) {
    // Test implementation using real components
}
```

## Code Style and Conventions

Follow standard Go conventions:

1. Use `gofmt` or `goimports` to format code:
   ```bash
   gofmt -w .
   ```

2. Follow Go naming conventions:
   - Use camelCase for variable names
   - Use PascalCase for exported names
   - Use short but descriptive names

3. Use descriptive error messages:
   ```go
   return nil, fmt.Errorf("failed to find interview: %w", err)
   ```

4. Document all exported functions and types using godoc format:
   ```go
   // InterviewService provides business logic for interview management.
   // It handles creation, retrieval, and transcription of interviews.
   type InterviewService struct {
       // ...
   }
   ```

## Documentation Standards

### Markdown Linting

All documentation files follow standard Markdown conventions to ensure compatibility with various Markdown parsers and tools.

1. Use the provided markdown linting script before submitting changes:
   ```bash
   ./scripts/lint_markdown.sh
   ```

2. Key Markdown standards:
   - Use fenced code blocks with language identifiers (e.g., ```bash, ```go)
   - Keep line length under 120 characters
   - Use backticks (`) for inline code
   - Use consistent heading styles (# for top level, ## for section, etc.)
   - Use ordered lists with incrementing numbers

3. Install markdown linting tools:
   ```bash
   npm install -g markdownlint-cli
   ```

These standards ensure that our markdown files render correctly in GitHub, documentation sites, and tools like Peek.

## Pull Request Process

1. Create a feature branch for your changes:
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. Make your changes and commit them with descriptive messages:
   ```bash
   git commit -m "Add support for PDF/A format in document processor"
   ```

3. Run tests and linters:
   ```bash
   go test ./...
   go vet ./...
   ./scripts/lint_markdown.sh  # If documentation was changed
   ```

4. Push your changes and create a pull request:
   ```bash
   git push origin feature/my-new-feature
   ```

5. Wait for code review and address any feedback.

## Common Patterns and Best Practices

### Error Handling

Use wrapped errors for better context:

```go
if err != nil {
    return fmt.Errorf("failed to initialize recognizer: %w", err)
}
```

### Dependency Injection

Services should receive their dependencies through constructors:

```go
func NewInterviewService(
    interviewRepo InterviewRepository,
    transcriptRepo TranscriptRepository,
    recognizer SpeechRecognizer,
) *InterviewService {
    return &InterviewService{
        interviewRepo:  interviewRepo,
        transcriptRepo: transcriptRepo,
        recognizer:     recognizer,
    }
}
```

### Interface Segregation

Keep interfaces focused on specific needs:

```go
// Good
type DocumentReader interface {
    Read(id string) (*Document, error)
}

type DocumentWriter interface {
    Write(doc *Document) error
}

// Use composition when needed
type DocumentRepository interface {
    DocumentReader
    DocumentWriter
}
``` 