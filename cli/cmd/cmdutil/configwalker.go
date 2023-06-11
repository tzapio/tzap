package cmdutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func SearchForTzapincludeAndGetRootDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get the current directory: %v", err)
	}
	cur := cwd
	for {

		includePath := filepath.Join(cur, ".tzapinclude")
		if _, err := os.Stat(includePath); err == nil {
			break
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("unexpected error when checking for .tzapinclude in %s: %v", cur, err)
		}

		if cur == filepath.Dir(cur) {
			return cwd, errors.New("reached the top of the directory tree. Did not find .tzapinclude in any parent directory. Did you run 'tzap init'?")
		}
		println("Did not find .tzapinclude in: " + cur)
		cur = filepath.Dir(cur)
	}

	if cur != cwd {
		println("Changing working directory to: " + cur)
		return cur, nil
	}

	return cwd, nil
}
