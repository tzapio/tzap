package cmdui

import (
	"os"
	"os/exec"
	"path"

	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
	"github.com/tzapio/tzap/pkg/util"
	"github.com/tzapio/tzap/pkg/util/stdin"
)

func (ui *CMDUI) EditFile(compareToFile *actionpb.FileWrite) error {
	// create temp file to edit.
	err := os.MkdirAll("./.tzap-data/edit", 0755)
	if err != nil {
		return err
	}
	tempFolder, err := os.MkdirTemp("./.tzap-data/edit/", "*")
	if err != nil {
		return err
	}

	tempFile := path.Join(tempFolder, path.Base(compareToFile.Fileout))
	if err := os.WriteFile(tempFile, []byte(compareToFile.Contentout), 0600); err != nil {
		return err
	}

	if ui.editor == "vscode" {
		if _, err := os.Stat(compareToFile.Fileout); err == nil {
			exec.Command("code", "-d", compareToFile.Fileout, tempFile).Run()
		} else {
			exec.Command("code", tempFile).Run()
		}
	}
	contentOut := editLoop(compareToFile.Contentout, tempFile, compareToFile.Fileout)
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := os.RemoveAll(tempFolder); err != nil {
		println(err)
	}

	println("Created file", compareToFile.Fileout)
	if err := util.MkdirPAndWriteFile(path.Join(cwd, compareToFile.Fileout), contentOut); err != nil {
		return err
	}
	if ui.editor == "vscode" {
		exec.Command("code", compareToFile.Fileout).Run()
	}
	return nil
}

func editLoop(changes, file string, fileOut string) string {
	for {
		println("\nProposing changes for: ", fileOut)
		println("Changes can be edited at: ", file, "\n")
		key := stdin.GetStdinInput("Edit file at file location.\n\n - press c and enter to open in vscode. \n - press v and enter to open in vim.  \n - press q or ctrl + c to exit and ignore changes \n - press y and enter to continue. \n\n")
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
		if key == "y" {
			bytes, err := os.ReadFile(file)
			if err != nil {
				panic(err)
			}
			return string(bytes)
		}
		if key == "q" || key == "exit" {
			println("Aborting")
			os.Exit(0)
		}
	}
}
