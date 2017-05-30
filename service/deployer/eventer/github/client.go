package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	microerror "github.com/giantswarm/microkit/error"
)

const (
	// deploymentUrlFormat is the string format for the GitHub
	// API call to fetch Deployments.
	deploymentUrlFormat = "https://api.github.com/repos/%s/%s/deployments"

	// deploymentStatusUrlFormat is the string format for the
	// GitHub API call to post Deployment Statuses.
	deploymentStatusUrlFormat = "https://api.github.com/repos/%s/%s/deployments/%v/statuses"

	// etagHeader is the header used for etag.
	// See: https://en.wikipedia.org/wiki/HTTP_ETag.
	etagHeader = "Etag"
)

// request makes a request, handling any metrics and logging.
func (e *githubEventer) request(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("token %s", e.oauthToken))

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	// Update rate limit metrics.
	if err := updateRateLimitMetrics(resp); err != nil {
		return nil, microerror.MaskAny(err)
	}

	return resp, err
}

// filterDeployments filters out any deployments we don't want,
// such as deployments for other installations.
func (e *githubEventer) filterDeployments(deployments []deployment) []deployment {
	matches := []deployment{}

	for _, deployment := range deployments {
		if deployment.Environment == e.environment {
			matches = append(matches, deployment)
		}
	}

	return matches
}

// fetchNewDeploymentEvents fetches any new GitHub Deployment Events for the
// given project.
func (e *githubEventer) fetchNewDeploymentEvents(project string, etagMap map[string]string) ([]deployment, error) {
	e.logger.Log("debug", "fetching deployments", "project", project)

	url := fmt.Sprintf(
		deploymentUrlFormat,
		e.organisation,
		project,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	// If we have an etag header for this project, then we have already
	// requested deployment events for it.
	// So, set the header so we only get notified of new events.
	if val, ok := etagMap[project]; ok {
		req.Header.Set("If-None-Match", val)
	}

	startTime := time.Now()

	resp, err := e.request(req)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}
	defer resp.Body.Close()

	updateDeploymentMetrics(e.organisation, project, resp.StatusCode, startTime)

	// Save the new etag header, so we don't get these deployment events again.
	etagMap[project] = resp.Header.Get(etagHeader)

	if resp.StatusCode == http.StatusNotModified {
		e.logger.Log("debug", "no new deployment events, continuing", "project", project)
		return []deployment{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, microerror.MaskAnyf(unexpectedStatusCode, fmt.Sprintf("received non-200 status code: %v", resp.StatusCode))
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	var deployments []deployment
	if err := json.Unmarshal(bytes, &deployments); err != nil {
		return nil, microerror.MaskAny(err)
	}

	deployments = e.filterDeployments(deployments)

	if len(deployments) > 0 {
		e.logger.Log("debug", "found new deployment events", "project", project)
	}

	return deployments, nil
}

// postDeploymentEventStatus posts a Deployment Status for the given Deployment.
func (e *githubEventer) postDeploymentStatus(project string, id int, state deploymentStatusState) error {
	e.logger.Log("debug", "posting deployment status", "project", project, "id", id, "state", state)

	url := fmt.Sprintf(
		deploymentStatusUrlFormat,
		e.organisation,
		project,
		id,
	)

	status := deploymentStatus{
		State: state,
	}

	payload, err := json.Marshal(status)
	if err != nil {
		return microerror.MaskAny(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return microerror.MaskAny(err)
	}

	startTime := time.Now()

	resp, err := e.request(req)
	if err != nil {
		return microerror.MaskAny(err)
	}
	defer resp.Body.Close()

	updateDeploymentStatusMetrics(e.organisation, project, resp.StatusCode, startTime)

	if resp.StatusCode != http.StatusCreated {
		return microerror.MaskAnyf(unexpectedStatusCode, fmt.Sprintf("received non-200 status code: %v", resp.StatusCode))
	}

	return nil
}
