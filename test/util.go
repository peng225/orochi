package test

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	randv2 "math/rand/v2"

	mgrclient "github.com/peng225/orochi/internal/manager/api/client"
)

func prepareBucket(t *testing.T) string {
	t.Helper()

	c, err := mgrclient.NewClient("http://localhost:8080")
	require.NoError(t, err)

	bucketName := fmt.Sprintf("test-bucket-%s", generateRandomStr(t, 8))
	data := fmt.Sprintf(`{"name": "%s"}`, bucketName)
	resp, err := c.CreateBucketWithBody(t.Context(), "application/json", strings.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	bucketIDStr := resp.Header.Get("X-Bucket-ID")
	bucketID, err := strconv.ParseInt(bucketIDStr, 10, 64)
	require.NoError(t, err)
	t.Cleanup(func() {
		teardownBucket(t, bucketID)
	})

	return bucketName
}

func teardownBucket(t *testing.T, bucketID int64) {
	t.Helper()

	c, err := mgrclient.NewClient("http://localhost:8080")
	require.NoError(t, err)

	resp, err := c.DeleteBucket(context.Background(), bucketID)
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, resp.StatusCode)
}

func generateRandomStr(t *testing.T, digit int) string {
	t.Helper()

	var result string
	for range digit {
		// Convert to the string consisting of a-z.
		result += string(byte(randv2.IntN(26)) + 97)
	}
	return result
}
