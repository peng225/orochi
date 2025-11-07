package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/peng225/orochi/internal/entity"
	mgrclient "github.com/peng225/orochi/internal/manager/api/client"
	"github.com/stretchr/testify/require"
)

func TestBucketCreateGet_Success(t *testing.T) {
	c, err := mgrclient.NewClient("http://localhost:8080")
	require.NoError(t, err)
	bucketName := "test-bucket"
	reqBody := fmt.Sprintf(`{"name": "%s"}`, bucketName)
	createRes, err := c.CreateBucketWithBody(t.Context(), "application/json", strings.NewReader(reqBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, createRes.StatusCode)
	idStr := createRes.Header.Get("X-Bucket-ID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	require.NoError(t, err)

	getRes, err := c.GetBucket(t.Context(), id)
	require.NoError(t, err)
	defer getRes.Body.Close()
	require.Equal(t, http.StatusOK, getRes.StatusCode)
	data, err := io.ReadAll(getRes.Body)
	require.NoError(t, err)
	var got entity.Bucket
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)
	require.Equal(t, id, got.ID)
	require.Equal(t, bucketName, got.Name)
	require.Equal(t, "active", got.Status)
}
