package concourse

import (
	"sort"
)

func NextBuildActions(ps *PipelineState) []BuildAction {
	buildActions := []BuildAction{}

	for _, job := range ps.Jobs {
		buildAction := BuildAction{
			Action:         "execute",
			JobName:        job.Name,
			ResourceInputs: []BuildResource{},
		}

		for _, input := range job.Inputs {
			resourceName := input.Name
			latestResourceVersion := latestAvailableResourceVersion(ps.Resources, resourceName)
			if !resourceVersionAlreadyProcessed(job, resourceName, latestResourceVersion.Sha256) {
				buildAction.ResourceInputs = append(buildAction.ResourceInputs, BuildResource{
					ResourceName:  resourceName,
					VersionSha256: latestResourceVersion.Sha256,
				})
			}
		}

		if len(buildAction.ResourceInputs) > 0 {
			buildActions = append(buildActions, buildAction)
		}
	}

	return buildActions
}

func mapByName(items []Job) map[string]Job {
	nameKeyedMap := make(map[string]Job)

	for _, item := range items {
		nameKeyedMap[item.Name] = item
	}

	return nameKeyedMap
}

func latestAvailableResourceVersion(resources []Resource, resourceName string) Version {
	resourcePos := find(len(resources), func(i int) bool {
		return resources[i].Name == resourceName
	})

	resource := resources[resourcePos]

	sort.Slice(resource.DiscoveredVersions, func(i, j int) bool {
		return resource.DiscoveredVersions[i].DiscoveryDate < resource.DiscoveredVersions[j].DiscoveryDate
	})

	return resource.DiscoveredVersions[len(resource.DiscoveredVersions)-1]
}

func resourceVersionAlreadyProcessed(job Job, resourceName string, version string) bool {
	for _, build := range job.Builds {
		for _, input := range build.BuildInputs {
			if input.VersionSha256 == version && input.ResourceName == resourceName {
				return true
			}
		}
	}

	return false
}

func find(n int, predicate func(int) bool) int {
	var i int
	for i = 0; i < n; i++ {
		if predicate(i) {
			return i
		}
	}
	return i + 1
}
