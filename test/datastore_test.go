package test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	mgrclient "github.com/peng225/orochi/internal/manager/api/client"
	"github.com/stretchr/testify/require"
)

func TestDatastore_Create_BadRequest(t *testing.T) {
	c, err := mgrclient.NewClient("http://localhost:8080")
	require.NoError(t, err)
	baseURL := "ftp://invalid-datastore"
	reqBody := fmt.Sprintf(`{"baseURL": "%s"}`, baseURL)
	createResp, err := c.CreateDatastoreWithBody(t.Context(), "application/json", strings.NewReader(reqBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, createResp.StatusCode)
}

func TestDatastore_Get_NotFound(t *testing.T) {
	c, err := mgrclient.NewClient("http://localhost:8080")
	require.NoError(t, err)
	resp, err := c.GetDatastore(t.Context(), int64(1000))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
