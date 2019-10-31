package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/giantswarm/microerror"
)

const (
	// deploymentUrlFormat is the string format for the GitHub
	// API call for Deployments.
	// See: https://developer.github.com/v3/repos/deployments/#list-deployments
	deploymentUrlFormat = "https://api.github.com/repos/%s/%s/deployments"

	// deploymentStatusUrlFormat is the string format for the
	// GitHub API call for Deployment Statuses.
	// See: https://developer.github.com/v3/repos/deployments/#create-a-deployment-status
	deploymentStatusUrlFormat = "https://api.github.com/repos/%s/%s/deployments/%v/statuses"

	// etagHeader is the header used for etag.
	// See: https://en.wikipedia.org/wiki/HTTP_ETag.
	etagHeader = "Etag"

	// rateLimitLimitHeader is the header set by GitHub to show the total
	// rate limit value.
	// See: https://developer.github.com/v3/#rate-limiting
	rateLimitLimitHeader = "X-RateLimit-Limit"

	// rateLimitRemainingHeader is the header set by GitHub to show the
	// remaining rate limit value.
	// See: https://developer.github.com/v3/#rate-limiting
	rateLimitRemainingHeader = "X-RateLimit-Remaining"

	// rateLimitResetHeader is the header set by GitHub to show the time at
	// wich the current rate limit window resets in UTC epoch seconds.
	// See: https://developer.github.com/v3/#rate-limiting
	rateLimitResetHeader = "X-RateLimit-Reset"

	// rateLimitExtraWait is the additional duration that should be waited
	// before rate limit bucket is expected to get refilled in GitHub API.
	rateLimitExtraWait = 5 * time.Second

	// rateLimitAlmostHitThreshold is the limit that makes our rate limiter to
	// wait for extra time until spending last tokens from GitHub token bucket.
	rateLimitAlmostHitThreshold = 5
)

// request makes a request, handling any metrics and logging.
func (e *GithubEventer) request(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("token %s", e.oauthToken))

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Update rate limit metrics.
	if err := updateRateLimitMetrics(resp); err != nil {
		return nil, microerror.Mask(err)
	}

	return resp, err
}

// filterDeploymentsByEnvironment filters out deployments that do not apply
// to this environment.
func (e *GithubEventer) filterDeploymentsByEnvironment(deployments []deployment) []deployment {
	matches := []deployment{}

	for _, deployment := range deployments {
		if deployment.Environment == e.environment {
			matches = append(matches, deployment)
		}
	}

	return matches
}

// filterDeploymentsByStatus filters out deployments that are finished -
// that is, there exists at least one status that is not pending.
func (e *GithubEventer) filterDeploymentsByStatus(deployments []deployment) []deployment {
	matches := []deployment{}

	for _, deployment := range deployments {
		// If there are any statuses apart from pending, we consider the
		// deployment finished, and do not act on it.
		isPending := true
		for _, status := range deployment.Statuses {
			if status.State != pendingState {
				isPending = false
				break
			}
		}

		if isPending {
			matches = append(matches, deployment)
		}
	}

	return matches
}

// fetchNewDeploymentEvents fetches any new GitHub Deployment Events for the
// given project.
func (e *GithubEventer) fetchNewDeploymentEvents(project string, etagMap map[string]string) ([]deployment, error) {
	url := fmt.Sprintf(
		deploymentUrlFormat,
		e.organisation,
		project,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// If we have an etag header for this project, then we have already
	// requested deployment events for it.
	// So, set the header so we only get notified of new events.
	if val, ok := etagMap[project]; ok {
		req.Header.Set("If-None-Match", val)
	}

	e.rateLimiter.Wait()

	startTime := time.Now()

	resp, err := e.request(req)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	defer resp.Body.Close()

	updateDeploymentMetrics(e.organisation, project, resp.StatusCode, startTime)

	err = e.updateRateLimiter(resp)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Save the new etag header, so we don't get these deployment events again.
	etagMap[project] = resp.Header.Get(etagHeader)

	if resp.StatusCode == http.StatusNotModified {
		return []deployment{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		e.logger.Log("error", "Error fetching deployment events", "project", project)
		return nil, microerror.Maskf(unexpectedStatusCode, fmt.Sprintf("received non-200 status code: %v", resp.StatusCode))
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var deployments []deployment
	if err := json.Unmarshal(bytes, &deployments); err != nil {
		return nil, microerror.Mask(err)
	}

	deployments = e.filterDeploymentsByEnvironment(deployments)

	for index, deployment := range deployments {
		deploymentStatuses, err := e.fetchDeploymentStatus(project, deployment)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		deployments[index].Statuses = deploymentStatuses
	}

	deployments = e.filterDeploymentsByStatus(deployments)

	if len(deployments) > 0 {
		e.logger.Log("debug", "found new deployment events", "project", project)
	}

	return deployments, nil
}

// fetchDeploymentStatus fetches Deployment Statuses for the given Deployment.
func (e *GithubEventer) fetchDeploymentStatus(project string, deployment deployment) ([]deploymentStatus, error) {
	url := fmt.Sprintf(
		deploymentStatusUrlFormat,
		e.organisation,
		project,
		deployment.ID,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	e.rateLimiter.Wait()

	startTime := time.Now()

	resp, err := e.request(req)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	defer resp.Body.Close()

	updateDeploymentStatusMetrics("GET", e.organisation, project, resp.StatusCode, startTime)

	err = e.updateRateLimiter(resp)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, microerror.Maskf(unexpectedStatusCode, fmt.Sprintf("received non-200 status code: %v", resp.StatusCode))
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var deploymentStatuses []deploymentStatus
	if err := json.Unmarshal(bytes, &deploymentStatuses); err != nil {
		return nil, microerror.Mask(err)
	}

	return deploymentStatuses, nil
}

// postDeploymentStatus posts a Deployment Status for the given Deployment.
func (e *GithubEventer) postDeploymentStatus(project string, id int, state deploymentStatusState) error {
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
		return microerror.Mask(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return microerror.Mask(err)
	}

	e.rateLimiter.Wait()

	startTime := time.Now()

	resp, err := e.request(req)
	if err != nil {
		return microerror.Mask(err)
	}
	defer resp.Body.Close()

	updateDeploymentStatusMetrics("POST", e.organisation, project, resp.StatusCode, startTime)

	err = e.updateRateLimiter(resp)
	if err != nil {
		return microerror.Mask(err)
	}

	if resp.StatusCode != http.StatusCreated {
		return microerror.Maskf(unexpectedStatusCode, fmt.Sprintf("received non-200 status code: %v", resp.StatusCode))
	}

	return nil
}

// updateRateLimiter updates latest rate limiting token bucket values from
// response received.
func (e *GithubEventer) updateRateLimiter(response *http.Response) error {
	rateLimitRemaining, err := parseRateLimitRemaining(response)
	if err != nil {
		return microerror.Mask(err)
	}

	rateLimitResetTime, err := parseRateLimitResetTime(response)
	if err != nil {
		return microerror.Mask(err)
	}

	timeToRefill := rateLimitResetTime.Sub(time.Now())

	// If we are close to hit the rate limit, wait some extra time.
	if rateLimitRemaining < rateLimitAlmostHitThreshold {
		if rateLimitRemaining == 0 {
			// We must have at least one token in bucket. Otherwise it would
			// deadlock.
			rateLimitRemaining = 1
		}

		// If all tokens have been spent from GitHub API, doubling the wait
		// period + adding some extra is expected to make that 1 token be
		// available when the GitHub bucket gets refilled. When there's more
		// than 1 still left, this just makes spending of them slower.
		timeToRefill = 2*timeToRefill + rateLimitExtraWait
	}

	e.rateLimiter.Update(timeToRefill, int64(rateLimitRemaining))

	return nil
}

// parseRateLimitValue parses GitHub API rate limit value from response
// headers.
func parseRateLimitValue(response *http.Response) (float64, error) {
	rateLimitLimitValue, err := strconv.ParseFloat(response.Header.Get(rateLimitLimitHeader), 64)
	if err != nil {
		return 0.0, microerror.Mask(err)
	}
	return rateLimitLimitValue, nil
}

// parseRateLimitRemaining parses remaining GitHub API request tokens before
// rate limiting prevents further requests.
func parseRateLimitRemaining(response *http.Response) (float64, error) {
	rateLimitRemainingValue, err := strconv.ParseFloat(response.Header.Get(rateLimitRemainingHeader), 64)
	if err != nil {
		return 0.0, microerror.Mask(err)
	}

	return rateLimitRemainingValue, nil
}

// parseRateLimitResetTime parses time when GitHub API's rate limit bucket gets
// refilled.
func parseRateLimitResetTime(response *http.Response) (time.Time, error) {
	rateLimitResetValue, err := strconv.ParseInt(response.Header.Get(rateLimitResetHeader), 10, 64)
	if err != nil {
		return time.Time{}, microerror.Mask(err)
	}

	return time.Unix(rateLimitResetValue, 0), nil
}
