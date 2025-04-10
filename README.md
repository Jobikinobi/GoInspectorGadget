# Police Investigation Document and Interview Processing System

A Go-based toolkit for managing police investigations, processing evidence, interviews, and documents with specialized multi-accent speech recognition support.

## Project Overview

This system is designed to aid law enforcement in managing investigations by providing:

1. Case management and evidence tracking
2. Document processing and analysis
3. Multi-accent speech recognition for interview transcription
4. Correspondence management with templates

## Key Features

### Case Management
- Track case details, status, and assignments
- Manage evidence and document collection
- Link cases to suspects, victims, and witnesses

### Document Processing
- Import and analyze various document types
- Extract text from PDFs and images
- Automatic document classification

### Speech Recognition
- Support for multiple accents:
  - Venezuelan Spanish
  - American English (with police terminology)
  - Generic accent fallback
- Speaker identification and diarization
- Interview transcription and analysis

### Correspondence Management
- Template-based correspondence generation
- Track communication threads
- Manage official communications

## Project Structure

```
/GoCode
├── cmd/                       # Command-line applications
│   ├── docprocessor/          # Document processing utility
│   └── investigator/          # Main investigator application
├── pkg/                       # Library packages
│   ├── audio/                 # Audio processing utilities
│   ├── casemanagement/        # Case management functionality
│   ├── casefile/              # Case file handling
│   ├── correspondence/        # Communication management
│   ├── document/              # Document processing
│   ├── evidence/              # Evidence tracking
│   ├── hashicorp/             # Credential management
│   ├── image/                 # Image processing utilities
│   ├── indexing/              # Search and indexing
│   ├── interview/             # Interview management
│   ├── security/              # Security and access control
│   └── speech/                # Speech recognition
├── go.mod                     # Go module definition
└── go.sum                     # Go module checksums
```

## Getting Started

### Prerequisites
- Go 1.18 or higher
- FFmpeg for audio processing
- Poppler for PDF processing

### Installation

```bash
# Clone the repository
git clone https://github.com/jth/claude/GoCode.git
cd GoCode

# Build the investigator tool
go build -o investigator ./cmd/investigator

# Build the document processor
go build -o docprocessor ./cmd/docprocessor
```

### Basic Usage

```bash
# Create a new case
./investigator case create --title "Missing Person" --desc "Investigation into missing person report" --type "Missing Person"

# Import a document into a case
./investigator doc import --path "/path/to/document.pdf" --case <case-id>

# Process an interview recording
./docprocessor --audio "/path/to/interview.mp3" --accent "venezuelan"
```

## Import Path Structure

This project uses the module path `github.com/jth/claude/GoCode` for all imports. If forking or modifying the code, you may need to update the import paths to match your repository structure.

## Development Status

This project is currently in development. Some features may be incomplete or function as mock implementations.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 