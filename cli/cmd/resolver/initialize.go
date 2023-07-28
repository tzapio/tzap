package resolver

import (
	"os"

	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
)

// initializeTzap initializes a Tzap instance with the provided configuration.
// It returns a pointer to the Tzap instance and an error, if any.
func initializeTzap() (*tzap.Tzap, error) {
	config := config.Configuration{
		OpenAIModel:   openai.GPT16,
		AutoMode:      false, // automode == yes
		TruncateLimit: 0,
		MD5Rewrites:   true,
		EnableLogs:    true,
		LoggerOutput:  ".tzap-data/logs",
		Temperature:   0,
		EmbeddingURL:  "",
		CompletionURL: "",
	}

	var connector types.TzapConnector

	apikey, err := tzapconnect.LoadOPENAI_API_KEY()
	if err != nil {
		return nil, err
	}
	connector = tzapconnect.WithConfig(apikey, config)

	t := tzap.NewWithConnector(connector)

	return t, nil
}

// initializeSession sets up the initial session and returns a pointer to the Tzap instance.
// It handles the setup based on the environment and project settings.
func initializeSession() *tzap.Tzap {
	tl.Logger.Println("Tzap initialized")

	// Search for .tzapinclude and get the root directory
	baseDir, err := cmdutil.SearchForTzapincludeAndGetRootDir()
	if err != nil {
		println("Warning: No .tzapinclude file found. Run 'tzap init' Using current directory as root.", err)
	} else {
		os.Chdir(baseDir)
	}
	tl.Logger.Println("Current working directory:", baseDir)

	// Initialize Tzap instance
	t, err := initializeTzap()
	if err != nil {
		panic(err)
	}

	tl.EnableUICompletionLogger()

	return t
}
