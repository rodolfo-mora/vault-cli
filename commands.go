package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Global flags
	vaultAddr    string
	vaultToken   string
	transitPath  string
	transitKey   string
	maxItems     int
	overrideFile bool
	outputFolder string
	inputFile    string
	inputText    string

	// Commands
	encryptCmd = &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt data using Vault transit engine",
		Long:  "Encrypt a string or file using HashiCorp Vault's transit engine",
		RunE:  runEncrypt,
	}

	decryptCmd = &cobra.Command{
		Use:   "decrypt",
		Short: "Decrypt data using Vault transit engine",
		Long:  "Decrypt a string or file using HashiCorp Vault's transit engine",
		RunE:  runDecrypt,
	}
)

func init() {
	// Global flags
	encryptCmd.Flags().StringVar(&vaultAddr, "vault-addr", "", "Vault server address (or use VAULT_ADDR env var)")
	encryptCmd.Flags().StringVar(&vaultToken, "vault-token", "", "Vault token (or use VAULT_TOKEN env var)")
	encryptCmd.Flags().StringVar(&transitPath, "transit-path", "transit", "Transit engine path")
	encryptCmd.Flags().StringVar(&transitKey, "transit-key", "", "Transit key name")
	encryptCmd.Flags().IntVar(&maxItems, "max-items", 100, "Maximum items per bulk operation")
	encryptCmd.Flags().BoolVar(&overrideFile, "override", false, "Override original file")
	encryptCmd.Flags().StringVar(&outputFolder, "output", "", "Output destination folder")
	encryptCmd.Flags().StringVar(&inputFile, "file", "", "File to encrypt")
	encryptCmd.Flags().StringVar(&inputText, "text", "", "Text to encrypt")

	decryptCmd.Flags().StringVar(&vaultAddr, "vault-addr", "", "Vault server address (or use VAULT_ADDR env var)")
	decryptCmd.Flags().StringVar(&vaultToken, "vault-token", "", "Vault token (or use VAULT_TOKEN env var)")
	decryptCmd.Flags().StringVar(&transitPath, "transit-path", "transit", "Transit engine path")
	decryptCmd.Flags().StringVar(&transitKey, "transit-key", "", "Transit key name")
	decryptCmd.Flags().IntVar(&maxItems, "max-items", 100, "Maximum items per bulk operation")
	decryptCmd.Flags().BoolVar(&overrideFile, "override", false, "Override original file")
	decryptCmd.Flags().StringVar(&outputFolder, "output", "", "Output destination folder")
	decryptCmd.Flags().StringVar(&inputFile, "file", "", "File to decrypt")
	decryptCmd.Flags().StringVar(&inputText, "text", "", "Text to decrypt")

	// Mark required flags
	encryptCmd.MarkFlagRequired("transit-key")
	decryptCmd.MarkFlagRequired("transit-key")
}

func getVaultConfig() (*VaultConfig, error) {
	// Use environment variables if not provided via flags
	if vaultAddr == "" {
		vaultAddr = viper.GetString("VAULT_ADDR")
	}
	if vaultToken == "" {
		vaultToken = viper.GetString("VAULT_TOKEN")
	}

	if vaultAddr == "" {
		return nil, fmt.Errorf("vault address is required (use --vault-addr or VAULT_ADDR env var)")
	}
	if vaultToken == "" {
		return nil, fmt.Errorf("vault token is required (use --vault-token or VAULT_TOKEN env var)")
	}

	return &VaultConfig{
		Addr:        vaultAddr,
		Token:       vaultToken,
		TransitPath: transitPath,
		TransitKey:  transitKey,
		MaxItems:    maxItems,
	}, nil
}

func validateInputs() error {
	if inputFile != "" && inputText != "" {
		return fmt.Errorf("cannot specify both --file and --text, choose one")
	}
	if inputFile == "" && inputText == "" {
		return fmt.Errorf("must specify either --file or --text")
	}
	return nil
}

func runEncrypt(cmd *cobra.Command, args []string) error {
	if err := validateInputs(); err != nil {
		return err
	}

	config, err := getVaultConfig()
	if err != nil {
		return err
	}

	client, err := NewVaultClient(config)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	if inputText != "" {
		return client.EncryptString(inputText)
	}

	if inputFile != "" {
		return client.EncryptFile(inputFile, outputFolder, overrideFile)
	}

	return fmt.Errorf("no input specified")
}

func runDecrypt(cmd *cobra.Command, args []string) error {
	if err := validateInputs(); err != nil {
		return err
	}

	config, err := getVaultConfig()
	if err != nil {
		return err
	}

	client, err := NewVaultClient(config)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	if inputText != "" {
		return client.DecryptString(inputText)
	}

	if inputFile != "" {
		return client.DecryptFile(inputFile, outputFolder, overrideFile)
	}

	return fmt.Errorf("no input specified")
}
