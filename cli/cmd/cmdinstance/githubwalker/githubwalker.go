package githubwalker

// root dir. https://api.github.com/repos/tzapio/tzap/contents
// list all files in root dir.
// if file: get contents: https://raw.githubusercontent.com/tzapio/tzap/main/README.md
// if dir, traverse

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/tzapio/tzap/cli/cmd/cmdutil/fileevaluator"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/types"
)

type GitHubWalker struct {
	baseURL string
	client  *http.Client
	e       *fileevaluator.FileEvaluator
}

func New(baseURL string, e *fileevaluator.FileEvaluator) *GitHubWalker {
	return &GitHubWalker{baseURL: baseURL, client: &http.Client{}, e: e}
}

func (g *GitHubWalker) walkDir(dir string) ([]types.FileReader, error) {
	url := fmt.Sprintf("%s/contents/%s", g.baseURL, dir)
	tl.Logger.Printf("Walking dir: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %v", err)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var contents []struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Type string `json:"type"`
		URL  string `json:"url"`
	}

	if err := json.Unmarshal(body, &contents); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	var fileReaders []types.FileReader

	for _, content := range contents {
		if content.Type == "dir" && g.e.ShouldTraverseDir(content.Path) {
			subpath := path.Join(dir, content.Name)
			subreaders, err := g.walkDir(subpath)
			if err != nil {
				return nil, fmt.Errorf("failed to walk subdirectory %s: %v", subpath, err)
			}
			fileReaders = append(fileReaders, subreaders...)
		} else if content.Type == "file" && g.e.ShouldKeepPath(content.Path) {

			fileURL := content.URL
			fileReader, err := g.getFile(fileURL)
			if err != nil {
				return nil, fmt.Errorf("failed to get file %s: %v", fileURL, err)
			}
			fileReaders = append(fileReaders, fileReader)
		}
	}

	return fileReaders, nil
}

func (g *GitHubWalker) getFile(fileURL string) (types.FileReader, error) {
	tl.Logger.Printf("Getting file: %s\n", fileURL)

	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %v", err)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	content := string(body)

	return &GitHubFileReader{content: content}, nil
}

type GitHubFileReader struct {
	name    string
	content string
}

// FilePath implements types.FileReader.
func (*GitHubFileReader) FilePath() string {
	panic("unimplemented")
}

// Open implements types.FileReader.
func (*GitHubFileReader) Open() (io.ReadCloser, error) {
	panic("unimplemented")
}

// Stat implements types.FileReader.
func (*GitHubFileReader) Stat() (fs.FileInfo, error) {
	panic("unimplemented")
}

func (f *GitHubFileReader) Read() (string, error) {
	return f.content, nil
}

func (g *GitHubWalker) GetFiles() ([]types.FileReader, error) {
	return g.walkDir("")
}
