**Atlantis YAML Generator**
--------------------------


`atlantis-yaml-generator` is a command-line tool designed to simplify the generation of Atlantis configuration files (atlantis.yaml).

It automates the process of detecting projects within your Terraform codebase, configuring workspaces, and generating the required configuration for the Atlantis CI/CD tool.

--------
**Features**
------------

Generate `atlantis.yaml` files with ease.
Automatically detect projects and workspaces based on your specified workflow.
Customize configurations such as automerge, parallel plan/apply, and more.
Flexible project inclusion and exclusion using regular expressions.

-----

**Usage**
---------

Run the tool using the following command:
`atlantis-yaml-generator` [flags]


**Available Flags:**



-  `--automerge string`: Atlantis automerge config value. (Equivalent envVar: `AUTOMERGE`)

-  `-r, --base-repo-name string`: Github Repo Name. (Equivalent envVar: `BASE_REPO_NAME`)

-  `-o, --base-repo-owner string`: Github Repo Owner Name. (Equivalent envVar: `BASE_REPO_OWNER`)

-  `-x, --excluded-projects string`: Atlantis regex filter to exclude projects. (Equivalent envVar: `EXCLUDED_PROJECTS`)

-  `-t, --gh-token string`: Github Token Value. (Equivalent envVar: `GH_TOKEN`)

-  `-h, --help`: help for atlantis-yaml-generator

-  `-z, --included-projects string`: Atlantis regex filter to only include projects. (Equivalent envVar: `INCLUDED_PROJECTS`)

-  `-f, --output-file string`: Atlantis output file name. (Equivalent envVar: `OUTPUT_FILE`)

-  `--parallel-apply string`: Atlantis parallel apply config value. (Equivalent envVar: `PARALLEL_APPLY`)

-  `--parallel-plan string`: Atlantis parallel plan config value. (Equivalent envVar: `PARALLEL_PLAN`)

-  `-q, --pattern-detector string`: discover projects based on files or directories names. (Equivalent envVar: `PATTERN_DETECTOR`)

-  `-p, --pull-num string`: Github Pull Request Number to check diffs. (Equivalent envVar: `PULL_NUM`)

-  `--terraform-base-dir string`: Basedir for terraform resources. (Equivalent envVar: `TERRAFORM_BASE_DIR`)

-  `-v, --version`: version for atlantis-yaml-generator

-  `-m, --when-modified string`: Atlantis When modified (list of strings) to run autoplan. (Equivalent envVar: `WHEN_MODIFIED`)

-  `-w, --workflow string`: Atlantis Workflow to be used. `single-workspace|multi-workspace`. (Equivalent envVar: `WORKFLOW`)


*Note that default values and required flags are defined in [`config.go`](pkg/config/config.go)*

--------

**Examples**

*Generate an `atlantis.yaml` file for a multi workspace workflow:*
  ```
  atlantis-yaml-generator  --base-repo-owner spendesk --base-repo-name terraform --pull-num 884 -w multi-workspace --pattern-detector workspace_vars --gh-token ghp_xxx
  ```
*Generate an `atlantis.yaml` file for a single workspace workflow:*
  ```
  atlantis-yaml-generator  --base-repo-owner spendesk --base-repo-name datadog-as-code --pull-num 123 -w single-workspace --pattern-detector main.tf --gh-token ghp_xxx
  ```
*Generate an `atlantis.yaml` file for a multiple workspace workflow with filtering:*
  ```
  atlantis-yaml-generator/atlantis-yaml-generator  --base-repo-owner spendesk --base-repo-name terraform --pull-num 884 -w multi-workspace --pattern-detector workspace_vars --gh-token ghp_xxxx --terraform-base-dir $(pwd) --excluded-projects "(databases-onboarding-tasks-service-production)$"
  ```

*Note that you can also use environment variables to pass the sensitive args*

-------

**Workflows**
-------------

Currenlty `atlantis-yaml-generator` support 2 workflow types:
- `single-workspace`: Intended for Terraform configurations that do not utilize multiple workspaces. In this context, the pattern detector parameters establish the criteria for identifying project folders. For instance, if the pattern-detector is set to main.tf (file), and this file is located at database/dev/main.tf, the resulting project would be labeled as database-env. Consequently, Terraform commands would be executed within the database/dev folder.
- `multiple-workspace`: Intended for Terraform configurations that utilize multiple workspaces. In this context, the pattern detector parameters establish the criteria for identifying project folders. For instance, if the pattern-detector is set to workspace_vars (folder), and there are several files located in this folder, i.e. (database/workspace_vars/dev.tf |database/workspace_vars/staging.tf ), the resulting projects would be labeled as (database-dev|database-staging). Consequently, Terraform commands would be executed within the database folder. Please be aware that in this scenario, the Atlantis workflow requires the use of the -var-file parameter, specifically in the form of -var-file=(workspace_vars/dev.tf|workspace_vars/staging.tf).

If there is the need for additional workflows, they can be easily added at code level. Current code is ready to easily add new workflows, while sharing common actions.
```
func detectProjectWorkspaces(foldersList []ProjectFolder, workflow string, patternDetector string, changedFiles []string) (updatedFoldersList []ProjectFolder, err error) {
	// Detect project workspaces based on the workflow
	switch workflow {
	case "single-workspace":
		updatedFoldersList, err = singleWorkspaceDetectProjectWorkspaces(foldersList)
	case "multi-workspace":
		updatedFoldersList, err = multiWorkspaceDetectProjectWorkspaces(changedFiles, foldersList, patternDetector)
	}
	// You can add more workflows rules here if required
	return updatedFoldersList, err
}
```

*Note that `atlantis-yaml-detector` will scan all files and folders under the current path. Alternatively you can pass the base path of terraform folders using --terraform-base-dir arg.*

-------

**Build**
---------

**Prerequisites**

*Before you proceed, ensure you have the following installed:*
- Go (1.20 or later)
- Make

--------
**Quick Start**

*Clone this repository:*
```
# git clone https://github.com/totmicro/atlantis-yaml-generator.git
# cd atlantis-yaml-generator
```
*Install missing dependencies:*
```
# make install
```

*Run tests:*
```
# make tests
```
*Build the project:*
```
# make build
```
*To build binaries for different platforms:*
```
# make build-all
```
--------
**Makefile Commands**

*Here are the available commands in the Makefile:*

*   `install`: Install missing dependencies.
*   `tests`: Run tests.
*   `coverage`: Run tests and get coverage.
*   `build`: Build the project.
*   `build-all`: Build amd64 and arm64 binaries for Linux and Darwin.
------------
**Contributing**

Feel free to contribute to this project. We welcome your input!

