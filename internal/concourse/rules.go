package concourse

import (
	"sort"
)

const (
	actionExecute = "execute"
)

func NextBuildActions(ps PipelineState) []BuildAction {
	buildActions := []BuildAction{}

	resources := nameMappedResources(ps.Resources)

	for _, job := range ps.Jobs {
		buildAction := BuildAction{
			JobName:        job.Name,
			ResourceInputs: []BuildResource{},
		}
		availableInputsCanTrigger := false

		resourceInputNames := gatherInputNames(job)
		for _, resourceInputName := range resourceInputNames {
			resourceVersions := resources[resourceInputName].Versions
			latestResourceVersionForJob := latestResourceVersionForJob(resourceVersions)

			buildResource := BuildResource{
				ResourceName:  resourceInputName,
				VersionSha256: latestResourceVersionForJob.Sha256,
			}
			buildAction.ResourceInputs = append(buildAction.ResourceInputs, buildResource)

			availableInputsCanTrigger = availableInputsCanTrigger || resourceVersionCanTriggerJob(latestResourceVersionForJob, job.Name)
		}

		if availableInputsCanTrigger {
			buildAction.Action = actionExecute
		}

		buildActions = append(buildActions, buildAction)
	}

	return buildActions
}

func gatherInputNames(job JobDefinition) []string {
	jds := []string{}
	for _, input := range job.Inputs {
		jds = append(jds, input.Name)
	}
	return jds
}

func latestResourceVersionForJob(resourceVersions []ResourceVersion) ResourceVersion {
	latestResourceVersion := ResourceVersion{}

	if len(resourceVersions) > 0 {
		sortByDiscoveryDateDescending(resourceVersions)
		latestResourceVersion = resourceVersions[0]
	}

	return latestResourceVersion
}

func sortByDiscoveryDateDescending(resourceVersions []ResourceVersion) {
	sort.Slice(resourceVersions, func(i, j int) bool {
		return resourceVersions[i].DiscoveryDate > resourceVersions[j].DiscoveryDate
	})
}

func resourceVersionCanTriggerJob(resourceVersion ResourceVersion, jobName string) bool {
	existingBuildFoundForJob := false

	for _, jobInputHistory := range resourceVersion.JobInputs {
		if jobInputHistory.JobName == jobName && len(jobInputHistory.Builds) > 0 {
			existingBuildFoundForJob = true
		}
	}

	return !existingBuildFoundForJob
}

func nameMappedResources(resources []PipelineResource) map[string]PipelineResource {
	mappedPipelineResources := make(map[string]PipelineResource)
	for _, resource := range resources {
		mappedPipelineResources[resource.Name] = resource
	}
	return mappedPipelineResources
}
