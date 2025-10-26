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

func TestDatastoreCreateGet_Success(t *testing.T) {
	// Create
	c, err := mgrclient.NewClient("http://localhost:8080")
	require.NoError(t, err)
	baseURL := "http://datastore"
	reqBody := fmt.Sprintf(`{"baseURL": "%s"}`, baseURL)
	createRes, err := c.CreateDatastoreWithBody(t.Context(), "application/json", strings.NewReader(reqBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, createRes.StatusCode)
	idStr := createRes.Header.Get("X-Datastore-ID")
	require.NotEmpty(t, idStr)

	// Get
	id, err := strconv.ParseInt(idStr, 10, 64)
	require.NoError(t, err)
	getRes, err := c.GetDatastore(t.Context(), id)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, getRes.StatusCode)
	resBody, err := io.ReadAll(getRes.Body)
	require.NoError(t, err)
	var ds entity.Datastore
	err = json.Unmarshal(resBody, &ds)
	require.NoError(t, err)
	require.Equal(t, id, ds.ID)
	require.Equal(t, baseURL, ds.BaseURL)
}
