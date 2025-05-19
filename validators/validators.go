package validators

import (
	"fmt"
	"os"
)

func ValidateSrcDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %w", err)
		}
		return fmt.Errorf("getting stats for path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}
	return nil
}
