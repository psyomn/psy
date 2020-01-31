package gh

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
	NodeId      string `json:"node_id"`
	Url         string `json:"url"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
	Description string `json:"description"`
}

type labelActions struct {
	Rename map[string]interface{}
	Create [][]string
}

type repo struct {
	owner  string
	repo   string
	labels []label

	token   string
	actions *labelActions
}

func (s *repo) poison(dryrun bool) {
	s.rename()
	s.create()
}

func (s *repo) rename() {
}

func (s *repo) create() {
}
