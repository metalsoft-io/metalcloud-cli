package job

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	sdk "github.com/metalsoft-io/metalcloud-sdk-go"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
)

var jobPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"JobId": {
			Title: "#",
			Order: 1,
		},
		"Status": {
			Order: 2,
		},
		"FunctionName": {
			Order: 3,
		},
		"CreatedTimestamp": {
			Title: "Created",
			Order: 4,
		},
		"JobGroupId": {
			Title: "Group",
			Order: 5,
		},
	},
}

type ListFlags struct {
	FilterJobId      []string
	FilterStatus     []string
	FilterJobGroupId []string
	SortBy           []string
}

func JobList(ctx context.Context, flags ListFlags) error {
	logger.Get().Info().Msg("Listing jobs")

	client := api.GetApiClient(ctx)
	request := client.JobAPI.GetJobs(ctx)

	if len(flags.FilterJobId) > 0 {
		request = request.FilterJobId(flags.FilterJobId)
	}
	if len(flags.FilterStatus) > 0 {
		request = request.FilterStatus(flags.FilterStatus)
	}
	if len(flags.FilterJobGroupId) > 0 {
		request = request.FilterJobGroupId(flags.FilterJobGroupId)
	}
	if len(flags.SortBy) > 0 {
		request = request.SortBy(flags.SortBy)
	}

	jobs, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(jobs, &jobPrintConfig)
}

func JobGet(ctx context.Context, jobId string) error {
	logger.Get().Info().Msgf("Get job '%s' details", jobId)

	id, err := getJobId(jobId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)
	job, httpRes, err := client.JobAPI.GetJob(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(job, &jobPrintConfig)
}

func JobSkip(ctx context.Context, jobId string) error {
	logger.Get().Info().Msgf("Skipping job '%s'", jobId)

	if err := validateJobId(jobId); err != nil {
		return err
	}

	httpRes, err := jobAction(ctx, jobId, "skip", nil)
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Printf("Job %s has been skipped.\n", jobId)
	return nil
}

func JobRetry(ctx context.Context, jobId string, retryEvenIfSuccessful bool) error {
	logger.Get().Info().Msgf("Retrying job '%s'", jobId)

	if err := validateJobId(jobId); err != nil {
		return err
	}

	retryInfo := sdk.NewJobRetryInfo()
	if retryEvenIfSuccessful {
		retryInfo.SetRetryEvenIfSuccessful(retryEvenIfSuccessful)
	}

	httpRes, err := jobAction(ctx, jobId, "retry", retryInfo)
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Printf("Job %s has been retried.\n", jobId)
	return nil
}

func JobKill(ctx context.Context, jobId string) error {
	logger.Get().Info().Msgf("Killing job '%s'", jobId)

	if err := validateJobId(jobId); err != nil {
		return err
	}

	commandInfo := sdk.NewJobCommandInfo()
	commandInfo.SetCommand("kill")
	commandInfo.SetExecuteImmediately(true)

	httpRes, err := jobAction(ctx, jobId, "issue-command", commandInfo)
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Printf("Job %s has been killed.\n", jobId)
	return nil
}

// jobAction performs a POST to /api/v2/jobs/{jobId}/actions/{action} bypassing
// the SDK's float32 jobId parameter which loses precision for IDs > 16,777,216.
func jobAction(ctx context.Context, jobId string, action string, body interface{}) (*http.Response, error) {
	client := api.GetApiClient(ctx)
	cfg := client.GetConfig()

	baseURL, err := cfg.ServerURL(0, nil)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v2/jobs/%s/actions/%s", baseURL, jobId, action)

	var reqBody *bytes.Buffer
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(data)
	} else {
		reqBody = &bytes.Buffer{}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if auth, ok := ctx.Value(sdk.ContextAccessToken).(string); ok {
		req.Header.Set("Authorization", "Bearer "+auth)
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return resp, err
	}

	// Read and re-wrap the body so that response_inspector.InspectResponse
	// can format it properly (it uses %s on httpRes.Body).
	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
	if err != nil {
		return resp, err
	}

	if resp.StatusCode >= 300 {
		return resp, fmt.Errorf("%s - %s", resp.Status, string(respBody))
	}

	return resp, nil
}

func validateJobId(jobId string) error {
	if _, err := strconv.Atoi(jobId); err != nil {
		err := fmt.Errorf("invalid job ID: '%s'", jobId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}
	return nil
}

func getJobId(jobId string) (float32, error) {
	id, err := strconv.ParseFloat(jobId, 32)
	if err != nil {
		err := fmt.Errorf("invalid job ID: '%s'", jobId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}
	return float32(id), nil
}
