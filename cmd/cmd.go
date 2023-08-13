package cmd

import (
	"fmt"
	"os"
	"strings"

	"errors"

	"github.com/totmicro/atlantis-yaml-generator/pkg/atlantis"
	"github.com/totmicro/atlantis-yaml-generator/pkg/github"
	"github.com/totmicro/atlantis-yaml-generator/pkg/helpers"
	"github.com/totmicro/atlantis-yaml-generator/pkg/version"

	"github.com/spf13/cobra"
)

const (
	defaultOutputFile        = "atlantis.yaml"
	defaultAutomerge         = "false"
	defaultParallelPlan      = "true"
	defaultParallelApply     = "true"
	defaultAtlantisTfRootDir = "./"
)

var (
	rootCmd = &cobra.Command{
		Use:           "atlantis-yaml-generator",
		Short:         fmt.Sprintf("Atlantis Yaml Generator tool (version %s)", version.VERSION),
		RunE:          genAtlantisYaml,
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version.VERSION,
	}
)

func Init() {

	rootCmd.PersistentFlags().StringP("atlantis-automerge", "", "", `Atlantis automerge config value.
Alternatively, you can set this parameter using ATLANTIS_AUTOMERGE environment variable.`)
	rootCmd.PersistentFlags().StringP("atlantis-parallel-apply", "", "", `Atlantis parallel apply config value.
Alternatively, you can set this parameter using ATLANTIS_PARALLEL_APPLY environment variable.`)
	rootCmd.PersistentFlags().StringP("atlantis-parallel-plan", "", "", `Atlantis parallel plan config value.
Alternatively, you can set this parameter using ATLANTIS_PARALLEL_PLAN environment variable.`)
	rootCmd.PersistentFlags().StringP("output-file", "f", "", `Atlantis output file name.
Alternatively, you can set this parameter using OUTPUT_FILE environment variable.`)
	rootCmd.PersistentFlags().StringP("atlantis-workflow", "w", "", `Atlantis Workflow to be used. [single-workspace|multi-workspace].
Alternatively, you can set this parameter using ATLANTIS_WORKFLOW environment variable.`)
	rootCmd.PersistentFlags().StringP("atlantis-when-modified", "m", "", `Atlantis When modified (list of strings) to run autoplan.
Alternatively, you can set this parameter using ATLANTIS_WHEN_MODIFIED environment variable.`)
	rootCmd.PersistentFlags().StringP("atlantis-excluded-projects", "x", "", `Atlantis regex filter to exclude desired projects.
Alternatively, you can set this parameter using ATLANTIS_EXCLUDED_PROJECTS environment variable.`)
	rootCmd.PersistentFlags().StringP("atlantis-included-projects", "z", "", `Atlantis regex filter to only include desired projects.
Alternatively, you can set this parameter using ATLANTIS_INCLUDED_PROJECTS environment variable.`)
	rootCmd.PersistentFlags().StringP("atlantis-projects-pattern-detector", "q", "", `Atlantis finds this pattern in folders to detect a project.
Alternatively, you can set this parameter using ATLANTIS_PROJECTS_PATTERN_DETECTOR environment variable.`)
	rootCmd.PersistentFlags().StringP("pull-num", "p", "", `Github Pull Number.
Alternatively, you can set this parameter using PULL_NUM environment variable.`)
	rootCmd.PersistentFlags().StringP("base-repo-name", "r", "", `Github Repo Name.
Alternatively, you can set this parameter using BASE_REPO_NAME environment variable.`)
	rootCmd.PersistentFlags().StringP("base-repo-owner", "o", "", `Github Repo Owner Name.
Alternatively, you can set this parameter using BASE_REPO_OWNER environment variable.`)
	rootCmd.PersistentFlags().StringP("gh-token", "t", "", `Github Token Value.
Alternatively, you can set this parameter using GH_TOKEN environment variable.`)
	rootCmd.PersistentFlags().StringP("atlantis-tf-root-dir", "d", "", `Basedir for tf resources.`)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func genAtlantisYaml(ccmd *cobra.Command, args []string) error {
	err := errors.New("")
	// Read flags and environment variables
	atlantisAutomerge := helpers.GetFlagOrEnv(ccmd, "atlantis-automerge", "ATLANTIS_AUTOMERGE", defaultAutomerge)
	atlantisParallelApply := helpers.GetFlagOrEnv(ccmd, "atlantis-parallel-apply", "ATLANTIS_PARALLEL_APPLY", defaultParallelApply)
	atlantisParallelPlan := helpers.GetFlagOrEnv(ccmd, "atlantis-parallel-plan", "ATLANTIS_PARALLEL_PLAN", defaultParallelPlan)
	atlantisTfRootDir := helpers.GetFlagOrEnv(ccmd, "atlantis-tf-root-dir", "ATLANTIS_TF_ROOT_DIR", defaultAtlantisTfRootDir)
	atlantisOutputFile := helpers.GetFlagOrEnv(ccmd, "output-file", "OUTPUT_FILE", defaultOutputFile)
	atlantisWorkflow := helpers.GetFlagOrEnv(ccmd, "atlantis-workflow", "ATLANTIS_WORKFLOW", "")
	atlantisWhenModified := helpers.GetFlagOrEnv(ccmd, "atlantis-when-modified", "ATLANTIS_WHEN_MODIFIED", "")
	atlantisExcludedProjects := helpers.GetFlagOrEnv(ccmd, "atlantis-excluded-projects", "ATLANTIS_EXCLUDED_PROJECTS", "")
	atlantisIncludedProjects := helpers.GetFlagOrEnv(ccmd, "atlantis-included-projects", "ATLANTIS_INCLUDED_PROJECTS", "")
	atlantisProjectsPatternDetector := helpers.GetFlagOrEnv(ccmd, "atlantis-projects-pattern-detector", "ATLANTIS_PROJECTS_PATTERN_DETECTOR", "")
	pullNum := helpers.GetFlagOrEnv(ccmd, "pull-num", "PULL_NUM", "")
	ghRepo := helpers.GetFlagOrEnv(ccmd, "base-repo-name", "BASE_REPO_NAME", "")
	ghRepoOwner := helpers.GetFlagOrEnv(ccmd, "base-repo-owner", "BASE_REPO_OWNER", "")
	ghToken := helpers.GetFlagOrEnv(ccmd, "gh-token", "GH_TOKEN", "")

	// Validate required parameters
	var requiredArgs = []string{"gh-token", "base-repo-owner", "base-repo-name", "pull-num", "atlantis-workflow"}
	var missingArgs []string

	for _, arg := range requiredArgs {
		switch arg {
		case "gh-token":
			if ghToken == "" {
				missingArgs = append(missingArgs, "gh-token")
			}
		case "base-repo-owner":
			if ghRepoOwner == "" {
				missingArgs = append(missingArgs, "base-repo-owner")
			}
		case "base-repo-name":
			if ghRepo == "" {
				missingArgs = append(missingArgs, "base-repo-name")
			}
		case "pull-num":
			if pullNum == "" {
				missingArgs = append(missingArgs, "pull-num")
			}
		case "atlantis-workflow":
			if atlantisWorkflow == "" {
				missingArgs = append(missingArgs, "atlantis-workflow")
			}
		}
	}
	if len(missingArgs) > 0 {
		errorMsg := fmt.Sprintf("missing required parameters: %s", strings.Join(missingArgs, ", "))
		fmt.Println("Use --help for more information.")
		return errors.New(errorMsg)
	}
	// Parse the WhenModified list
	atlantisWhenModifiedList := parseWhenModifiedList(atlantisWhenModified)

	// Create GitHub and Atlantis parameters
	ghParams := github.GithubRequest{
		Repo:              ghRepo,
		Owner:             ghRepoOwner,
		PullRequestNumber: pullNum,
		AuthToken:         ghToken,
	}

	atlantisParams := atlantis.Parameters{
		Automerge:               atlantisAutomerge,
		ParallelApply:           atlantisParallelApply,
		ParallelPlan:            atlantisParallelPlan,
		TfRootDir:               atlantisTfRootDir,
		OutputFile:              atlantisOutputFile,
		Workflow:                atlantisWorkflow,
		WhenModified:            atlantisWhenModifiedList,
		ExcludedProjects:        atlantisExcludedProjects,
		IncludedProjects:        atlantisIncludedProjects,
		ProjectsPatternDetector: atlantisProjectsPatternDetector,
	}

	err = atlantis.GenerateAtlantisYAML(ghParams, atlantisParams)
	return err
}

func parseWhenModifiedList(whenModified string) []string {
	if whenModified == "" {
		return atlantis.DefaultWhenModified
	}
	return strings.Split(whenModified, ",")
}
