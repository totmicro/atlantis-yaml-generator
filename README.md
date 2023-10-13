**Atlantis YAML Generator**
--------------------------


`atlantis-yaml-generator` is a command-line tool designed to simplify the generation of Atlantis configuration files (atlantis.yaml).

It automates the process of detecting projects within your Terraform codebase, configuring workspaces, and generating the required configuration for the Atlantis CI/CD tool.

--------
**Features**
------------

- Generate `atlantis.yaml` files with ease.
- Automatically detect projects and workspaces based on your specified workflow.
- Customize configurations such as automerge, parallel plan/apply, and more.
- Flexible project inclusion and exclusion using regular expressions.
- Create an Atlantis project by taking into account the files modified in the PR. This approach
 ensures that only projects relevant to the current pull request are generated, avoiding the creation of any unrelated projects.
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
  atlantis-yaml-generator  --base-repo-owner totmicro --base-repo-name terraform --pull-num 884 -w multi-workspace --pattern-detector workspace_vars --gh-token ghp_xxx
  ```
*Generate an `atlantis.yaml` file for a single workspace workflow:*
  ```
  atlantis-yaml-generator  --base-repo-owner totmicro --base-repo-name datadog-as-code --pull-num 123 -w single-workspace --pattern-detector main.tf --gh-token ghp_xxx
  ```
*Generate an `atlantis.yaml` file for a multiple workspace workflow with filtering:*
  ```
  atlantis-yaml-generator/atlantis-yaml-generator  --base-repo-owner totmicro --base-repo-name terraform --pull-num 884 -w multi-workspace --pattern-detector workspace_vars --gh-token ghp_xxxx --terraform-base-dir $(pwd) --excluded-projects "(databases-onboarding-tasks-service-production)$"
  ```

*Use environment variables to pass the sensitive args*

*When you run this command within an Atlantis workflow, it will make an effort to automatically identify the GitHub token by extracting it from the URL in the .git/config file.*

*When you run this command within an Atlantis workflow, `base-repo` `base-repo-owner` and `pull-num` parameters will be automatically identified.*

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

**Atlantis integration**
-------------

Incorporating this generator into Atlantis requires setting up a [`pre-workflow-run`](https://www.runatlantis.io/docs/pre-workflow-hooks.html) as follows:

```yaml
pre_workflow_hooks:
  - run: >
      atlantis-yaml-generator  -w single-workspace --pattern-detector main.tf
```

Sample output:

<details><summary>Files changed in the PR</summary>

```
project/one/main.tf
project/two/main.tf
```
</details>

<details><summary>Rendered atlantis.yaml file</summary>

```
version: 3
automerge: true
parallel_apply: true
parallel_plan: true
projects:
    - name: project-one
      workspace: default
      workflow: single-workspace
      dir: project/one
      autoplan:
        enabled: true
        when_modified:
            - '**/*.tf'
            - '**/*.tfvars'
            - '**/*.json'
            - '**/*.tpl'
            - '**/*.tmpl'
            - '**/*.xml'
    - name: project-two
      workspace: default
      workflow: single-workspace
      dir: project/two
      autoplan:
        enabled: true
        when_modified:
            - '**/*.tf'
            - '**/*.tfvars'
            - '**/*.json'
            - '**/*.tpl'
            - '**/*.tmpl'
            - '**/*.xml'
```
</details>

Additionally, it's essential to ensure that the atlantis-yaml-generator binary is included on your Atlantis server.

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
**Limitations**
---------
- This project is specifically designed for use with the GitHub SCM.
------------
**Contributing**

Feel free to contribute to this project. We welcome your input!

