/*
Package gh is a simple package to interface with github

It's not supposed to do everything. My specific use is to set github
labels, according to a configuration that I think is OK.

usage:
  psy gh test-labels <config.yaml> <owner> <repo>
    will do a dry run of label nomenclature application

  psy gh poison <config.yaml> <owner> <repo>
    will edit the labels to respect the given specs.


The config file should look something like this:

  ---
  # Github label configuration.

  # T = Type
  # L = Lifecycle

  rename:
    bug:         ["T-bug",       "a00000", "A software fault or failure"]
    duplicate:   ["L-duplicate", "16336d", "dOops!"]
    enhancement: ["T-upkeep",    "009a00", "Anything that rejuvenates"]
    invalid:     ["L-invalid",   "ffbb0e", "Anything that is invalid"]
    question:    ["T-question",  "ffaa99", "Questions are good!"]
    wontfix:     ["L-wontfix",   "231f20", "Things not to fix"]

  create:
    - ["T-perf",    "f51919", "Anything concerning perf"]
    - ["T-feature", "078a00", "Anything that is adds new behavior"]
    - ["T-doc",     "001c8a", "Anything that adds documentation"]
    - ["T-hotfix",  "ff5100", "Critical things to deploy"]
    - ["L-ready",   "ffe100", "Anything ready for merging and releasing"]


Copyright 2020 Simon Symeonidis (psyomn)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package gh

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/psyomn/psy/common"

	"github.com/go-yaml/yaml"
)

func usage() error {
	fmt.Println("usage: ")
	fmt.Println("  help - print this")
	fmt.Println("  generate-config - print configuration to stdout")
	fmt.Println("  list-labels <owner> <repo> - list labels on a github repository")
	fmt.Println("  dryrun-poison <config.yml> <owner> <repo> - dryrun changes to labels")
	fmt.Println("  poison-labels <config.yml> <owner> <repo> - run actual changes to labels")
	return errors.New("wrong usage")
}

// Run will run the subcommand with provided parameters
func Run(args common.RunParams) common.RunReturn {
	ghToken, err := getGithubToken()

	if err != nil {
		return err
	}

	if len(args) == 1 {
		if args[0] == "generate-config" {
			return generateConfig()
		} else if args[0] == "help" {
			usage()
			return nil
		} else {
			return usage()
		}
	} else if len(args) == 3 {
		if args[0] != "list-labels" {
			return usage()
		}
		owner := args[1]
		repoName := args[2]
		return listLabels(owner, repoName, ghToken)
	} else if len(args) == 4 {
		dryRun := false
		if args[0] == "dryrun-poison" {
			dryRun = true
		} else if args[0] == "poison-labels" {
		} else {
			return usage()
		}

		config := args[1]
		owner := args[2]
		repoName := args[3]

		return poison(dryRun, config, owner, repoName, ghToken)
	}

	return usage()
}

func getGithubToken() (string, error) {
	// Maybe sourced as env variable for now, maybe something else
	// in the future

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	}

	return "", errors.New("no github token set in environment variable GITHUB_TOKEN")
}

func parseLabelConfig(path string) (*labelActions, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadAll(fh)

	var actions labelActions
	err = yaml.Unmarshal(contents, &actions)

	return &actions, err
}

func generateConfig() error {
	fmt.Println(`---
# Github label configuration.

# T = Type
# L = Lifecycle

rename:
  bug:         ["T-bug",       "a00000", "A software fault or failure"]
  duplicate:   ["L-duplicate", "16336d", "dOops!"]
  enhancement: ["T-upkeep",    "009a00", "Anything that rejuvenates"]
  invalid:     ["L-invalid",   "ffbb0e", "Anything that is invalid"]
  question:    ["T-question",  "ffaa99", "Questions are good!"]
  wontfix:     ["L-wontfix",   "231f20", "Things not to fix"]

create:
  - ["T-perf",    "f51919", "Anything concerning perf"]
  - ["T-feature", "078a00", "Anything that is adds new behavior"]
  - ["T-doc",     "001c8a", "Anything that adds documentation"]
  - ["T-hotfix",  "ff5100", "Critical things to deploy"]
  - ["L-ready",   "ffe100", "Anything ready for merging and releasing"]`)

	return nil
}

func listLabels(owner, repoName, token string) error {
	repo := repo{
		owner: owner,
		repo:  repoName,
		token: token,
	}

	repo.getLabels()
	fmt.Println(repo.String())
	return nil
}

func poison(dryRun bool, configFileName, owner, repoName, ghToken string) error {
	config, err := parseLabelConfig(configFileName)
	if err != nil {
		return err
	}

	repo := repo{
		owner:   owner,
		repo:    repoName,
		token:   ghToken,
		actions: config,
	}

	if err := repo.dryRun(dryRun).poison(); err != nil {
		return err
	}

	fmt.Println(repo.String())

	return nil
}
