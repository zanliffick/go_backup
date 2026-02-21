Project Overview
GoBackup is a CLI utility written in Go designed to automate the creation of compressed local backups and their secure transmission to a remote server.

Core Workflow
Initialization: Check for an existing RSA SSH key; generate one if missing.

Configuration: Parse a configjson file for paths and credentials.

Compression: Archive the target directory into a .tar.gz format with a timestamped filename.

Transmission: Securely transfer the archive to the remote server via SCP using the generated RSA key.

Cleanup: Remove the local temporary archive after a successful transfer.

1. Technical Stack
Language: Go (Golang)

Compression: archive/tar, compress/gzip

Security/Transport: crypto/rsa, golang.org/x/crypto/ssh

Format: JSON (Configuration)

2. Configuration Schema
The application expects a config.json in the root directory:

JSON
{
  "server_ip": "192.168.1.100",
  "username": "backup_user",
  "local_dir": "/home/user/documents",
  "remote_dir": "/backups/storage",
  "file_prefix": "doc_backup"
}
3. Component Specifications
A. SSH Key Generator
On startup, the app checks ~/.ssh/id_rsa_gobackup. If not found:

Generate a 2048-bit RSA key.

Save the private key locally for the SCP client.

Output the public key to the console, instructing the user to add it to the remote server's authorized_keys.

B. Archiver Logic
The backup file must follow this naming convention:
{file_prefix}_{YYYYMMDD}.tar.gz

C. SSH/SCP Transport
The application will utilize the x/crypto/ssh package to:

Establish a connection using the generated RSA private key.

Open an SCP session (or use an SFTP subsystem).

Stream the .tar.gz file to the remote_dir defined in the config.

4. Error Handling & Safety
Pre-flight Check: Verify local_dir exists before starting compression.

Integrity: Ensure the local .tar.gz is fully written before initiating the upload.

Security: File permissions for the generated RSA key must be restricted to 0600.

5. Execution Flow
Load Config → 2.  Ensure SSH Key → 3.  Build Tarball → 4.  Transfer to Remote → 5.  Exit