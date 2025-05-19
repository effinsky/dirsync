package sync_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dirsync/sync"
)

type check func(srcDir string, dstDir string, err error, t *testing.T)

func hasErrorMessage(expMsg string) check {
	return func(srcDir, dstDir string, err error, t *testing.T) {
		if err == nil {
			t.Errorf("Expected an error starting with '%s' but got no error", expMsg)
			return
		}
		if !strings.Contains(err.Error(), expMsg) {
			t.Errorf(
				"Mismatch: expected error message to contain '%s', but got '%s'",
				expMsg, err,
			)
		}
	}
}

func hasNoError() check {
	return func(srcDir, dstDir string, err error, t *testing.T) {
		if err != nil {
			t.Errorf("Expected no error; got: %v", err)
		}
	}
}

func assertFileExistsInDst(filename string) check {
	return func(srcDir, dstDir string, err error, t *testing.T) {
		dstFile := filepath.Join(dstDir, filename)
		if !fileExists(dstFile) {
			t.Errorf("Expected file to exist: %s", dstFile)
		}
	}
}

func assertFileMissingInDst(filename string) check {
	return func(srcDir, dstDir string, err error, t *testing.T) {
		dstFile := filepath.Join(dstDir, filename)
		if fileExists(dstFile) {
			t.Errorf("Expected file to be deleted: %s", dstFile)
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func TestSync(t *testing.T) {
	testcases := []struct {
		name          string
		setup         func() (string, string, error)
		deleteMissing bool
		checks        []check
	}{
		{
			name: "source directory does not exist",
			setup: func() (string, string, error) {
				srcDir := filepath.Join(os.TempDir(), "nonexistent_dir_test")
				dstDir, err := os.MkdirTemp("", "dst")
				if err != nil {
					return "", "", err
				}
				return srcDir, dstDir, nil
			},
			deleteMissing: false,
			checks: []check{
				hasErrorMessage("error walking source folder:"),
			},
		},
		{
			name: "destination directory does not exist",
			setup: func() (string, string, error) {
				srcDir, err := os.MkdirTemp("", "src")
				if err != nil {
					return "", "", err
				}
				dstDir := filepath.Join(os.TempDir(), "nonexistent_dst_test")
				err = os.WriteFile(filepath.Join(srcDir, "testfile.txt"), []byte("Hello, world!"), 0o644)
				if err != nil {
					return "", "", err
				}
				return srcDir, dstDir, nil
			},
			deleteMissing: false,
			checks: []check{
				hasNoError(),
				assertFileExistsInDst("testfile.txt"),
			},
		},
		{
			name: "file copied from source to destination",
			setup: func() (string, string, error) {
				srcDir, err := os.MkdirTemp("", "src")
				if err != nil {
					return "", "", err
				}
				dstDir, err := os.MkdirTemp("", "dst")
				if err != nil {
					return "", "", err
				}
				err = os.WriteFile(filepath.Join(srcDir, "file.txt"), []byte("content"), 0o644)
				if err != nil {
					return "", "", err
				}
				return srcDir, dstDir, nil
			},
			deleteMissing: false,
			checks: []check{
				hasNoError(),
				assertFileExistsInDst("file.txt"),
				assertFileExistsInDst("file.txt"),
			},
		},
		{
			name: "file deleted in destination if not in source and deleteMissing=true",
			setup: func() (string, string, error) {
				srcDir, err := os.MkdirTemp("", "src")
				if err != nil {
					return "", "", err
				}
				dstDir, err := os.MkdirTemp("", "dst")
				if err != nil {
					return "", "", err
				}
				err = os.WriteFile(filepath.Join(dstDir, "orphan.txt"), []byte("orphan"), 0o644)
				if err != nil {
					return "", "", err
				}
				return srcDir, dstDir, nil
			},
			deleteMissing: true,
			checks: []check{
				hasNoError(),
				assertFileMissingInDst("orphan.txt"),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			srcDir, dstDir, err := tc.setup()
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			defer func() {
				if err := os.RemoveAll(srcDir); err != nil {
					t.Errorf("Failed to remove source test dirs: %v\n", err)
				}
				if err := os.RemoveAll(dstDir); err != nil {
					t.Errorf("Failed to remove destination test dirs: %v\n", err)
				}
			}()

			err = sync.Dirs(srcDir, dstDir, tc.deleteMissing)

			for _, check := range tc.checks {
				check(srcDir, dstDir, err, t)
			}
		})
	}
}
