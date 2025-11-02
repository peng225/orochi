package test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	gwclient "github.com/peng225/orochi/internal/gateway/api/client"
	mgrclient "github.com/peng225/orochi/internal/manager/api/client"
	"github.com/stretchr/testify/require"
)

func prepare(t *testing.T) {
	t.Helper()

	c, err := mgrclient.NewClient("http://localhost:8080")
	require.NoError(t, err)
	ports := []string{"8082", "8083", "8084"}
	for i, port := range ports {
		data := fmt.Sprintf(`{"baseURL": "http://datastore%d:%s"}`, i, port)
		res, err := c.CreateDatastoreWithBody(t.Context(), "application/json", strings.NewReader(data))
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)
	}

	data := `{"name": "test-bucket"}`
	res, err := c.CreateBucketWithBody(t.Context(), "application/json", strings.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestObjectCreateGet_Success(t *testing.T) {
	prepare(t)

	c, err := gwclient.NewClient("http://localhost:8081")
	require.NoError(t, err)

	bucket := "test-bucket"
	object := "test-object"
	createRes, err := c.CreateObjectWithBody(t.Context(), bucket, object,
		"application/octet-stream", strings.NewReader("test-data"))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, createRes.StatusCode)

	getRes, err := c.GetObject(t.Context(), bucket, object)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, getRes.StatusCode)
	defer getRes.Body.Close()
	data, err := io.ReadAll(getRes.Body)
	require.NoError(t, err)
	require.Equal(t, "test-data", string(data))
}
