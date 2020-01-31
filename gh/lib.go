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
    - ["T-perf",    "000000", "Anything concerning perf"]
    - ["T-feature", "000000", "Anything that is adds new behavior"]
    - ["T-doc",     "000000", "Anything that adds documentation"]
    - ["T-hotfix",  "000000", "Critical things to deploy"]
    - ["L-ready",   "000000", "Anything ready for merging and releasing"]


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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/psyomn/psy/common"

	"github.com/go-yaml/yaml"
)

func Run(args common.RunParams) common.RunReturn {
	if len(args) != 4 {
		fmt.Println("usage: ")
		// nice to have
		// fmt.Println("  list-labels <owner> <repo>")
		fmt.Println("  test-labels <config.yaml> <owner> <repo>")
		fmt.Println("  poison-labels <config.yaml> <owner> <repo>")
		return errors.New("wrong usage")
	}

	ghToken, err := getGithubToken()
	if err != nil {
		return err
	}

	config, err := parseLabelConfig(args[1])
	if err != nil {
		return err
	}

	owner := args[2]
	repoName := args[3]
	labels, err := getLabels(owner, repoName)
	if err != nil {
		return err
	}

	repo := repo{
		owner:   owner,
		repo:    repoName,
		labels:  labels,
		token:   ghToken,
		actions: config,
	}

	repo.poison(false)

	fmt.Println(repo)

	return nil
}

func getGithubToken() (string, error) {
	// Maybe sourced as env variable for now, maybe something else
	// in the future

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	} else {
		return "", errors.New("no github token set")
	}
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

func getLabels(owner, repo string) ([]label, error) {
	const url = "https://api.github.com/"

	resp, err := http.Get(url + "repos/" + owner + "/" + repo + "/labels")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var labels []label
	err = json.Unmarshal(bytes, &labels)
	if err != nil {
		return nil, err
	}

	return labels, nil
}
