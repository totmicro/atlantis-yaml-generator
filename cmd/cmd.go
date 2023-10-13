package cmd

import (
	"fmt"
	"os"

	"github.com/totmicro/atlantis-yaml-generator/pkg/atlantis"
	"github.com/totmicro/atlantis-yaml-generator/pkg/config"

	"github.com/totmicro/atlantis-yaml-generator/pkg/version"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:           "atlantis-yaml-generator",
		Short:         fmt.Sprintf("Atlantis Yaml Generator tool (version %s)", version.GetVersion()),
		RunE:          runE,
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version.GetVersion(),
	}
)

// Init initializes the command line parser and executes the root command.
func Init() {
	initFlags(rootCmd)
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

// initFlags initializes all flags for the root command.
func initFlags(cmd *cobra.Command) {
	for i := range config.ParameterList {
		cmd.PersistentFlags().StringP(
			config.ParameterList[i].Name,
			config.ParameterList[i].Shorthand,
			"",
			config.GenerateDescription(
				config.ParameterList[i].Name,
				config.ParameterList[i].Description))
	}
}

// After command args are parsed in the above call, config.Init() function is invoked to check environment vars.
// This approach enables to define args using either command-line or environment variables.

// runE is the actual execution of the root command.
func runE(ccmd *cobra.Command, args []string) (err error) {
	err = config.Init(ccmd)
	if err != nil {
		return err
	}

	err = atlantis.GenerateAtlantisYAML()
	return err
}
