package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

const keyPath = "~/.ssh/id_rsa_gobackup"

func expandTilde(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return filepath.Join(home, path[1:])
}

// ensureSSHKey checks for an existing RSA key at ~/.ssh/id_rsa_gobackup.
// If not found, it generates a new 2048-bit RSA key pair, saves the private key
// with 0600 permissions, and prints the public key for the user to add to
// the remote server's authorized_keys.
// Returns the path to the private key file.
func ensureSSHKey() (string, error) {
	privKeyPath := expandTilde(keyPath)
	pubKeyPath := privKeyPath + ".pub"

	if _, err := os.Stat(privKeyPath); err == nil {
		fmt.Printf("SSH key found at %s\n", privKeyPath)
		fmt.Printf("Public key to authorize on remote server (%s):\n", pubKeyPath)
		if pub, err := os.ReadFile(pubKeyPath); err == nil {
			fmt.Printf("  %s", string(pub))
		}
		return privKeyPath, nil
	}

	fmt.Println("No SSH key found. Generating a new 2048-bit RSA key pair...")

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", fmt.Errorf("generating RSA key: %w", err)
	}

	// Encode private key to PEM
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	})

	// Ensure ~/.ssh directory exists
	sshDir := filepath.Dir(privKeyPath)
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return "", fmt.Errorf("creating .ssh directory: %w", err)
	}

	// Write private key with restricted permissions (0600)
	if err := os.WriteFile(privKeyPath, privPEM, 0600); err != nil {
		return "", fmt.Errorf("writing private key: %w", err)
	}

	// Generate and write public key
	pubKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", fmt.Errorf("generating public key: %w", err)
	}
	pubKeyBytes := ssh.MarshalAuthorizedKey(pubKey)

	if err := os.WriteFile(pubKeyPath, pubKeyBytes, 0644); err != nil {
		return "", fmt.Errorf("writing public key: %w", err)
	}

	fmt.Printf("\nNew SSH key generated and saved to: %s\n\n", privKeyPath)
	fmt.Println("Add the following public key to the remote server's ~/.ssh/authorized_keys:")
	fmt.Println("---")
	fmt.Print(string(pubKeyBytes))
	fmt.Println("---")

	return privKeyPath, nil
}

// loadPrivateKey reads and parses the RSA private key at the given path.
func loadPrivateKey(path string) (ssh.Signer, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading private key %s: %w", path, err)
	}

	signer, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return nil, fmt.Errorf("parsing private key: %w", err)
	}

	return signer, nil
}
