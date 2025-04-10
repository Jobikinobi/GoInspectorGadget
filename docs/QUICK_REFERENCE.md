# GoInspectorGadget Quick Reference Guide

## Case Management

| Task | Command |
|------|---------|
| Create a case | `investigator case create --title "Title" --desc "Description" --type "Type"` |
| List all cases | `investigator case list` |
| Open a case | `investigator case open CASE-ID` |

## Document Management

| Task | Command |
|------|---------|
| Import document | `investigator doc import --path "/path/to/doc.pdf" --case CASE-ID` |

## Evidence Management

| Task | Command |
|------|---------|
| Add evidence | `investigator evidence add --desc "Description" --type "TYPE" --case CASE-ID` |
| List evidence | `investigator evidence list CASE-ID` |

## Interview Management

| Task | Command |
|------|---------|
| Add interview | `investigator interview add --title "Title" --type "TYPE" --case CASE-ID` |
| Transcribe interview | `investigator interview transcribe --id INT-ID` |

## Correspondence

| Task | Command |
|------|---------|
| Create from template | `investigator correspondence create --template TEMPLATE-ID --recipient "Name" --case CASE-ID` |
| Create custom | `investigator correspondence create --type "TYPE" --subject "Subject" --body "Content" --recipient "Name" --case CASE-ID` |
| List templates | `investigator correspondence templates` |
| List correspondence | `investigator correspondence list CASE-ID` |
| Send correspondence | `investigator correspondence send --id CORR-ID` |

## Audio Processing

| Task | Command |
|------|---------|
| Process audio | `docprocessor --type audio --audio "/path/to/file.wav" --accent "TYPE"` |
| Process interview | `docprocessor --type interview --input "/path/to/file.wav" --output "transcript.txt"` |

## Common Options

### Case Types
- Homicide
- Theft
- Assault
- Missing Person
- Fraud
- Other

### Evidence Types
- PHYSICAL
- DIGITAL
- DOCUMENTARY
- TESTIMONIAL
- DEMONSTRATIVE

### Interview Types
- WITNESS
- SUSPECT
- VICTIM
- EXPERT
- INFORMER

### Correspondence Types
- EMAIL
- LETTER
- MEMO
- SUBPOENA
- WARRANT

### Accent Types
- venezuelan
- american
- generic
- auto (default) 