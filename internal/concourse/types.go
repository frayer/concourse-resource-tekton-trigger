package concourse

type PipelineState struct {
	Name      string             `json:"name"`
	Resources []PipelineResource `json:"resources"`
	Jobs      []JobDefinition    `json:"jobs"`
}

type PipelineResource struct {
	Name     string            `json:"name"`
	Versions []ResourceVersion `json:"versions"`
}

type ResourceVersion struct {
	Sha256        string                      `json:"sha256"`
	DiscoveryDate string                      `json:"discoveryDate"`
	JobInputs     []ResourceVersionJobHistory `json:"jobInputs"`
	JobOutputs    []ResourceVersionJobHistory `json:"jobOutputs"`
}

type ResourceVersionJobHistory struct {
	JobName string  `json:"jobName"`
	Builds  []Build `json:"builds"`
}

type Build struct {
	Id     uint   `json:"id"`
	Status string `json:"status"`
}

type JobDefinition struct {
	Name    string        `json:"name"`
	Inputs  []JobResource `json:"inputs"`
	Outputs []JobResource `json:"outputs"`
}

type JobResource struct {
	Name    string   `json:"name"`
	Trigger bool     `json:"trigger"`
	Passed  []string `json:"passed"`
}

type BuildResource struct {
	ResourceName  string `json:"resourceName"`
	VersionSha256 string `json:"versionSha256"`
}

type BuildAction struct {
	Action         string          `json:"action"`
	JobName        string          `json:"jobName"`
	ResourceInputs []BuildResource `json:"resourceInputs"`
}
