package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// transferFile connects to the remote server via SSH using the provided signer
// and uploads localFile to remoteDir via SFTP.
func transferFile(localFile, serverIP, username, remoteDir string, signer ssh.Signer) error {
	// Print the fingerprint of the key being used so the user can confirm it
	// matches what is in the remote server's authorized_keys.
	pubKey := signer.PublicKey()
	fp := md5.Sum(pubKey.Marshal()) //nolint:gosec
	fingerprint := fmt.Sprintf("%x", fp)
	fmt.Printf("Using key fingerprint (MD5): %s\n", fingerprint)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		// HostKeyCallback accepts any host key. For production use, replace
		// with ssh.FixedHostKey or a known_hosts-based callback.
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec
		BannerCallback: func(banner string) error {
			fmt.Printf("Server banner: %s\n", banner)
			return nil
		},
		Timeout: 15 * time.Second,
	}

	addr := fmt.Sprintf("%s:22", serverIP)
	fmt.Printf("Connecting to %s@%s...\n", username, addr)

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		// Distinguish common failure modes to give the user actionable output.
		var netErr net.Error
		if os.IsTimeout(err) || (asNetError(err, &netErr) && netErr.Timeout()) {
			return fmt.Errorf("SSH connection timed out connecting to %s — check the IP and that port 22 is reachable", addr)
		}
		return fmt.Errorf("SSH dial to %s: %w\n\nHint: ensure the public key printed above is in the remote user's ~/.ssh/authorized_keys", addr, err)
	}
	defer conn.Close()

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("creating SFTP client: %w", err)
	}
	defer sftpClient.Close()

	// Ensure the remote directory exists
	if err := sftpClient.MkdirAll(remoteDir); err != nil {
		return fmt.Errorf("creating remote directory %s: %w", remoteDir, err)
	}

	srcFile, err := os.Open(localFile)
	if err != nil {
		return fmt.Errorf("opening local file %s: %w", localFile, err)
	}
	defer srcFile.Close()

	remotePath := filepath.Join(remoteDir, filepath.Base(localFile))
	dstFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("creating remote file %s: %w", remotePath, err)
	}
	defer dstFile.Close()

	written, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("uploading file: %w", err)
	}

	fmt.Printf("Uploaded %d bytes to %s:%s\n", written, serverIP, remotePath)
	return nil
}

func asNetError(err error, target *net.Error) bool {
	e, ok := err.(net.Error)
	if ok {
		*target = e
	}
	return ok
}
