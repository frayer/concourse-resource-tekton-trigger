package concourse

import (
	"testing"
)

//                      +---------+                            +---------+
// +-------------+      |         |       +------------+       |         |
// | resource-x  +------+         +-------+ resource-y +-------+         |
// +------+------+      |         |       +------------+       |         |
//        |             |         |                            |         |
//        |             |  job-a  |                            |  job-b  |
//        |             |         |                            |         |
//        |             |         |               +------------+         |
//        |             |         |               |            |         |
//        |             |         |               |            |         |
//        |             +---------+               |            +---------+
//        |                                       |
//        +---------------------------------------+

func getMockPipelineState() *PipelineState {
	return &PipelineState{
		Name:      "pipeline-a",
		Resources: getMockResources(),
		Jobs: []Job{
			{
				Name: "job-a",
				Inputs: []ResourceStep{
					{Name: "resource-x"},
				},
				Outputs: []ResourceStep{
					{Name: "resource-y"},
				},
				Builds: []Build{
					{
						Id:          "1",
						BuildInputs: []BuildResource{{ResourceName: "resource-x", VersionSha256: "12345"}},
						Status:      "passed",
					},
					{
						Id:          "2",
						BuildInputs: []BuildResource{{ResourceName: "resource-x", VersionSha256: "23456"}},
						Status:      "passed",
					},
				},
			},
			{
				Name: "job-b",
				Inputs: []ResourceStep{
					{Name: "resource-x", Passed: []string{"job-a"}},
					{Name: "resource-y"},
				},
				Builds: []Build{
					{Id: "1", BuildInputs: []BuildResource{
						{ResourceName: "resource-x", VersionSha256: "12345"},
						{ResourceName: "resource-y", VersionSha256: "abcde"},
					}},
				},
			},
		},
	}
}

func getMockResources() []Resource {
	return []Resource{
		{
			Name: "resource-x",
			DiscoveredVersions: []Version{
				{Sha256: "12345", DiscoveryDate: "1"},
				{Sha256: "23456", DiscoveryDate: "2"},
				{Sha256: "34567", DiscoveryDate: "2"},
			},
		},
		{
			Name: "resource-y",
			DiscoveredVersions: []Version{
				{Sha256: "12345", DiscoveryDate: "1"},
				{Sha256: "23456", DiscoveryDate: "2"},
			},
		},
	}
}

func TestMock(t *testing.T) {
	pipeline := getMockPipelineState()
	expected := "pipeline-a"
	received := pipeline.Name
	if received != expected {
		t.Errorf("expected %s, but got %s", expected, received)
	}
}
