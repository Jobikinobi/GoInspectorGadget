# GoInspectorGadget User Manual

Version 1.0.0

## Table of Contents

1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Getting Started](#getting-started)
4. [Case Management](#case-management)
5. [Document Processing](#document-processing)
6. [Evidence Management](#evidence-management)
7. [Interview Management](#interview-management)
8. [Correspondence](#correspondence)
9. [Audio Processing](#audio-processing)
10. [Command Reference](#command-reference)
11. [Best Practices](#best-practices)
12. [Troubleshooting](#troubleshooting)
13. [Technical Support](#technical-support)

## Introduction

GoInspectorGadget is a comprehensive investigation management system designed for law enforcement agencies to streamline the processes of collecting, managing, and analyzing evidence, conducting interviews, preparing documentation, and maintaining chains of custody.

### System Overview

The system consists of two main components:

1. **Investigator Tool**: A command-line interface for case management, evidence tracking, interview management, and correspondence.
2. **Document Processor**: A specialized tool for processing documents and audio files, including transcription with multi-accent support.

### Key Features

- Case management with customizable case types and statuses
- Document import and analysis
- Evidence tracking with chain of custody
- Interview recording and transcription
- Automated correspondence with templates
- Multi-accent speech recognition

## Installation

### System Requirements

- Operating System: Linux, macOS, or Windows
- Processor: 2.0 GHz dual-core processor or better
- Memory: 4 GB RAM minimum (8 GB recommended)
- Disk Space: 1 GB available space
- Go 1.23.0 or higher
- FFmpeg for audio processing

### Installation Steps

#### From Source

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

3. Add the binaries to your PATH:
   
   For Linux/macOS:
   ```bash
   export PATH=$PATH:$(pwd)/bin
   ```
   
   For Windows:
   ```cmd
   set PATH=%PATH%;%cd%\bin
   ```

#### Optional Components

For optimal speech recognition performance, you may need to install additional language models:

- Venezuelan Spanish model: Download from our model repository
- American English with police terminology: Download from our model repository

## Getting Started

### Initial Setup

After installation, you can verify the installation by running:

```bash
investigator help
```

This should display the help menu with available commands.

### Working Directory

By default, GoInspectorGadget stores its data in:

- Linux/macOS: `$HOME/investigator-simulator`
- Windows: `%USERPROFILE%\investigator-simulator`

You can customize this location using environment variables if needed.

### First Steps

1. Create your first case:
   ```bash
   investigator case create --title "Sample Case" --desc "Description" --type "Theft"
   ```

2. Note down the Case ID that is generated, as you'll need it for subsequent operations

## Case Management

Case management is the core functionality of GoInspectorGadget, providing a central point for organizing all information related to an investigation.

### Creating a Case

```bash
investigator case create --title "Case Title" --desc "Case Description" --type "Case Type"
```

Available case types include:
- Homicide
- Theft
- Assault
- Missing Person
- Fraud
- Other

### Listing Cases

To view all cases in the system:

```bash
investigator case list
```

### Opening a Case

To open a specific case and set it as the current working case:

```bash
investigator case open CASE-1234567890
```

When a case is open, you can perform operations without specifying the case ID each time.

### Case Status

Cases can have the following statuses:
- Open
- Pending
- Closed
- Suspended

## Document Processing

GoInspectorGadget can import and analyze various document types, including PDFs, images, and text files.

### Importing Documents

To import a document into a case:

```bash
investigator doc import --path "/path/to/document.pdf" --case CASE-1234567890
```

If a case is already open, you can omit the --case parameter:

```bash
investigator doc import --path "/path/to/document.pdf"
```

### Supported Document Types

- PDF documents
- Text files
- Image files (JPG, PNG, TIFF)
- Word documents (DOCX)

### Document Analysis

Documents are automatically analyzed upon import:
- Text extraction
- OCR for image-based documents
- Metadata analysis

## Evidence Management

GoInspectorGadget provides comprehensive tools for tracking physical and digital evidence.

### Adding Evidence

To add evidence to a case:

```bash
investigator evidence add --desc "Evidence Description" --type "PHYSICAL" --case CASE-1234567890
```

Evidence types include:
- PHYSICAL: Physical items
- DIGITAL: Digital files, data, etc.
- DOCUMENTARY: Paper documents
- TESTIMONIAL: Witness testimony
- DEMONSTRATIVE: Maps, charts, etc.

### Listing Evidence

To list all evidence for a case:

```bash
investigator evidence list CASE-1234567890
```

Or if a case is already open:

```bash
investigator evidence list
```

### Chain of Custody

Each piece of evidence automatically maintains a chain of custody that records:
- Who collected the evidence
- When it was collected
- Location of collection
- Current storage location
- Any transfers or handling

## Interview Management

GoInspectorGadget allows you to manage interview records and transcribe audio recordings.

### Adding an Interview

To add a new interview record:

```bash
investigator interview add --title "Witness Interview - John Smith" --type "WITNESS" --case CASE-1234567890
```

Interview types include:
- WITNESS: Interviews with witnesses
- SUSPECT: Interviews with suspects
- VICTIM: Interviews with victims
- EXPERT: Interviews with expert witnesses
- INFORMER: Interviews with confidential informants

### Transcribing Interviews

To transcribe an interview recording:

```bash
investigator interview transcribe --id INT-1234567890
```

### Interview Statuses

Interviews can have the following statuses:
- SCHEDULED: Interview is planned
- COMPLETED: Interview has been conducted
- CANCELLED: Interview did not take place
- POSTPONED: Interview has been rescheduled

## Correspondence

GoInspectorGadget includes a correspondence management system for generating and tracking official communications.

### Creating Correspondence

Using a template:

```bash
investigator correspondence create --template TEMPLATE-ID --recipient "John Smith" --case CASE-1234567890
```

Custom correspondence:

```bash
investigator correspondence create --type "EMAIL" --subject "Subject Line" --body "Content" --recipient "John Smith" --case CASE-1234567890
```

### Listing Templates

To view available correspondence templates:

```bash
investigator correspondence templates
```

### Listing Correspondence

To list all correspondence for a case:

```bash
investigator correspondence list CASE-1234567890
```

### Sending Correspondence

To mark correspondence as sent:

```bash
investigator correspondence send --id CORR-1234567890
```

### Correspondence Types

- EMAIL: Electronic mail
- LETTER: Physical letter
- MEMO: Internal memorandum
- SUBPOENA: Legal summons
- WARRANT: Legal warrant

## Audio Processing

GoInspectorGadget includes a powerful audio processing system that supports transcription with accent detection.

### Processing Audio Files

Basic audio processing:

```bash
docprocessor --type audio --audio "/path/to/recording.wav"
```

With accent specification:

```bash
docprocessor --type audio --audio "/path/to/recording.wav" --accent "venezuelan"
```

### Processing Interview Recordings

For interview recordings with multiple speakers:

```bash
docprocessor --type interview --input "/path/to/interview.wav" --output "transcript.txt"
```

This process includes:
1. Speaker diarization (identifying different speakers)
2. Accent detection for each speaker
3. Transcription using the appropriate language model
4. Combining into a timestamped transcript

### Supported Accents

- venezuelan: Venezuelan Spanish
- american: American English (with police terminology)
- generic: General accent (fallback)
- auto: Automatic detection (default)

### GPU Acceleration

For faster processing, GPU acceleration can be enabled:

```bash
docprocessor --type audio --audio "/path/to/recording.wav" --gpu
```

## Command Reference

### Investigator Commands

| Command | Description |
|---------|-------------|
| `investigator case create` | Create a new case |
| `investigator case open` | Open an existing case |
| `investigator case list` | List all cases |
| `investigator doc import` | Import a document |
| `investigator evidence add` | Add new evidence |
| `investigator evidence list` | List evidence for a case |
| `investigator interview add` | Add a new interview |
| `investigator interview transcribe` | Transcribe an interview recording |
| `investigator correspondence create` | Create new correspondence |
| `investigator correspondence list` | List correspondence for a case |
| `investigator correspondence send` | Mark correspondence as sent |
| `investigator correspondence templates` | List available templates |
| `investigator help` | Display help information |

### Document Processor Commands

| Command | Description |
|---------|-------------|
| `docprocessor --type audio` | Process an audio file |
| `docprocessor --type interview` | Process an interview recording |

## Best Practices

### Case Management

- Create a new case as soon as an investigation begins
- Use descriptive case titles that include key information
- Include the case type for easier categorization

### Evidence Handling

- Add evidence to the system immediately after collection
- Provide detailed descriptions for each piece of evidence
- Include exact locations where evidence was found

### Interview Management

- Prepare interview questions in advance
- Record all interviews when possible
- Transcribe interviews promptly

### Document Processing

- Organize documents by type before importing
- Verify document content after import
- Back up important documents separately

## Troubleshooting

### Common Issues

#### Case Creation Fails

**Issue**: Unable to create a new case.
**Solution**: Ensure you have specified required fields (title). Check permissions on the data directory.

#### Audio Processing Fails

**Issue**: Audio transcription fails or produces poor results.
**Solution**: 
- Check that the audio file is in a supported format
- Ensure the file isn't corrupted
- Try specifying the accent manually instead of auto-detection
- For large files, enable GPU acceleration if available

#### Import Path Errors

**Issue**: Errors about missing packages or incorrect import paths.
**Solution**: 
- Verify that the module path in go.mod is correct
- Run `go mod tidy` to clean up dependencies
- Rebuild the application

### Error Messages

| Error Message | Meaning | Solution |
|---------------|---------|----------|
| "Error: Case title is required" | The case title was not provided | Specify a title with the --title flag |
| "Error: Input file is required" | No audio file specified | Provide the path to the audio file |
| "Error: Audio file not found" | The specified audio file doesn't exist | Check the file path and try again |
| "Error: No case specified and no case is currently open" | No case ID provided and none is open | Either open a case first or specify the case ID |

## Technical Support

For technical assistance, please contact our support team:

- Email: support@goinspectorgadget.example.com
- Phone: 1-800-INSPECT
- Support Hours: Monday-Friday, 9am-5pm Eastern Time

When reporting issues, please include:
- The exact command you were running
- Any error messages displayed
- The version of GoInspectorGadget you're using
- Your operating system and version 