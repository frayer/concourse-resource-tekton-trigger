package concourse

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gotest.tools/assert"
)

type TypesTestSuite struct {
	suite.Suite
}

func (suite *TypesTestSuite) SetupTest() {
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}

func (suite *TypesTestSuite) TestGetResourceByName() {
	ps := PipelineState{
		Resources: []PipelineResource{
			{Name: "resource-a"},
			{Name: "resource-b"},
			{Name: "resource-c"},
		},
	}

	assert.Equal(suite.T(), "resource-b", ps.resource("resource-b").Name)
}

func (suite *TypesTestSuite) TestGetJobByName() {
	ps := PipelineState{
		Jobs: []JobDefinition{
			{Name: "job-a"},
			{Name: "job-b"},
			{Name: "job-c"},
		},
	}

	assert.Equal(suite.T(), "job-b", ps.job("job-b").Name)
}

func (suite *TypesTestSuite) TestJobInputOutputHelpers() {
	jd := JobDefinition{
		Name:    "job-a",
		Inputs:  []JobResource{{Name: "resource-a"}, {Name: "resource-b"}, {Name: "resource-c"}},
		Outputs: []JobResource{{Name: "resource-b"}, {Name: "resource-d"}, {Name: "resource-e"}},
	}

	assert.Equal(suite.T(), true, jd.getsFromResource("resource-a"))
	assert.Equal(suite.T(), true, jd.getsFromResource("resource-b"))
	assert.Equal(suite.T(), true, jd.getsFromResource("resource-c"))
	assert.Equal(suite.T(), false, jd.getsFromResource("resource-d"))
	assert.Equal(suite.T(), false, jd.getsFromResource("resource-e"))

	assert.Equal(suite.T(), false, jd.putsToResource("resource-a"))
	assert.Equal(suite.T(), true, jd.putsToResource("resource-b"))
	assert.Equal(suite.T(), false, jd.putsToResource("resource-c"))
	assert.Equal(suite.T(), true, jd.putsToResource("resource-d"))
	assert.Equal(suite.T(), true, jd.putsToResource("resource-e"))
}
