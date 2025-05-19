package sync

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Dirs synchronizes the contents of the source directory (srcDir) with the
// destination directory (dstDir). It ensures that all files present in the
// source are copied to the destination, and optionally deletes files from
// the destination that are missing in the source if deleteMissing is true.
// If dstDir does not exist, it will be created.
//
// Parameters:
//   - srcDir: Path to the source directory (must exist and be valid).
//   - dstDir: Path to the destination directory (created if missing).
//   - deleteMissing: Whether to delete files in dstDir that are absent in srcDir.
//
// Returns an error if synchronization fails for any reason.
func Dirs(srcDir, dstDir string, deleteMissing bool) error {
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return fmt.Errorf("creating destination directory: %w", err)
	}

	dstFiles := make(map[string]os.FileInfo)

	// Walk through the destination directory to collect all files and directories
	err := filepath.Walk(dstDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == dstDir {
			return nil
		}
		relPath, err := filepath.Rel(dstDir, path)
		if err != nil {
			return err
		}
		dstFiles[relPath] = info
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking destination folder: %w", err)
	}

	// Walk through the source directory to synchronize contents
	err = filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking source folder: %w", err)
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil // Skip symbolic links
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(dstPath, 0o755)
		}

		srcInfo, err := d.Info()
		if err != nil {
			return err
		}

		dstInfo, exists := dstFiles[relPath]
		delete(dstFiles, relPath)

		if exists {
			if needsUpdate(srcInfo, dstInfo) {
				return copyFile(path, dstPath)
			}
		} else {
			return copyFile(path, dstPath)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if deleteMissing {
		// Delete files and directories in the destination that are not in the source
		for relPath, info := range dstFiles {
			dstPath := filepath.Join(dstDir, relPath)
			if err := os.RemoveAll(dstPath); err != nil {
				return fmt.Errorf("failed to delete %s: %w", dstPath, err)
			}
			// Ensure directories are deleted properly
			if info.IsDir() {
				if err := os.RemoveAll(dstPath); err != nil {
					return fmt.Errorf("failed to delete directory %s: %w", dstPath, err)
				}
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.Chtimes(dst, time.Now(), srcInfo.ModTime()); err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}

func needsUpdate(src, dst os.FileInfo) bool {
	return src.Size() != dst.Size() || !src.ModTime().Equal(dst.ModTime())
}
