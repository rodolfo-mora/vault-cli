package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "vault-cli",
		Short: "A CLI tool for HashiCorp Vault transit operations",
		Long: `A comprehensive CLI tool for encrypting and decrypting data using HashiCorp Vault's transit engine.
Supports both string and file operations with bulk processing capabilities.`,
	}

	// Add subcommands
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
