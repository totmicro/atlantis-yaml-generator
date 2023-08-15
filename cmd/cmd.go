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

type RequiredArgs struct {
	GhToken          string
	BaseRepoOwner    string
	BaseRepoName     string
	PullNum          string
	AtlantisWorkflow string
}

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
	atlantisAutomerge := helpers.GetFlagOrEnv(ccmd, "atlantis-automerge", "ATLANTIS_AUTOMERGE", atlantis.DefaultAutomerge)
	atlantisParallelApply := helpers.GetFlagOrEnv(ccmd, "atlantis-parallel-apply", "ATLANTIS_PARALLEL_APPLY", atlantis.DefaultParallelApply)
	atlantisParallelPlan := helpers.GetFlagOrEnv(ccmd, "atlantis-parallel-plan", "ATLANTIS_PARALLEL_PLAN", atlantis.DefaultParallelPlan)
	atlantisTfRootDir := helpers.GetFlagOrEnv(ccmd, "atlantis-tf-root-dir", "ATLANTIS_TF_ROOT_DIR", atlantis.DefaultAtlantisTfRootDir)
	atlantisOutputFile := helpers.GetFlagOrEnv(ccmd, "output-file", "OUTPUT_FILE", atlantis.DefaultOutputFile)
	atlantisWorkflow := helpers.GetFlagOrEnv(ccmd, "atlantis-workflow", "ATLANTIS_WORKFLOW", "")
	atlantisWhenModified := helpers.GetFlagOrEnv(ccmd, "atlantis-when-modified", "ATLANTIS_WHEN_MODIFIED", "")
	atlantisExcludedProjects := helpers.GetFlagOrEnv(ccmd, "atlantis-excluded-projects", "ATLANTIS_EXCLUDED_PROJECTS", "")
	atlantisIncludedProjects := helpers.GetFlagOrEnv(ccmd, "atlantis-included-projects", "ATLANTIS_INCLUDED_PROJECTS", "")
	atlantisProjectsPatternDetector := helpers.GetFlagOrEnv(ccmd, "atlantis-projects-pattern-detector", "ATLANTIS_PROJECTS_PATTERN_DETECTOR", "")
	pullNum := helpers.GetFlagOrEnv(ccmd, "pull-num", "PULL_NUM", "")
	baseRepoName := helpers.GetFlagOrEnv(ccmd, "base-repo-name", "BASE_REPO_NAME", "")
	baseRepoOwner := helpers.GetFlagOrEnv(ccmd, "base-repo-owner", "BASE_REPO_OWNER", "")
	ghToken := helpers.GetFlagOrEnv(ccmd, "gh-token", "GH_TOKEN", "")

	// Validate required parameter
	var reqArgs = RequiredArgs{
		GhToken:          ghToken,
		BaseRepoOwner:    baseRepoOwner,
		BaseRepoName:     baseRepoName,
		PullNum:          pullNum,
		AtlantisWorkflow: atlantisWorkflow,
	}

	err = helpers.CheckRequiredArgs(reqArgs)
	if err != nil {
		fmt.Println("Run", ccmd.CommandPath(), "--help for more information.")
		return err
	}

	// Define the WhenModified list
	atlantisWhenModifiedList := defineWhenModifiedList(atlantisWhenModified)

	// Define pattern detector
	atlantisProjectsPatternDetector = defineProjectPatternDetector(atlantisProjectsPatternDetector, atlantisWorkflow)

	// Create GitHub and Atlantis parameters
	ghParams := github.GithubRequest{
		Repo:              baseRepoName,
		Owner:             baseRepoOwner,
		PullRequestNumber: pullNum,
		AuthToken:         ghToken,
	}

	atlantisParams := atlantis.Parameters{
		Automerge:        atlantisAutomerge,
		ParallelApply:    atlantisParallelApply,
		ParallelPlan:     atlantisParallelPlan,
		TfRootDir:        atlantisTfRootDir,
		OutputFile:       atlantisOutputFile,
		Workflow:         atlantisWorkflow,
		WhenModified:     atlantisWhenModifiedList,
		ExcludedProjects: atlantisExcludedProjects,
		IncludedProjects: atlantisIncludedProjects,
		PatternDetector:  atlantisProjectsPatternDetector,
	}

	err = atlantis.GenerateAtlantisYAML(ghParams, atlantisParams)
	return err
}

func defineWhenModifiedList(whenModified string) []string {
	if whenModified == "" {
		return atlantis.DefaultWhenModified
	}
	return strings.Split(whenModified, ",")
}

func defineProjectPatternDetector(patternDetector, workflow string) string {
	value, found := atlantis.WorkflowPatternDetectorMap[workflow]
	if !found {
		err := fmt.Errorf("Workflow %s not found", workflow)
		return err.Error()
	}
	if patternDetector == "" {
		return value
	}
	return patternDetector
}
