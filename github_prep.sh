#!/bin/bash

# GitHub preparation script for the GoCode project

echo "Preparing the GoCode project for GitHub..."

# Update .gitignore
echo "Creating comprehensive .gitignore..."
cat > .gitignore << 'EOF'
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# Go workspace file
go.work

# IDE files
.idea/
.vscode/
*.swp
*.swo

# OS specific files
.DS_Store
Thumbs.db

# Binary output directory
bin/

# Build artifacts
docprocessor
investigator
**/docprocessor/docprocessor
**/investigator/investigator

# Backup files
*.bak
EOF

# Create a README file if it doesn't exist or update it
echo "Updating README.md..."
cat > README.md << 'EOF'
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
│   ├── casemanagement/        # Case management functionality
│   ├── casefile/              # Case file handling
│   ├── correspondence/        # Communication management
│   ├── document/              # Document processing
│   ├── evidence/              # Evidence tracking
│   ├── interview/             # Interview management
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
go build ./cmd/investigator

# Build the document processor
go build ./cmd/docprocessor
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

## License

This project is licensed under the MIT License - see the LICENSE file for details.
EOF

# Create a LICENSE file if it doesn't exist
if [ ! -f LICENSE ]; then
  echo "Creating MIT LICENSE file..."
  cat > LICENSE << 'EOF'
MIT License

Copyright (c) 2023 JTH

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
EOF
fi

# Final verification build
echo "Running final verification build..."
go build ./cmd/docprocessor
go build ./cmd/investigator

echo "GitHub preparation complete!"
echo "You can now push your repository to GitHub:"
echo "1. git add ."
echo "2. git commit -m \"Prepare for GitHub release\""
echo "3. git remote add origin https://github.com/yourusername/repo.git"
echo "4. git push -u origin main" 