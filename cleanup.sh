#!/bin/bash

# Cleanup script for the GoCode project

echo "Cleaning up the GoCode project..."

# Remove any build artifacts and potential duplicate files
rm -f ./cmd/docprocessor/docprocessor
rm -f ./cmd/investigator/investigator
rm -f ./docprocessor
rm -f ./investigator

# Move pdf.go to a backup location if it exists
if [ -f ./pkg/document/pdf.go ]; then
  echo "Backing up pdf.go..."
  mv ./pkg/document/pdf.go ./pkg/document/pdf.go.bak
fi

# Create the document directory structure if it doesn't exist
mkdir -p ./pkg/document
mkdir -p ./pkg/casemanagement
mkdir -p ./cmd/docprocessor
mkdir -p ./cmd/investigator

# Fix import statements if needed
echo "Fixing import statements..."
find . -name "*.go" -exec sed -i '' 's|github.com/jth/docprocessor|github.com/jth/claude/GoCode|g' {} \;

# Rebuild the project
echo "Building the project..."
go build ./cmd/docprocessor
go build ./cmd/investigator

echo "Cleanup complete!"

# Show the project structure
echo "Current project structure:"
find . -type f -name "*.go" | sort 