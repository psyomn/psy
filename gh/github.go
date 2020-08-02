package gh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Sample github response
// [{
//      node_id: "MDU6TGFiZWw2ODUwOTU0Ng==",
//      url: "https://api.github.com/repos/psyomn/notes/labels/bug",
//      name: "bug",
//      color: "fc2929",
//      default: true,
// 	description: null
// }]

type label struct {
	NodeID      string `json:"node_id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
	Description string `json:"description"`
}

func (s *label) String() string {
	return fmt.Sprintf("[#%s %s]: %s (default: %v)",
		s.Color,
		s.Name,
		s.Description,
		s.Default)
}

type labelActions struct {
	Rename map[string]interface{}
	Create [][]string
}

type repo struct {
	owner  string
	repo   string
	labels []label

	token    string
	actions  *labelActions
	isDryRun bool
}

func (s *repo) create() {
}

func (s *repo) getLabels() error {
	const url = "https://api.github.com/"

	resp, err := http.Get(url + "repos/" + s.owner + "/" + s.repo + "/labels")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &s.labels)
	if err != nil {
		return err
	}

	return nil
}

func (s *repo) dryRun(run bool) *repo {
	s.isDryRun = run
	return s
}

func (s *repo) poison() error {
	err := s.getLabels()
	if err != nil {
		return err
	}

	s.rename()
	s.create()

	return nil
}

func (s *repo) rename() {
}

func (s *repo) String() string {
	var build strings.Builder
	build.Grow(256)

	fmt.Fprintf(&build, "owner  : %s\n", s.owner)
	fmt.Fprintf(&build, "repo   : %s\n", s.repo)
	fmt.Fprintf(&build, "dryrun : %v\n", s.isDryRun)

	fmt.Fprintf(&build, "labels :\n")
	for _, label := range s.labels {
		fmt.Fprintf(&build, "  - %v\n", label.String())
	}

	return build.String()
}
