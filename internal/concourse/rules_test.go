package concourse

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/suite"
	"gotest.tools/assert"
)

type RulesTestSuite struct {
	suite.Suite

	ps      PipelineState
	actions []BuildAction
}

func (suite *RulesTestSuite) SetupTest() {
	jsonFile, _ := ioutil.ReadFile("testdata/PipelineState-reference.json")
	suite.ps = PipelineState{}
	json.Unmarshal(jsonFile, &suite.ps)

	suite.actions = NextBuildActions(suite.ps)
}

func TestRulesTestSuite(t *testing.T) {
	suite.Run(t, new(RulesTestSuite))
}

func (suite *RulesTestSuite) TestBasicJsonLoad() {
	assert.Equal(suite.T(), "pipeline-a", suite.ps.Name)
}

func buildActionByJobName(bas []BuildAction, jobName string) BuildAction {
	for _, ba := range bas {
		if ba.JobName == jobName {
			return ba
		}
	}

	return BuildAction{}
}

func buildResourceByResourceName(brs []BuildResource, resourceName string) BuildResource {
	for _, br := range brs {
		if br.ResourceName == resourceName {
			return br
		}
	}

	return BuildResource{}
}

func (suite *RulesTestSuite) TestBuildAction() {
	assert.Equal(suite.T(), 3, len(suite.actions))

	// job-1
	baJob1 := buildActionByJobName(suite.actions, "job-1")
	buildResourceA := buildResourceByResourceName(baJob1.ResourceInputs, "resource-a")
	buildResourceB := buildResourceByResourceName(baJob1.ResourceInputs, "resource-b")

	assert.Equal(suite.T(), "job-1", baJob1.JobName)
	assert.Equal(suite.T(), 2, len(baJob1.ResourceInputs))
	assert.Equal(suite.T(), "resource-a", buildResourceA.ResourceName)
	assert.Equal(suite.T(), "001", buildResourceA.VersionSha256)
	assert.Equal(suite.T(), "resource-b", buildResourceB.ResourceName)
	assert.Equal(suite.T(), "003", buildResourceB.VersionSha256)

	// job-2
	baJob2 := buildActionByJobName(suite.actions, "job-2")

	assert.Equal(suite.T(), "job-2", baJob2.JobName)
	assert.Equal(suite.T(), 0, len(baJob2.ResourceInputs))

	// job-3
	baJob3 := buildActionByJobName(suite.actions, "job-3")
	buildResourceA = buildResourceByResourceName(baJob3.ResourceInputs, "resource-a")
	buildResourceB = buildResourceByResourceName(baJob3.ResourceInputs, "resource-b")

	assert.Equal(suite.T(), "job-3", baJob3.JobName)
	assert.Equal(suite.T(), 1, len(baJob3.ResourceInputs))
	assert.Equal(suite.T(), "resource-b", buildResourceB.ResourceName)
	assert.Equal(suite.T(), "001", buildResourceB.VersionSha256)
}

func (suite *RulesTestSuite) TestSuccessfulForAllProfiles() {
	// All successful Job Inputs
	resourceVersion := ResourceVersion{
		JobInputs: []ResourceVersionJobHistory{
			{JobName: "job-1", Builds: []Build{
				{Status: BuildStatusSuccess},
			}},
			{JobName: "job-2", Builds: []Build{
				{Status: BuildStatusSuccess},
			}},
		},
	}
	resourceJobProfiles := []upstreamJobResourceProfile{
		{jobName: "job-1", resourceDirection: ResourceDirectionInput},
		{jobName: "job-2", resourceDirection: ResourceDirectionInput},
	}
	assert.Equal(suite.T(), true, successfulForAllProfiles(resourceVersion, resourceJobProfiles))

	// All Job Inputs, one error status
	resourceVersion = ResourceVersion{
		JobInputs: []ResourceVersionJobHistory{
			{JobName: "job-1", Builds: []Build{
				{Status: BuildStatusError},
			}},
			{JobName: "job-2", Builds: []Build{
				{Status: BuildStatusSuccess},
			}},
		},
	}
	resourceJobProfiles = []upstreamJobResourceProfile{
		{jobName: "job-1", resourceDirection: ResourceDirectionInput},
		{jobName: "job-2", resourceDirection: ResourceDirectionInput},
	}
	assert.Equal(suite.T(), false, successfulForAllProfiles(resourceVersion, resourceJobProfiles))

	// One successful Job Input, one successful Job Output
	resourceVersion = ResourceVersion{
		JobInputs: []ResourceVersionJobHistory{
			{JobName: "job-1", Builds: []Build{
				{Status: BuildStatusSuccess},
			}},
		},
		JobOutputs: []ResourceVersionJobHistory{
			{JobName: "job-2", Builds: []Build{
				{Status: BuildStatusSuccess},
			}},
		},
	}
	resourceJobProfiles = []upstreamJobResourceProfile{
		{jobName: "job-1", resourceDirection: ResourceDirectionInput},
		{jobName: "job-2", resourceDirection: ResourceDirectionOutput},
	}
	assert.Equal(suite.T(), true, successfulForAllProfiles(resourceVersion, resourceJobProfiles))

	// One successful Job Input, one error status Job Output
	resourceVersion = ResourceVersion{
		JobInputs: []ResourceVersionJobHistory{
			{JobName: "job-1", Builds: []Build{
				{Status: BuildStatusSuccess},
			}},
		},
		JobOutputs: []ResourceVersionJobHistory{
			{JobName: "job-2", Builds: []Build{
				{Status: BuildStatusError},
			}},
		},
	}
	resourceJobProfiles = []upstreamJobResourceProfile{
		{jobName: "job-1", resourceDirection: ResourceDirectionInput},
		{jobName: "job-2", resourceDirection: ResourceDirectionOutput},
	}
	assert.Equal(suite.T(), false, successfulForAllProfiles(resourceVersion, resourceJobProfiles))
}
