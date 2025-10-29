package test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	mgrclient "github.com/peng225/orochi/internal/manager/api/client"
	"github.com/stretchr/testify/require"
)

func TestBucketCreate_Success(t *testing.T) {
	c, err := mgrclient.NewClient("http://localhost:8080")
	require.NoError(t, err)
	bucketName := "test-bucket"
	reqBody := fmt.Sprintf(`{"name": "%s"}`, bucketName)
	createRes, err := c.CreateBucketWithBody(t.Context(), "application/json", strings.NewReader(reqBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, createRes.StatusCode)
	idStr := createRes.Header.Get("X-Bucket-ID")
	require.NotEmpty(t, idStr)
}
