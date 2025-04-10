package hashicorp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// CredentialManager provides secure credential management for speech services
type CredentialManager struct {
	credentialFile string
	credentials    map[string]map[string]string
	initialized    bool
}

// NewCredentialManager creates a new CredentialManager
func NewCredentialManager(credentialFile string) (*CredentialManager, error) {
	if credentialFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		credentialFile = filepath.Join(homeDir, ".media-processor", "credentials.json")
	}

	cm := &CredentialManager{
		credentialFile: credentialFile,
		credentials:    make(map[string]map[string]string),
	}

	// Ensure directory exists
	dir := filepath.Dir(credentialFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create credential directory: %w", err)
	}

	// Load credentials if file exists
	if _, err := os.Stat(credentialFile); err == nil {
		if err := cm.loadCredentials(); err != nil {
			return nil, err
		}
	}

	cm.initialized = true
	return cm, nil
}

// GetCredential retrieves a credential
func (c *CredentialManager) GetCredential(engineName, key string) (string, error) {
	if !c.initialized {
		return "", fmt.Errorf("credential manager not initialized")
	}

	// Check if engine exists
	engine, ok := c.credentials[engineName]
	if !ok {
		return "", fmt.Errorf("engine '%s' not found", engineName)
	}

	// Check if key exists
	value, ok := engine[key]
	if !ok {
		return "", fmt.Errorf("key '%s' not found in engine '%s'", key, engineName)
	}

	return value, nil
}

// StoreCredential stores a credential
func (c *CredentialManager) StoreCredential(engineName, key, value string) error {
	if !c.initialized {
		return fmt.Errorf("credential manager not initialized")
	}

	// Create engine if it doesn't exist
	if _, ok := c.credentials[engineName]; !ok {
		c.credentials[engineName] = make(map[string]string)
	}

	// Store value
	c.credentials[engineName][key] = value

	// Save to file
	return c.saveCredentials()
}

// GetAPIKey is a convenience method to retrieve API keys for speech services
func (c *CredentialManager) GetAPIKey(service string) (string, error) {
	return c.GetCredential("speech-services", service+"-api-key")
}

// FromEnvironment creates a CredentialManager from environment variables
func FromEnvironment() (*CredentialManager, error) {
	credentialFile := os.Getenv("CREDENTIAL_FILE")
	return NewCredentialManager(credentialFile)
}

// loadCredentials loads credentials from file
func (c *CredentialManager) loadCredentials() error {
	data, err := os.ReadFile(c.credentialFile)
	if err != nil {
		return fmt.Errorf("failed to read credential file: %w", err)
	}

	if err := json.Unmarshal(data, &c.credentials); err != nil {
		return fmt.Errorf("failed to parse credential file: %w", err)
	}

	return nil
}

// saveCredentials saves credentials to file
func (c *CredentialManager) saveCredentials() error {
	data, err := json.MarshalIndent(c.credentials, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	if err := os.WriteFile(c.credentialFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write credential file: %w", err)
	}

	return nil
}
