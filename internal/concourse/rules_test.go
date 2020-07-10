package concourse

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gotest.tools/assert"
)

type RulesTestSuite struct {
	suite.Suite
	ps      *PipelineState
	actions []BuildAction
}

func (suite *RulesTestSuite) SetupTest() {
	suite.ps = getMockPipelineState()
	suite.actions = NextBuildActions(suite.ps)
}

func TestRulesTestSuite(t *testing.T) {
	suite.Run(t, new(RulesTestSuite))
}

func (suite *RulesTestSuite) TestSingleActionReturned() {
	assert.Equal(suite.T(), 2, len(suite.actions))
}

func (suite *RulesTestSuite) TestActionTypeIsExecute() {
	assert.Equal(suite.T(), "execute", suite.actions[0].Action)
}

func (suite *RulesTestSuite) TestBuildAction() {
	assert.Equal(suite.T(), "job-a", suite.actions[0].JobName)
	assert.Equal(suite.T(), 1, len(suite.actions[0].ResourceInputs))
	assert.Equal(suite.T(), "resource-x", suite.actions[0].ResourceInputs[0].ResourceName)
	assert.Equal(suite.T(), "34567", suite.actions[0].ResourceInputs[0].VersionSha256)

	assert.Equal(suite.T(), "job-b", suite.actions[1].JobName)
	assert.Equal(suite.T(), 2, len(suite.actions[1].ResourceInputs))
	assert.Equal(suite.T(), "resource-x", suite.actions[1].ResourceInputs[0].ResourceName)
	assert.Equal(suite.T(), "34567", suite.actions[1].ResourceInputs[0].VersionSha256)
	assert.Equal(suite.T(), "resource-y", suite.actions[1].ResourceInputs[1].ResourceName)
	assert.Equal(suite.T(), "lmnop", suite.actions[1].ResourceInputs[1].VersionSha256)
}
