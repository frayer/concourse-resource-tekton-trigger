package concourse

const (
	ActionExecute = "execute"

	BuildStatusSuccess = "success"
	BuildStatusError   = "error"

	ResourceDirectionInput  = "input"
	ResourceDirectionOutput = "output"
)

type PipelineState struct {
	Name      string             `json:"name"`
	Resources []PipelineResource `json:"resources"`
	Jobs      []JobDefinition    `json:"jobs"`
}

func (ps *PipelineState) resource(name string) *PipelineResource {
	for _, pr := range ps.Resources {
		if pr.Name == name {
			return &pr
		}
	}

	return &PipelineResource{}
}

func (ps *PipelineState) job(name string) *JobDefinition {
	for _, j := range ps.Jobs {
		if j.Name == name {
			return &j
		}
	}

	return &JobDefinition{}
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

func (jd *JobDefinition) getsFromResource(name string) bool {
	for _, input := range jd.Inputs {
		if input.Name == name {
			return true
		}
	}

	return false
}

func (jd *JobDefinition) putsToResource(name string) bool {
	for _, input := range jd.Outputs {
		if input.Name == name {
			return true
		}
	}

	return false
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
