package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Step 1: Load config
	fmt.Println("Loading configuration from config.json...")
	cfg, err := loadConfig("config.json")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	fmt.Printf("  Server:     %s\n", cfg.ServerIP)
	fmt.Printf("  Username:   %s\n", cfg.Username)
	fmt.Printf("  Local dir:  %s\n", cfg.LocalDir)
	fmt.Printf("  Remote dir: %s\n", cfg.RemoteDir)
	fmt.Printf("  Prefix:     %s\n\n", cfg.FilePrefix)

	// Step 2: Ensure SSH key exists
	privKeyPath, err := ensureSSHKey()
	if err != nil {
		return fmt.Errorf("ensure SSH key: %w", err)
	}

	signer, err := loadPrivateKey(privKeyPath)
	if err != nil {
		return fmt.Errorf("load private key: %w", err)
	}

	// Step 3: Build tarball in the OS temp directory
	fmt.Printf("\nArchiving %s...\n", cfg.LocalDir)
	archivePath, err := buildTarball(cfg.LocalDir, os.TempDir(), cfg.FilePrefix)
	if err != nil {
		return fmt.Errorf("build tarball: %w", err)
	}
	fmt.Printf("Archive created: %s\n", archivePath)

	// Step 4: Transfer archive to remote server
	fmt.Printf("\nTransferring archive to %s:%s...\n", cfg.ServerIP, cfg.RemoteDir)
	if err := transferFile(archivePath, cfg.ServerIP, cfg.Username, cfg.RemoteDir, signer); err != nil {
		return fmt.Errorf("transfer file: %w", err)
	}

	// Step 5: Cleanup local archive
	fmt.Printf("\nCleaning up local archive %s...\n", archivePath)
	if err := os.Remove(archivePath); err != nil {
		return fmt.Errorf("removing local archive: %w", err)
	}

	fmt.Println("\nBackup completed successfully.")
	return nil
}
