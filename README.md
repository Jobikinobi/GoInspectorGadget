# GoInspectorGadget

A comprehensive investigation management system for law enforcement agencies.

## Overview

GoInspectorGadget provides an integrated platform for managing investigations, evidence, documents, interviews, and correspondence. Built with Go, it offers a command-line interface for efficient case management and processing.

## Features

- **Case Management**: Create, view, and manage investigation cases
- **Document Processing**: Import, analyze, and extract content from documents
- **Evidence Tracking**: Record and track physical and digital evidence
- **Interview Management**: Manage and transcribe interview recordings
- **Correspondence Handling**: Create and send official correspondence using templates
- **Audio Processing**: Transcribe audio with support for multiple languages and accents

## Documentation

- [Installation Guide](docs/INSTALLATION.md) - Detailed instructions for installing and configuring the system
- [User Manual](docs/USER_MANUAL.md) - Comprehensive guide to using all features
- [Quick Reference](docs/QUICK_REFERENCE.md) - Command reference for common tasks
- [Developer Guide](docs/DEVELOPER_GUIDE.md) - Guide for extending and contributing to the system

## Installation

### Prerequisites

- Go 1.23.0 or higher
- FFmpeg (for audio processing)

### Install from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/jth/claude/GoInspectorGadget.git
   cd GoInspectorGadget
   ```

2. Build the binaries:
   ```bash
   go build -o bin/investigator cmd/investigator/main.go
   go build -o bin/docprocessor cmd/docprocessor/*.go
   ```

3. Add the binaries to your PATH or move them to a location in your PATH:
   ```bash
   export PATH=$PATH:$(pwd)/bin
   ```

For detailed installation instructions, see the [Installation Guide](docs/INSTALLATION.md).

## Usage

### Case Management

Create a new case:
```bash
investigator case create --title "Missing Person - John Doe" --desc "Investigation into disappearance of John Doe" --type "Missing Person"
```

List all cases:
```bash
investigator case list
```

Open a specific case:
```bash
investigator case open CASE-1234567890
```

### Document Management

Import a document to a case:
```bash
investigator doc import --path "/path/to/document.pdf" --case CASE-1234567890
```

### Evidence Management

Add evidence to a case:
```bash
investigator evidence add --desc "Wallet found at scene" --type "PHYSICAL" --case CASE-1234567890
```

List evidence for a case:
```bash
investigator evidence list CASE-1234567890
```

### Interview Management

Add an interview:
```bash
investigator interview add --title "Witness Interview - Jane Smith" --type "WITNESS" --case CASE-1234567890
```

Transcribe an interview recording:
```bash
investigator interview transcribe --id INT-1234567890
```

### Correspondence

Create correspondence using a template:
```bash
investigator correspondence create --template TEMPLATE-ID --recipient "John Smith" --case CASE-1234567890
```

List available correspondence templates:
```bash
investigator correspondence templates
```

Send correspondence:
```bash
investigator correspondence send --id CORR-1234567890
```

### Audio Processing

Process an audio file for transcription:
```bash
docprocessor --type audio --audio "/path/to/recording.wav" --accent "auto"
```

Process an interview recording with speaker diarization:
```bash
docprocessor --type interview --input "/path/to/interview.wav" --output "transcript.txt"
```

For more detailed usage instructions, see the [User Manual](docs/USER_MANUAL.md).

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
- `docs/`: Documentation
  - `INSTALLATION.md`: Detailed installation instructions
  - `USER_MANUAL.md`: Comprehensive user guide
  - `QUICK_REFERENCE.md`: Command reference

## License

This project is licensed under the terms specified in the LICENSE file.
