#!/bin/bash

# Script to lint all markdown files in the project
# Requires markdownlint-cli to be installed
# npm install -g markdownlint-cli

# Check if markdownlint is installed
if ! command -v markdownlint &> /dev/null; then
  echo "markdownlint could not be found. Please install it with:"
  echo "npm install -g markdownlint-cli"
  exit 1
fi

# Find all markdown files and run markdownlint
echo "Checking markdown files..."

# Using find to get all markdown files
FILES=$(find . -type f -name "*.md" -not -path "./node_modules/*" -not -path "./.git/*")

if [ -z "$FILES" ]; then
  echo "No markdown files found"
  exit 0
fi

# Run markdownlint with our config file
markdownlint -c .markdownlintrc $FILES

# Check the exit code
if [ $? -eq 0 ]; then
  echo "✅ All markdown files look good!"
else
  echo "❌ Some issues were found in markdown files. Please fix them."
  exit 1
fi 