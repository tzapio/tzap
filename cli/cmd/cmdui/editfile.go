package cmdui

import (
	"os"
	"os/exec"

	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
	"github.com/tzapio/tzap/pkg/util/stdin"
)

func (ui *CMDUI) EditFile(compareToFile *actionpb.FileWrite) error {
	// create temp file to edit.
	temp, err := os.CreateTemp("", "tzapcontent*")
	if err != nil {
		return err
	}

	if _, err := temp.Write([]byte(compareToFile.Contentout)); err != nil {
		return err
	}

	if _, err := os.Stat(compareToFile.Fileout); err == nil {
		exec.Command("code", "-d", compareToFile.Fileout, temp.Name()).Run()
	} else {
		exec.Command("code", temp.Name()).Run()
	}
	contentOut := editLoop(compareToFile.Contentout, temp.Name(), compareToFile.Fileout)
	if err := os.WriteFile(compareToFile.Fileout, []byte(contentOut), 0755); err != nil {
		return err
	}
	exec.Command("code", compareToFile.Fileout).Run()
	return nil
}

func editLoop(changes, file string, fileOut string) string {
	for {
		println("\n\nFile to edit for changes: ", file)
		println("")
		println("\nFile that will be saved to: ", fileOut, "\n")
		key := stdin.GetStdinInput("Edit file at file location.\n\n - press c and enter to open in vscode. \n - press v and enter to open in vim. \n - press enter to continue. \n\n")
		if key == "v" {
			// open vim
			cmd := exec.Command("vim", file)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stderr
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				panic(err)
			}
			key = ""
		}
		if key == "c" {
			// open code
			exec.Command("code", file).Run()
		}
		if key == "" || key == "y" {
			bytes, err := os.ReadFile(file)
			if err != nil {
				panic(err)
			}
			return string(bytes)
		}
		if key == "e" || key == "exit" {
			panic("Aborting")
		}
	}
}
