package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/peng225/orochi/internal/entity"
	mgrclient "github.com/peng225/orochi/internal/manager/api/client"
	"github.com/stretchr/testify/require"
)

func TestBucket_CreateGetDelete(t *testing.T) {
	c, err := mgrclient.NewClient(managerBaseURL)
	require.NoError(t, err)
	bucketName := "test-bucket"
	reqBody := fmt.Sprintf(`{"name": "%s", "ecConfig": "2D1P"}`, bucketName)
	createResp, err := c.CreateBucketWithBody(t.Context(), "application/json", strings.NewReader(reqBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, createResp.StatusCode)
	idStr := createResp.Header.Get("X-Bucket-ID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	require.NoError(t, err)

	getResp, err := c.GetBucket(t.Context(), id)
	require.NoError(t, err)
	defer getResp.Body.Close()
	require.Equal(t, http.StatusOK, getResp.StatusCode)
	data, err := io.ReadAll(getResp.Body)
	require.NoError(t, err)
	var got entity.Bucket
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)
	require.Equal(t, id, got.ID)
	require.Equal(t, bucketName, got.Name)
	require.Equal(t, "active", got.Status)

	deleteResp, err := c.DeleteBucket(t.Context(), id)
	require.NoError(t, err)
	defer deleteResp.Body.Close()
	require.Equal(t, http.StatusAccepted, deleteResp.StatusCode)
	require.Eventually(t, func() bool {
		getResp, err := c.GetBucket(t.Context(), id)
		if err != nil {
			return false
		}
		defer getResp.Body.Close()
		_, err = io.Copy(io.Discard, getResp.Body)
		if err != nil {
			return false
		}
		return getResp.StatusCode == http.StatusNotFound
	}, 10*time.Second, 1*time.Second)
}

func TestBucket_Create_BadRequest(t *testing.T) {
	c, err := mgrclient.NewClient(managerBaseURL)
	require.NoError(t, err)
	bucketName := "test-bucket//"
	reqBody := fmt.Sprintf(`{"name": "%s", "ecConfig": "2D1P"}`, bucketName)
	createResp, err := c.CreateBucketWithBody(t.Context(), "application/json", strings.NewReader(reqBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, createResp.StatusCode)
}

func TestBucket_Get_NotFound(t *testing.T) {
	c, err := mgrclient.NewClient(managerBaseURL)
	require.NoError(t, err)
	resp, err := c.GetBucket(t.Context(), int64(1000))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
