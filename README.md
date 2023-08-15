
Atlantis YAML Generator
=======================


`atlantis-yaml-generator` is a command-line tool designed to simplify the generation of Atlantis configuration files (atlantis.yaml). It automates the process of detecting projects within your Terraform codebase, configuring workspaces, and generating the required configuration for the Atlantis CI/CD tool.

Features
--------

*   Generate `atlantis.yaml` files with ease.
*   Automatically detect projects and workspaces based on your specified workflow.
*   Customize configurations such as automerge, parallel plan/apply, and more.
*   Flexible project inclusion and exclusion using regular expressions.

Installation
------------

To install `atlantis-yaml-generator`, you need to have Go (Golang) installed. Use the following command to install the tool:

    go get github.com/yourusername/atlantis-yaml-generator

Usage
-----

Run the tool using the following command:

    atlantis-yaml-generator [flags]

Available flags:

*   `--atlantis-automerge`: Atlantis automerge config value.
*   `--atlantis-parallel-apply`: Atlantis parallel apply config value.
*   `--atlantis-parallel-plan`: Atlantis parallel plan config value.
*   `--output-file`: Atlantis output file name.
*   `--atlantis-workflow`: Atlantis Workflow to be used. Options: single-workspace or multi-workspace.
*   ... and more flags for specifying various parameters.

For detailed usage instructions and available flags, run:

    atlantis-yaml-generator --help

Examples
--------

Generate an `atlantis.yaml` file for a single workspace workflow:

    atlantis-yaml-generator --atlantis-workflow single-workspace --output-file atlantis.yaml

Generate an `atlantis.yaml` file for a multi-workspace workflow:

    atlantis-yaml-generator --atlantis-workflow multi-workspace --output-file atlantis.yaml

Contributing
------------

Contributions are welcome! Feel free to open issues and submit pull requests on the GitHub repository.

License
-------

This project is licensed under the MIT License.
