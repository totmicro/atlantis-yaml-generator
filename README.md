**Atlantis YAML Generator**
--------------------------


`atlantis-yaml-generator` is a command-line tool designed to simplify the generation of Atlantis configuration files (atlantis.yaml).

It automates the process of detecting projects within your Terraform codebase, configuring workspaces, and generating the required configuration for the Atlantis CI/CD tool.

--------
**Features**
------------

- Generate `atlantis.yaml` files with ease.
- Automatically detect projects and workspaces based on discovery mode rules.
- Customize configurations such as automerge, parallel plan/apply, and more.
- Flexible project inclusion and exclusion using regular expressions.
- Create an Atlantis project by taking into account the files modified in the PR. This approach
 ensures that only projects relevant to the current pull request are generated, avoiding the creation of any unrelated projects.
-----

**Usage**
---------

Run the tool using the following command:
`atlantis-yaml-generator` [flags]


| Flag                          | Description                                                    | Equivalent envVar   | Default Value |
| ------------------------------ | -------------------------------------------------------------- | ------------------- | ------------- |
| `--automerge`          | Atlantis automerge config value.                               | `AUTOMERGE`         | `true`          |
| `-r, --base-repo-name` | Github Repo Name.                                              | `BASE_REPO_NAME`    |               |
| `-o, --base-repo-owner`| Github Repo Owner Name.                                        | `BASE_REPO_OWNER`   |               |
| `-x, --excluded-projects`| Atlantis regex filter to exclude projects.                    | `EXCLUDED_PROJECTS` |               |
| `-d, --discovery-mode`| mode used to discover projects                                  | `DISCOVERY_MODE` | `single-workspace`|
| `-t, --gh-token`       | Github Token Value.                                            | `GH_TOKEN`          |               |
| `-h, --help`                  | Help for atlantis-yaml-generator.                               |                     |               |
| `-z, --included-projects`| Atlantis regex filter to only include projects.              | `INCLUDED_PROJECTS` |               |
| `-f, --output-file`    | Atlantis output file name.                                     | `OUTPUT_FILE`       | `atlantis.yaml`          |
| `-e, --output-type`    | Atlantis YAML output type [file stdout]                      | `OUTPUT_TYPE`       | `file`          |
| `--parallel-apply`     | Atlantis parallel apply config value.                         | `PARALLEL_APPLY`    | `true`          |
| `--parallel-plan`      | Atlantis parallel plan config value.                          | `PARALLEL_PLAN`    | `true`          |
| `-q, --pattern-detector`| Discover projects based on files, directories names or regex.      | `PATTERN_DETECTOR`  | `main.tf`      |
| `-u, --pr-filter`      | Filter projects based on the PR changes (Only for github SCM).| `PR_FILTER`       | `false`          |
| `-p, --pull-num`       | Github Pull Request Number to check diffs.                    | `PULL_NUM`          |               |
| `--terraform-base-dir` | Basedir for terraform resources.                               | `TERRAFORM_BASE_DIR`| `./`            |
| `-v, --version`               | Version for atlantis-yaml-generator.                           |                     |               |
| `-m, --when-modified`  | Atlantis When modified (list of strings) to run autoplan.    | `WHEN_MODIFIED`     | `**/*.tf,**/*.tfvars,**/*.json,**/*.tpl,**/*.tmpl,**/*.xml` |
| `-w, --workflow`       | Atlantis Workflow to be used.                                | `WORKFLOW`     |          |


*Note that default values and required flags are defined in [`config.go`](pkg/config/config.go)*

--------

**Examples**
------------

<details><summary>Multi workspace workflow</summary>

```
# atlantis-yaml-generator -d multi-workspace -w myWorkflow --pattern-detector workspace_vars -e stdout

version: 3
automerge: true
parallel_apply: true
parallel_plan: true
projects:
    - name: project_one-dev
      workspace: dev
      workflow: myWorkflow
      dir: project_one
      autoplan:
        enabled: true
        when_modified:
            - '**/*.tf'
            - '**/*.tfvars'
            - '**/*.json'
            - '**/*.tpl'
            - '**/*.tmpl'
            - '**/*.xml'
    - name: project_one-production
      workspace: production
      workflow: myWorkflow
      dir: project_one
      autoplan:
        enabled: true
        when_modified:
            - '**/*.tf'
            - '**/*.tfvars'
            - '**/*.json'
            - '**/*.tpl'
            - '**/*.tmpl'
            - '**/*.xml'
    - name: project_one-staging
      workspace: staging
      workflow: myWorkflow
      dir: project_one
      autoplan:
        enabled: true
        when_modified:
            - '**/*.tf'
            - '**/*.tfvars'
            - '**/*.json'
            - '**/*.tpl'
            - '**/*.tmpl'
            - '**/*.xml'
    - name: project_two-production
      workspace: production
      workflow: myWorkflow
      dir: project_two
      autoplan:
        enabled: true
        when_modified:
            - '**/*.tf'
            - '**/*.tfvars'
            - '**/*.json'
            - '**/*.tpl'
            - '**/*.tmpl'
            - '**/*.xml'
    - name: project_two-staging
      workspace: staging
      workflow: myWorkflow
      dir: project_two
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

<details><summary>Single workspace discovery mode</summary>

```
# atlantis-yaml-generator -d single-workspace -w myWorkflow -e stdout --pattern-detector main.tf

version: 3
automerge: true
parallel_apply: true
parallel_plan: true
projects:
    - name: project_one
      workspace: default
      workflow: myWorkflow
      dir: project_one
      autoplan:
        enabled: true
        when_modified:
            - '**/*.tf'
            - '**/*.tfvars'
            - '**/*.json'
            - '**/*.tpl'
            - '**/*.tmp
```
</details>

<details><summary>Multiple workspace workflow with regex project filtering</summary>

```
# atlantis-yaml-generator -w multi-workspace -e stdout --pattern-detector workspace_vars --included-projects "(^project_two-staging|production)$"

version: 3
automerge: true
parallel_apply: true
parallel_plan: true
projects:
    - name: project_one-production
      workspace: production
      workflow: multi-workspace
      dir: project_one
      autoplan:
        enabled: true
        when_modified:
            - '**/*.tf'
            - '**/*.tfvars'
            - '**/*.json'
            - '**/*.tpl'
            - '**/*.tmpl'
            - '**/*.xml'
    - name: project_two-production
      workspace: production
      workflow: multi-workspace
      dir: project_two
      autoplan:
        enabled: true
        when_modified:
            - '**/*.tf'
            - '**/*.tfvars'
            - '**/*.json'
            - '**/*.tpl'
            - '**/*.tmpl'
            - '**/*.xml'
    - name: project_two-staging
      workspace: staging
      workflow: multi-workspace
      dir: project_two
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

<details><summary>Multi workspace workflow with PR filter</summary>

```
# atlantis-yaml-generator -w multi-workspace --pattern-detector workspace_vars -e stdout --pr-filter true --pull-num 1 --base-repo-name atlantis-yaml-generator --base-repo-owner totmicro --gh-token ghp_xxxx

version: 3
automerge: true
parallel_apply: true
parallel_plan: true
projects:
    - name: examples-multi-workspace/multiple-projects/project_one-dev
      workspace: dev
      workflow: multi-workspace
      dir: examples/multi-workspace/multiple-projects/project_one
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
[PR used in the above example](https://github.com/totmicro/atlantis-yaml-generator/pull/1/files)

</details>

-------
*Use environment variables to pass the sensitive args*

*When you run this command within an Atlantis workflow, it will make an effort to automatically identify the GitHub token by extracting it from the URL in the .git/config file.*

*When you run this command within an Atlantis workflow, `base-repo` `base-repo-owner` and `pull-num` parameters will be automatically identified. (Only github SCM)*

-------

**Discovery Modes**
-------------

Currenlty `atlantis-yaml-generator` support 2 discovery modes:
- `single-workspace`: Intended for Terraform configurations that do not utilize multiple workspaces. In this context, the pattern detector parameters establish the criteria for identifying project folders. For instance, if the pattern-detector is set to main.tf (file), and this file is located at database/dev/main.tf, the resulting project would be labeled as database-env. Consequently, Terraform commands would be executed within the database/dev folder.
- `multiple-workspace`: Intended for Terraform configurations that utilize multiple workspaces. In this context, the pattern detector parameters establish the criteria for identifying project folders. For instance, if the pattern-detector is set to workspace_vars (folder), and there are several files located in this folder, i.e. (database/workspace_vars/dev.tf |database/workspace_vars/staging.tf ), the resulting projects would be labeled as (database-dev|database-staging). Consequently, Terraform commands would be executed within the database folder. Please be aware that in this scenario, the Atlantis workflow requires the use of the -var-file parameter, specifically in the form of -var-file=(workspace_vars/dev.tf|workspace_vars/staging.tf).

If there is the need for additional discovery modes, they can be easily added at code level. Current code is ready to easily add new discovery modes, while sharing common actions.
```
func detectProjectWorkspaces(foldersList []ProjectFolder, discoveryMode string, patternDetector string, changedFiles []string) (updatedFoldersList []ProjectFolder, err error) {
	// Detect project workspaces based on the discovery mode
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
      atlantis-yaml-generator  -d single-workspace -w myWorkflow --pattern-detector main.tf
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
      workflow: myWorkflow
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
      workflow: myWorkflow
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
- `pr-filter` parameter is specifically designed for use with the GitHub SCM.
------------
**Contributing**

Feel free to contribute to this project. We welcome your input!

