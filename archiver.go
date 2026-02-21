package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// buildTarball creates a .tar.gz archive of srcDir.
// The archive is written to destDir with the naming convention:
// {filePrefix}_{YYYYMMDD}.tar.gz
// Returns the full path to the created archive.
func buildTarball(srcDir, destDir, filePrefix string) (string, error) {
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return "", fmt.Errorf("source directory does not exist: %s", srcDir)
	}

	timestamp := time.Now().Format("20060102")
	archiveName := fmt.Sprintf("%s_%s.tar.gz", filePrefix, timestamp)
	archivePath := filepath.Join(destDir, archiveName)

	outFile, err := os.Create(archivePath)
	if err != nil {
		return "", fmt.Errorf("creating archive file %s: %w", archivePath, err)
	}
	defer outFile.Close()

	gzWriter := gzip.NewWriter(outFile)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	srcDir = filepath.Clean(srcDir)

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Build a relative path for the tar header
		relPath, err := filepath.Rel(filepath.Dir(srcDir), path)
		if err != nil {
			return fmt.Errorf("computing relative path: %w", err)
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return fmt.Errorf("creating tar header for %s: %w", path, err)
		}
		header.Name = relPath

		// Handle symlinks
		if info.Mode()&os.ModeSymlink != 0 {
			linkTarget, err := os.Readlink(path)
			if err != nil {
				return fmt.Errorf("reading symlink %s: %w", path, err)
			}
			header.Linkname = linkTarget
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("writing tar header for %s: %w", path, err)
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("opening file %s: %w", path, err)
		}
		defer f.Close()

		if _, err := io.Copy(tarWriter, f); err != nil {
			return fmt.Errorf("archiving file %s: %w", path, err)
		}

		return nil
	})
	if err != nil {
		// Clean up partial archive on failure
		os.Remove(archivePath)
		return "", fmt.Errorf("walking source directory: %w", err)
	}

	// Flush all writers to ensure the archive is fully written before returning
	if err := tarWriter.Close(); err != nil {
		os.Remove(archivePath)
		return "", fmt.Errorf("finalizing tar archive: %w", err)
	}
	if err := gzWriter.Close(); err != nil {
		os.Remove(archivePath)
		return "", fmt.Errorf("finalizing gzip stream: %w", err)
	}
	if err := outFile.Sync(); err != nil {
		os.Remove(archivePath)
		return "", fmt.Errorf("flushing archive to disk: %w", err)
	}

	return archivePath, nil
}
