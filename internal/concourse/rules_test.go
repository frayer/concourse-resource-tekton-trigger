package concourse_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	c "github.com/frayer/concourse-resource-tekton-trigger/internal/concourse"
	"github.com/stretchr/testify/suite"
	"gotest.tools/assert"
)

type RulesTestSuite struct {
	suite.Suite

	ps      c.PipelineState
	actions []c.BuildAction
}

func (suite *RulesTestSuite) SetupTest() {
	jsonFile, _ := ioutil.ReadFile("testdata/PipelineState-reference.json")
	suite.ps = c.PipelineState{}
	json.Unmarshal(jsonFile, &suite.ps)
	suite.actions = c.NextBuildActions(suite.ps)
}

func TestRulesTestSuite(t *testing.T) {
	suite.Run(t, new(RulesTestSuite))
}

func (suite *RulesTestSuite) TestBasicJsonLoad() {
	assert.Equal(suite.T(), "pipeline-a", suite.ps.Name)
}

func buildActionByJobName(bas []c.BuildAction, jobName string) c.BuildAction {
	for _, ba := range bas {
		if ba.JobName == jobName {
			return ba
		}
	}

	return c.BuildAction{}
}

func buildResourceByResourceName(brs []c.BuildResource, resourceName string) c.BuildResource {
	for _, br := range brs {
		if br.ResourceName == resourceName {
			return br
		}
	}

	return c.BuildResource{}
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
	assert.Equal(suite.T(), "003", buildResourceB.VersionSha256)
}
