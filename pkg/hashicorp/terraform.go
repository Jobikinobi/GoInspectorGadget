package hashicorp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
)

// SpeechInfraManager manages infrastructure for speech recognition
type SpeechInfraManager struct {
	workingDir  string
	tf          *tfexec.Terraform
	initialized bool
}

// NewSpeechInfraManager creates a new SpeechInfraManager
func NewSpeechInfraManager(tfPath, workingDir string) (*SpeechInfraManager, error) {
	// Validate terraform path
	if tfPath == "" {
		var err error
		tfPath, err = findTerraform()
		if err != nil {
			return nil, fmt.Errorf("terraform executable not found: %w", err)
		}
	}

	// Ensure working directory exists
	if _, err := os.Stat(workingDir); os.IsNotExist(err) {
		if err := os.MkdirAll(workingDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create working directory: %w", err)
		}
	}

	// Initialize terraform
	tf, err := tfexec.NewTerraform(workingDir, tfPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize terraform: %w", err)
	}

	return &SpeechInfraManager{
		workingDir:  workingDir,
		tf:          tf,
		initialized: true,
	}, nil
}

// findTerraform attempts to locate the terraform executable in PATH
func findTerraform() (string, error) {
	// Look for terraform in PATH
	path, err := exec.LookPath("terraform")
	if err != nil {
		return "", fmt.Errorf("terraform not found in PATH: %w", err)
	}
	return path, nil
}

// Initialize initializes the Terraform configuration
func (m *SpeechInfraManager) Initialize(ctx context.Context) error {
	if !m.initialized {
		return fmt.Errorf("terraform manager not initialized")
	}

	// Run terraform init
	return m.tf.Init(ctx, tfexec.Upgrade(true))
}

// Plan creates a Terraform plan
func (m *SpeechInfraManager) Plan(ctx context.Context) (bool, error) {
	if !m.initialized {
		return false, fmt.Errorf("terraform manager not initialized")
	}

	// Run terraform plan
	return m.tf.Plan(ctx)
}

// Apply applies the Terraform configuration
func (m *SpeechInfraManager) Apply(ctx context.Context) error {
	if !m.initialized {
		return fmt.Errorf("terraform manager not initialized")
	}

	// Run terraform apply
	return m.tf.Apply(ctx)
}

// Destroy destroys the infrastructure
func (m *SpeechInfraManager) Destroy(ctx context.Context) error {
	if !m.initialized {
		return fmt.Errorf("terraform manager not initialized")
	}

	// Run terraform destroy
	return m.tf.Destroy(ctx)
}

// GetOutput retrieves an output from the Terraform state
func (m *SpeechInfraManager) GetOutput(ctx context.Context, name string) (string, error) {
	if !m.initialized {
		return "", fmt.Errorf("terraform manager not initialized")
	}

	// Get all outputs
	outputs, err := m.tf.Output(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get terraform outputs: %w", err)
	}

	// Find the specific output
	output, ok := outputs[name]
	if !ok {
		return "", fmt.Errorf("output '%s' not found", name)
	}

	// Extract value
	var value string
	if err := json.Unmarshal(output.Value, &value); err != nil {
		return "", fmt.Errorf("failed to unmarshal output value: %w", err)
	}

	return value, nil
}

// CreateSpeechRecognitionInfra sets up infrastructure for speech recognition
func (m *SpeechInfraManager) CreateSpeechRecognitionInfra(ctx context.Context, configPath string) error {
	if !m.initialized {
		return fmt.Errorf("terraform manager not initialized")
	}

	// Copy configuration files if provided
	if configPath != "" {
		if err := copyTerraformConfig(configPath, m.workingDir); err != nil {
			return fmt.Errorf("failed to copy terraform configuration: %w", err)
		}
	} else {
		// Create default configuration
		if err := createDefaultSpeechInfraConfig(m.workingDir); err != nil {
			return fmt.Errorf("failed to create default configuration: %w", err)
		}
	}

	// Initialize, plan, and apply
	if err := m.Initialize(ctx); err != nil {
		return fmt.Errorf("terraform init failed: %w", err)
	}

	if _, err := m.Plan(ctx); err != nil {
		return fmt.Errorf("terraform plan failed: %w", err)
	}

	if err := m.Apply(ctx); err != nil {
		return fmt.Errorf("terraform apply failed: %w", err)
	}

	return nil
}

// copyTerraformConfig copies terraform configuration files to the working directory
func copyTerraformConfig(sourcePath, destPath string) error {
	// Check if source path exists
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("source path error: %w", err)
	}

	// Copy differently based on whether source is a file or directory
	if sourceInfo.IsDir() {
		// Copy all .tf and .tfvars files
		entries, err := os.ReadDir(sourcePath)
		if err != nil {
			return fmt.Errorf("failed to read source directory: %w", err)
		}

		for _, entry := range entries {
			if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".tf") || strings.HasSuffix(entry.Name(), ".tfvars")) {
				sourcefile := filepath.Join(sourcePath, entry.Name())
				destfile := filepath.Join(destPath, entry.Name())

				data, err := os.ReadFile(sourcefile)
				if err != nil {
					return fmt.Errorf("failed to read file %s: %w", sourcefile, err)
				}

				if err := os.WriteFile(destfile, data, 0644); err != nil {
					return fmt.Errorf("failed to write file %s: %w", destfile, err)
				}
			}
		}
	} else {
		// Copy the single file
		data, err := os.ReadFile(sourcePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", sourcePath, err)
		}

		destfile := filepath.Join(destPath, filepath.Base(sourcePath))
		if err := os.WriteFile(destfile, data, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", destfile, err)
		}
	}

	return nil
}

// createDefaultSpeechInfraConfig creates a default terraform configuration for speech recognition
func createDefaultSpeechInfraConfig(destPath string) error {
	// Create main.tf
	mainTf := `
provider "aws" {
  region = var.aws_region
}

# S3 bucket for media storage
resource "aws_s3_bucket" "media_bucket" {
  bucket = var.media_bucket_name
  acl    = "private"

  tags = {
    Name        = "MediaProcessingBucket"
    Environment = var.environment
    Project     = "SpeechRecognition"
  }
}

# IAM role for speech processing
resource "aws_iam_role" "speech_processing_role" {
  name = "speech_processing_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

# Lambda function for speech processing
resource "aws_lambda_function" "speech_processor" {
  function_name = "speech_processor"
  role          = aws_iam_role.speech_processing_role.arn
  handler       = "index.handler"
  runtime       = "nodejs14.x"
  timeout       = 300
  memory_size   = 1024

  environment {
    variables = {
      MEDIA_BUCKET = aws_s3_bucket.media_bucket.bucket
    }
  }
}

output "media_bucket_name" {
  value = aws_s3_bucket.media_bucket.bucket
}

output "lambda_function_name" {
  value = aws_lambda_function.speech_processor.function_name
}
`

	// Create variables.tf
	variablesTf := `
variable "aws_region" {
  description = "AWS region for resources"
  default     = "us-west-2"
}

variable "media_bucket_name" {
  description = "Name of the S3 bucket for media storage"
  default     = "speech-recognition-media"
}

variable "environment" {
  description = "Deployment environment"
  default     = "dev"
}
`

	// Write the files
	if err := os.WriteFile(filepath.Join(destPath, "main.tf"), []byte(mainTf), 0644); err != nil {
		return fmt.Errorf("failed to write main.tf: %w", err)
	}

	if err := os.WriteFile(filepath.Join(destPath, "variables.tf"), []byte(variablesTf), 0644); err != nil {
		return fmt.Errorf("failed to write variables.tf: %w", err)
	}

	return nil
}
