package concourse

type PipelineState struct {
	Name      string     `json:"name"`
	Resources []Resource `json:"resources"`
	Jobs      []Job      `json:"jobs"`
}

type Job struct {
	Name    string         `json:"name"`
	Inputs  []ResourceStep `json:"inputs,omitempty"`
	Outputs []ResourceStep `json:"outputs,omitempty"`
	Builds  []Build        `json:"builds,omitempty"`
}

type ResourceStep struct {
	Name    string   `json:"resource"`
	Passed  []string `json:"passed,omitempty"`
	Trigger bool     `json:"trigger,omitempty"`
}

type Build struct {
	Id           string          `json:"id"`
	BuildInputs  []BuildResource `json:"buildInputs"`
	BuildOutputs []BuildResource `json:"buildOutputs"`
	Status       string          `json:"string"`
}

type BuildResource struct {
	ResourceName  string `json:"resourceName"`
	VersionSha256 string `json:"versionSha256"`
}

type Resource struct {
	Name               string    `json:"name"`
	DiscoveredVersions []Version `json:"discoveredVersions"`
}

type Version struct {
	Sha256        string `json:"sha256"`
	DiscoveryDate string `json:"discoveryDate"`
}

type BuildAction struct {
	Action         string          `json:"action"`
	JobName        string          `json:"jobName"`
	ResourceInputs []BuildResource `json:"resourceInputs"`
}
