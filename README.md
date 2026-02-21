# GoBackup

A CLI utility written in Go that automates the creation of compressed local backups and their secure transmission to a remote server.

## What it does

1. Checks for an RSA SSH key at `~/.ssh/id_rsa_gobackup`, generating one if it doesn't exist
2. Reads `config.json` for source path and remote server credentials
3. Compresses the target directory into a timestamped `.tar.gz` archive
4. Uploads the archive to the remote server via SFTP
5. Removes the local archive after a successful transfer

## Requirements

- Go 1.22 or later
- The remote server must have your public key added to `~/.ssh/authorized_keys`

## Setup

### 1. Clone the repository

```bash
git clone <repo-url>
cd go_backup
```

### 2. Create your config file

```bash
cp config.json.example config.json
```

Edit `config.json` with your values:

```json
{
  "server_ip": "192.168.1.100",
  "username": "backup_user",
  "local_dir": "/home/user/documents",
  "remote_dir": "/backups/storage",
  "file_prefix": "doc_backup"
}
```

| Field | Description |
|---|---|
| `server_ip` | IP address of the remote backup server |
| `username` | SSH username on the remote server |
| `local_dir` | Local directory to back up |
| `remote_dir` | Destination directory on the remote server |
| `file_prefix` | Prefix for the archive filename |

### 3. Authorize the SSH key on the remote server

On the first run, GoBackup will generate an RSA key pair at `~/.ssh/id_rsa_gobackup` and print the public key to the console. Copy that key and append it to `~/.ssh/authorized_keys` on the remote server:

```bash
echo "<printed-public-key>" >> ~/.ssh/authorized_keys
```

## Usage

### Build

```bash
make
```

### Run

```bash
make run
```

The program reads `config.json` from the current directory every time it runs. `make run` will error with a helpful message if the config file is not found.

### Clean

```bash
make clean
```

Removes the compiled `gobackup` binary.

## Archive naming

Archives follow the convention:

```
{file_prefix}_{YYYYMMDD}.tar.gz
```

For example, with `file_prefix` set to `doc_backup` and run on February 21 2026:

```
doc_backup_20260221.tar.gz
```

## SSH key details

- Key location: `~/.ssh/id_rsa_gobackup`
- Key type: RSA 2048-bit
- Private key permissions: `0600`
- If the key already exists at startup, it is reused without regeneration
