package concourse

import (
	"sort"
)

// upstreamJobResourceProfile holds information about whether a Resource is used
// as an Input or Output (get or put, respectively) on a Job. This is useful
// when trying to find candidate Resource Versions when one Job has a "passed"
// dependency on one or more upstream Jobs.
type upstreamJobResourceProfile struct {
	jobName           string
	resourceName      string
	resourceDirection string
}

// NextBuildActions determines the collection of BuildActions which could
// next execute based on the given PipelineState
func NextBuildActions(ps PipelineState) []BuildAction {
	buildActions := []BuildAction{}

	for _, job := range ps.Jobs {
		buildAction := BuildAction{
			JobName:        job.Name,
			ResourceInputs: []BuildResource{},
		}
		availableInputsCanTriggerJob := false

		for _, input := range job.Inputs {
			resource := ps.resource(input.Name)
			var latestVersionForJob ResourceVersion

			upstreamJobResourceProfiles := getUpstreamJobResourceProfile(ps, input)
			latestVersionForJob = latestResourceVersionForJob(resource.Versions, upstreamJobResourceProfiles)

			buildResource := BuildResource{
				ResourceName:  input.Name,
				VersionSha256: latestVersionForJob.Sha256,
			}
			buildAction.ResourceInputs = append(buildAction.ResourceInputs, buildResource)

			availableInputsCanTriggerJob = availableInputsCanTriggerJob || resourceVersionCanTriggerJob(latestVersionForJob, job.Name)
		}

		if availableInputsCanTriggerJob {
			buildAction.Action = ActionExecute
		}

		buildActions = append(buildActions, buildAction)
	}

	return buildActions
}

func getUpstreamJobResourceProfile(ps PipelineState, input JobResource) []upstreamJobResourceProfile {
	profiles := []upstreamJobResourceProfile{}

	resource := ps.resource(input.Name)
	if hasUpstreamJobDependencies(input) {
		for _, passedJobName := range input.Passed {
			upstreamJob := ps.job(passedJobName)
			profile := upstreamJobResourceProfile{jobName: upstreamJob.Name, resourceName: resource.Name}
			if upstreamJob.getsFromResource(resource.Name) {
				profile.resourceDirection = ResourceDirectionInput
			}
			if upstreamJob.putsToResource(resource.Name) {
				profile.resourceDirection = ResourceDirectionOutput
			}
			profiles = append(profiles, profile)
		}
	}

	return profiles
}

func hasUpstreamJobDependencies(jobResource JobResource) bool {
	return len(jobResource.Passed) > 0
}

func latestResourceVersionForJob(versions []ResourceVersion, profiles []upstreamJobResourceProfile) ResourceVersion {
	latest := ResourceVersion{}

	sortByDiscoveryDateDescending(versions)
	for _, resourceVersion := range versions {
		if successfulForAllProfiles(resourceVersion, profiles) {
			latest = resourceVersion
			break
		}
	}

	return latest
}

func successfulForAllProfiles(resourceVersion ResourceVersion, resourceJobProfiles []upstreamJobResourceProfile) bool {
	result := true

	for _, jobProfile := range resourceJobProfiles {
		switch jobProfile.resourceDirection {
		case ResourceDirectionInput:
			result = result && buildForJobSuccessful(resourceVersion.JobInputs, jobProfile.jobName)
		case ResourceDirectionOutput:
			result = result && buildForJobSuccessful(resourceVersion.JobOutputs, jobProfile.jobName)
		}
	}

	return result
}

func buildForJobSuccessful(jobHistories []ResourceVersionJobHistory, jobName string) bool {
	for _, jobHistory := range jobHistories {
		if jobHistory.JobName == jobName {
			for _, build := range jobHistory.Builds {
				if build.Status == BuildStatusSuccess {
					return true
				}
			}
		}
	}

	return false
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
