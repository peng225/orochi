package test

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/peng225/orochi/internal/entity"
	mgrclient "github.com/peng225/orochi/internal/manager/api/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatastore_Create_BadRequest(t *testing.T) {
	c, err := mgrclient.NewClient(managerBaseURL)
	require.NoError(t, err)
	baseURL := "ftp://invalid-datastore"
	reqBody := fmt.Sprintf(`{"baseURL": "%s"}`, baseURL)
	createResp, err := c.CreateDatastoreWithBody(t.Context(), "application/json", strings.NewReader(reqBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, createResp.StatusCode)
}

func TestDatastore_Get_NotFound(t *testing.T) {
	c, err := mgrclient.NewClient(managerBaseURL)
	require.NoError(t, err)
	resp, err := c.GetDatastore(t.Context(), int64(1000))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestDatastore_List(t *testing.T) {
	c, err := mgrclient.NewClient(managerBaseURL)
	require.NoError(t, err)

	// Without query parameters.
	resp1, err := c.ListDatastores(t.Context(), nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp1.StatusCode)
	defer resp1.Body.Close()
	data, err := io.ReadAll(resp1.Body)
	require.NoError(t, err)
	dss := make([]*entity.Datastore, 0)
	err = json.Unmarshal(data, &dss)
	require.NoError(t, err)
	for _, ds := range dss {
		assert.Equal(t, entity.DatastoreStatusActive, ds.Status)
	}
	require.Equal(t, strconv.FormatInt(math.MaxInt64, 10), resp1.Header.Get("X-Next-Datastore-ID"))

	// With limit parameter.
	limit := 2
	resp2, err := c.ListDatastores(t.Context(), &mgrclient.ListDatastoresParams{
		Limit: &limit,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)
	defer resp2.Body.Close()
	data, err = io.ReadAll(resp2.Body)
	require.NoError(t, err)
	dss = make([]*entity.Datastore, 0)
	err = json.Unmarshal(data, &dss)
	require.NoError(t, err)
	require.Len(t, dss, limit)
	for _, ds := range dss {
		assert.Equal(t, entity.DatastoreStatusActive, ds.Status)
	}
	nextDatastoreID, err := strconv.ParseInt(resp2.Header.Get("X-Next-Datastore-ID"), 10, 64)
	require.NoError(t, err)

	// With startFrom parameter.
	resp3, err := c.ListDatastores(t.Context(), &mgrclient.ListDatastoresParams{
		StartFrom: &nextDatastoreID,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp3.StatusCode)
	defer resp3.Body.Close()
	data, err = io.ReadAll(resp3.Body)
	require.NoError(t, err)
	dss = make([]*entity.Datastore, 0)
	err = json.Unmarshal(data, &dss)
	require.NoError(t, err)
	require.Less(t, 0, len(dss))
}
