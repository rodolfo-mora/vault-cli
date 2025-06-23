package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

// VaultConfig holds the configuration for Vault client
type VaultConfig struct {
	Addr        string
	Token       string
	TransitPath string
	TransitKey  string
	MaxItems    int
}

// VaultClient wraps the Vault API client
type VaultClient struct {
	client *vault.Client
	config *VaultConfig
}

// BulkEncryptRequest represents a bulk encryption request
type BulkEncryptRequest struct {
	Plaintext string `json:"plaintext"`
}

// BulkEncryptResponse represents a bulk encryption response
type BulkEncryptResponse struct {
	Data struct {
		BatchResults []struct {
			Ciphertext string `json:"ciphertext"`
		} `json:"batch_results"`
	} `json:"data"`
}

// BulkDecryptRequest represents a bulk decryption request
type BulkDecryptRequest struct {
	Ciphertext string `json:"ciphertext"`
}

// BulkDecryptResponse represents a bulk decryption response
type BulkDecryptResponse struct {
	Data struct {
		BatchResults []struct {
			Plaintext string `json:"plaintext"`
		} `json:"batch_results"`
	} `json:"data"`
}

// NewVaultClient creates a new Vault client
func NewVaultClient(config *VaultConfig) (*VaultClient, error) {
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = config.Addr

	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	client.SetToken(config.Token)

	return &VaultClient{
		client: client,
		config: config,
	}, nil
}

// EncryptString encrypts a single string
func (vc *VaultClient) EncryptString(text string) error {
	// Base64 encode the input text
	encodedText := base64.StdEncoding.EncodeToString([]byte(text))

	// Prepare the encryption request
	data := map[string]interface{}{
		"plaintext": encodedText,
	}

	// Make the API call
	secret, err := vc.client.Logical().Write(
		fmt.Sprintf("%s/encrypt/%s", vc.config.TransitPath, vc.config.TransitKey),
		data,
	)
	if err != nil {
		return fmt.Errorf("failed to encrypt string: %w", err)
	}

	// Extract the ciphertext
	ciphertext, ok := secret.Data["ciphertext"].(string)
	if !ok {
		return fmt.Errorf("invalid response format: ciphertext not found")
	}

	fmt.Printf("Encrypted text: %s\n", ciphertext)
	return nil
}

// DecryptString decrypts a single string
func (vc *VaultClient) DecryptString(ciphertext string) error {
	// Prepare the decryption request
	data := map[string]interface{}{
		"ciphertext": ciphertext,
	}

	// Make the API call
	secret, err := vc.client.Logical().Write(
		fmt.Sprintf("%s/decrypt/%s", vc.config.TransitPath, vc.config.TransitKey),
		data,
	)
	if err != nil {
		return fmt.Errorf("failed to decrypt string: %w", err)
	}

	// Extract the plaintext
	encodedPlaintext, ok := secret.Data["plaintext"].(string)
	if !ok {
		return fmt.Errorf("invalid response format: plaintext not found")
	}

	// Base64 decode the plaintext
	plaintextBytes, err := base64.StdEncoding.DecodeString(encodedPlaintext)
	if err != nil {
		return fmt.Errorf("failed to decode plaintext: %w", err)
	}

	fmt.Printf("Decrypted text: %s\n", string(plaintextBytes))
	return nil
}

// EncryptFile encrypts a file using bulk operations
func (vc *VaultClient) EncryptFile(inputPath, outputFolder string, override bool) error {
	// Read the input file
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Split content into lines for bulk processing
	lines := strings.Split(string(content), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1] // Remove trailing empty line
	}

	// Process in batches
	var encryptedLines []string
	for i := 0; i < len(lines); i += vc.config.MaxItems {
		end := i + vc.config.MaxItems
		if end > len(lines) {
			end = len(lines)
		}

		batch := lines[i:end]
		encryptedBatch, err := vc.encryptBatch(batch)
		if err != nil {
			return fmt.Errorf("failed to encrypt batch %d-%d: %w", i, end, err)
		}
		encryptedLines = append(encryptedLines, encryptedBatch...)
	}

	// Determine output path
	var outputPath string
	if override {
		outputPath = inputPath
	} else if outputFolder != "" {
		filename := filepath.Base(inputPath)
		outputPath = filepath.Join(outputFolder, filename+".encrypted")
	} else {
		outputPath = inputPath + ".encrypted"
	}

	// Create output directory if it doesn't exist
	if outputFolder != "" && !override {
		if err := os.MkdirAll(outputFolder, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Write encrypted content
	encryptedContent := strings.Join(encryptedLines, "\n")
	if err := os.WriteFile(outputPath, []byte(encryptedContent), 0644); err != nil {
		return fmt.Errorf("failed to write encrypted file: %w", err)
	}

	fmt.Printf("File encrypted successfully: %s\n", outputPath)
	return nil
}

// DecryptFile decrypts a file using bulk operations
func (vc *VaultClient) DecryptFile(inputPath, outputFolder string, override bool) error {
	// Read the input file
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Split content into lines for bulk processing
	lines := strings.Split(string(content), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1] // Remove trailing empty line
	}

	// Process in batches
	var decryptedLines []string
	for i := 0; i < len(lines); i += vc.config.MaxItems {
		end := i + vc.config.MaxItems
		if end > len(lines) {
			end = len(lines)
		}

		batch := lines[i:end]
		decryptedBatch, err := vc.decryptBatch(batch)
		if err != nil {
			return fmt.Errorf("failed to decrypt batch %d-%d: %w", i, end, err)
		}
		decryptedLines = append(decryptedLines, decryptedBatch...)
	}

	// Determine output path
	var outputPath string
	if override {
		outputPath = inputPath
	} else if outputFolder != "" {
		filename := filepath.Base(inputPath)
		// Remove .encrypted extension if present
		if strings.HasSuffix(filename, ".encrypted") {
			filename = strings.TrimSuffix(filename, ".encrypted")
		}
		outputPath = filepath.Join(outputFolder, filename)
	} else {
		// Remove .encrypted extension if present
		if strings.HasSuffix(inputPath, ".encrypted") {
			outputPath = strings.TrimSuffix(inputPath, ".encrypted")
		} else {
			outputPath = inputPath + ".decrypted"
		}
	}

	// Create output directory if it doesn't exist
	if outputFolder != "" && !override {
		if err := os.MkdirAll(outputFolder, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Write decrypted content
	decryptedContent := strings.Join(decryptedLines, "\n")
	if err := os.WriteFile(outputPath, []byte(decryptedContent), 0644); err != nil {
		return fmt.Errorf("failed to write decrypted file: %w", err)
	}

	fmt.Printf("File decrypted successfully: %s\n", outputPath)
	return nil
}

// encryptBatch encrypts a batch of lines using Vault's bulk API
func (vc *VaultClient) encryptBatch(lines []string) ([]string, error) {
	if len(lines) == 0 {
		return []string{}, nil
	}

	// Prepare batch request
	var batchRequests []BulkEncryptRequest
	for _, line := range lines {
		if line == "" {
			continue // Skip empty lines
		}
		// Base64 encode each line
		encodedLine := base64.StdEncoding.EncodeToString([]byte(line))
		batchRequests = append(batchRequests, BulkEncryptRequest{
			Plaintext: encodedLine,
		})
	}

	if len(batchRequests) == 0 {
		return []string{}, nil
	}

	// Prepare the bulk encryption request
	data := map[string]interface{}{
		"batch_input": batchRequests,
	}

	// Make the API call
	secret, err := vc.client.Logical().Write(
		fmt.Sprintf("%s/encrypt/%s", vc.config.TransitPath, vc.config.TransitKey),
		data,
	)
	if err != nil {
		return nil, fmt.Errorf("bulk encryption failed: %w", err)
	}

	// Parse the response
	responseBytes, err := json.Marshal(secret.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var response BulkEncryptResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract ciphertexts
	var results []string
	for _, result := range response.Data.BatchResults {
		results = append(results, result.Ciphertext)
	}

	return results, nil
}

// decryptBatch decrypts a batch of lines using Vault's bulk API
func (vc *VaultClient) decryptBatch(lines []string) ([]string, error) {
	if len(lines) == 0 {
		return []string{}, nil
	}

	// Prepare batch request
	var batchRequests []BulkDecryptRequest
	for _, line := range lines {
		if line == "" {
			continue // Skip empty lines
		}
		batchRequests = append(batchRequests, BulkDecryptRequest{
			Ciphertext: line,
		})
	}

	if len(batchRequests) == 0 {
		return []string{}, nil
	}

	// Prepare the bulk decryption request
	data := map[string]interface{}{
		"batch_input": batchRequests,
	}

	// Make the API call
	secret, err := vc.client.Logical().Write(
		fmt.Sprintf("%s/decrypt/%s", vc.config.TransitPath, vc.config.TransitKey),
		data,
	)
	if err != nil {
		return nil, fmt.Errorf("bulk decryption failed: %w", err)
	}

	// Parse the response
	responseBytes, err := json.Marshal(secret.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var response BulkDecryptResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract and decode plaintexts
	var results []string
	for _, result := range response.Data.BatchResults {
		// Base64 decode the plaintext
		plaintextBytes, err := base64.StdEncoding.DecodeString(result.Plaintext)
		if err != nil {
			return nil, fmt.Errorf("failed to decode plaintext: %w", err)
		}
		results = append(results, string(plaintextBytes))
	}

	return results, nil
}
