package cmdutil

import (
	"fmt"
	"os"
	"path/filepath"
)

func SearchForTzapincludeAndChangeDir() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get the current directory: %v", err)
	}
	cur := cwd
	for {

		includePath := filepath.Join(cur, ".tzapinclude")
		if _, err := os.Stat(includePath); err == nil {
			break
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("unexpected error when checking for .tzapinclude in %s: %v", cur, err)
		}

		if cur == filepath.Dir(cur) {
			panic("Reached the top of the directory tree. Did not find .tzapinclude in any parent directory. Did you run 'tzap init'?")
		}
		fmt.Printf("Did not find .tzapinclude in: %s\n", cur)
		cur = filepath.Dir(cur)
	}

	if cur != cwd {
		fmt.Printf("Changing working directory to: %s\n", cur)

		if err := os.Chdir(cur); err != nil {
			return fmt.Errorf("failed to change the directory to %s: %v", cur, err)
		}
	}

	return nil
}
