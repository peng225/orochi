package test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	gwclient "github.com/peng225/orochi/internal/gateway/api/client"
	mgrclient "github.com/peng225/orochi/internal/manager/api/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prepare(t *testing.T) {
	t.Helper()

	c, err := mgrclient.NewClient("http://localhost:8080")
	require.NoError(t, err)
	data := `{"baseURL": "http://datastore:8082"}`
	res, err := c.CreateDatastoreWithBody(t.Context(), "application/json", strings.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestObjectCreateGet(t *testing.T) {
	prepare(t)

	c, err := gwclient.NewClient("http://localhost:8081")
	require.NoError(t, err)

	bucket := "test-bucket"
	object := "test-object"
	createRes, err := c.CreateObjectWithBody(t.Context(), bucket, object,
		"application/octet-stream", strings.NewReader("test-data"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, createRes.StatusCode)

	getRes, err := c.GetObject(t.Context(), bucket, object)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, getRes.StatusCode)
	defer getRes.Body.Close()
	data, err := io.ReadAll(getRes.Body)
	require.NoError(t, err)
	require.Equal(t, "test-data", string(data))
}
