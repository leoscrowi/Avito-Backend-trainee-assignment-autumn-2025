package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/leoscrowi/pr-assignment-service/domain"
	"github.com/leoscrowi/pr-assignment-service/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPullRequestMerge_Success(t *testing.T) {
	teamName := "test_pr_team_merge_success"
	team := map[string]interface{}{
		"team_name": teamName,
		"members": []map[string]interface{}{
			{"user_id": "test_u1_merge_success", "username": "TestAlice", "is_active": true},
			{"user_id": "test_u2_merge_success", "username": "TestBob", "is_active": true},
			{"user_id": "test_u3_merge_success", "username": "TestCharlie", "is_active": true},
		},
	}

	respAdd := helpers.PostJSON(t, "/team/add", team, helpers.AdminToken)
	_ = respAdd.Body.Close()
	helpers.RequireStatusCode(t, respAdd, http.StatusCreated)

	pr := map[string]interface{}{
		"pull_request_id":   "pr-1001_merge_success",
		"pull_request_name": "Add search feature",
		"author_id":         "test_u1_merge_success",
	}

	createResp := helpers.PostJSON(t, "/pullRequest/create", pr, helpers.AdminToken)
	_ = createResp.Body.Close()
	require.Equal(t, http.StatusCreated, createResp.StatusCode)

	mergeRequest := map[string]interface{}{
		"pull_request_id": "pr-1001_merge_success",
	}

	resp := helpers.PatchJSON(t, "/pullRequest/merge", mergeRequest, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var out struct {
		PR domain.PullRequest `json:"pr"`
	}
	err := json.NewDecoder(resp.Body).Decode(&out)
	require.NoError(t, err)

	assert.Equal(t, "pr-1001_merge_success", out.PR.PullRequestID)
	assert.Equal(t, "Add search feature", out.PR.PullRequestName)
	assert.Equal(t, "test_u1_merge_success", out.PR.AuthorID)
	assert.Equal(t, domain.MERGED, out.PR.Status)
	assert.NotNil(t, out.PR.MergedAt)
	assert.WithinDuration(t, time.Now(), out.PR.MergedAt, 5*time.Second)
}

func TestPullRequestMerge_Idempotent(t *testing.T) {
	teamName := "test_pr_team_merge_idempotent"
	team := map[string]interface{}{
		"team_name": teamName,
		"members": []map[string]interface{}{
			{"user_id": "test_u1_merge_idempotent", "username": "TestAlice", "is_active": true},
			{"user_id": "test_u2_merge_idempotent", "username": "TestBob", "is_active": true},
		},
	}

	respAdd := helpers.PostJSON(t, "/team/add", team, helpers.AdminToken)
	_ = respAdd.Body.Close()
	helpers.RequireStatusCode(t, respAdd, http.StatusCreated)

	pr := map[string]interface{}{
		"pull_request_id":   "pr-1002_merge_idempotent",
		"pull_request_name": "Fix bug",
		"author_id":         "test_u1_merge_idempotent",
	}

	createResp := helpers.PostJSON(t, "/pullRequest/create", pr, helpers.AdminToken)
	_ = createResp.Body.Close()
	require.Equal(t, http.StatusCreated, createResp.StatusCode)

	mergeRequest := map[string]interface{}{
		"pull_request_id": "pr-1002_merge_idempotent",
	}

	resp1 := helpers.PatchJSON(t, "/pullRequest/merge", mergeRequest, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp1.Body)
	require.Equal(t, http.StatusOK, resp1.StatusCode)

	var out1 struct {
		PR domain.PullRequest `json:"pr"`
	}
	err := json.NewDecoder(resp1.Body).Decode(&out1)
	require.NoError(t, err)
	assert.Equal(t, domain.MERGED, out1.PR.Status)
	firstMergeTime := out1.PR.MergedAt

	resp2 := helpers.PatchJSON(t, "/pullRequest/merge", mergeRequest, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp2.Body)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	var out2 struct {
		PR domain.PullRequest `json:"pr"`
	}
	err = json.NewDecoder(resp2.Body).Decode(&out2)
	require.NoError(t, err)

	assert.Equal(t, domain.MERGED, out2.PR.Status)
	assert.Equal(t, firstMergeTime, out2.PR.MergedAt, "merge time should not change on subsequent calls")
}

func TestPullRequestMerge_NotFound(t *testing.T) {
	mergeRequest := map[string]interface{}{
		"pull_request_id": "pr-nonexistent_merge_notfound",
	}

	resp := helpers.PatchJSON(t, "/pullRequest/merge", mergeRequest, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	var errorResp struct {
		Error struct {
			Code    domain.ErrorCode `json:"code"`
			Message string           `json:"message"`
		} `json:"error"`
	}

	err := json.Unmarshal(helpers.ReadBody(t, resp), &errorResp)
	require.NoError(t, err)

	assert.True(t, errorResp.Error.Code == domain.NOT_FOUND,
		"expected NOT_FOUND, got: ", errorResp.Error.Code)
}

func TestPullRequestMerge_MissingPullRequestID(t *testing.T) {
	mergeRequest := map[string]interface{}{}

	resp := helpers.PatchJSON(t, "/pullRequest/merge", mergeRequest, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var errorResp struct {
		Error struct {
			Code    domain.ErrorCode `json:"code"`
			Message string           `json:"message"`
		} `json:"error"`
	}

	err := json.Unmarshal(helpers.ReadBody(t, resp), &errorResp)
	require.NoError(t, err)

	assert.True(t, errorResp.Error.Code == domain.BAD_REQUEST,
		"expected INVALID_INPUT, got: ", errorResp.Error.Code)
}
