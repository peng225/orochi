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

func TestObject_CreateAndGet(t *testing.T) {
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

func TestObject_Delete(t *testing.T) {
	prepare(t)

	c, err := gwclient.NewClient("http://localhost:8081")
	require.NoError(t, err)

	bucket := "test-bucket"
	object := "test-object2"
	createRes, err := c.CreateObjectWithBody(t.Context(), bucket, object,
		"application/octet-stream", strings.NewReader("test-data"))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, createRes.StatusCode)

	delRes, err := c.DeleteObject(t.Context(), bucket, object)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, delRes.StatusCode)
	// Check idempotency.
	delRes, err = c.DeleteObject(t.Context(), bucket, object)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, delRes.StatusCode)

	getRes, err := c.GetObject(t.Context(), bucket, object)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, getRes.StatusCode)
}

func TestObject_List(t *testing.T) {
	prepare(t)

	c, err := gwclient.NewClient("http://localhost:8081")
	require.NoError(t, err)

	bucket := "test-bucket"
	objects := []string{
		"test-object0",
		"test-object1",
		"test-object2",
	}
	for _, o := range objects {
		createRes, err := c.CreateObjectWithBody(t.Context(), bucket, o,
			"application/octet-stream", strings.NewReader("test-data"))
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, createRes.StatusCode)
	}

	// Without limit parameter.
	res1, err := c.ListObjects(t.Context(), bucket, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res1.StatusCode)
	defer res1.Body.Close()
	data, err := io.ReadAll(res1.Body)
	require.NoError(t, err)
	objList := make([]string, 0, 3)
	err = json.Unmarshal(data, &objList)
	require.NoError(t, err)
	require.Len(t, objList, len(objects))
	for i, o := range objList {
		require.Equal(t, objects[i], o)
	}
	require.Equal(t, strconv.FormatInt(math.MaxInt64, 10), res1.Header.Get("X-Next-Object-ID"))

	// With limit parameter.
	limit := 2
	res2, err := c.ListObjects(t.Context(), bucket, &gwclient.ListObjectsParams{
		XLimit: &limit,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res2.StatusCode)
	defer res2.Body.Close()
	data, err = io.ReadAll(res2.Body)
	require.NoError(t, err)
	objList = make([]string, 0, 2)
	err = json.Unmarshal(data, &objList)
	require.NoError(t, err)
	require.Len(t, objList, limit)
	for i, o := range objList {
		require.Equal(t, objects[i], o)
	}
	nextObjectID, err := strconv.ParseInt(res2.Header.Get("X-Next-Object-ID"), 10, 64)
	require.NoError(t, err)

	// With firstObjectID parameter.
	res3, err := c.ListObjects(t.Context(), bucket, &gwclient.ListObjectsParams{
		XFirstObjectID: &nextObjectID,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res3.StatusCode)
	defer res3.Body.Close()
	data, err = io.ReadAll(res3.Body)
	require.NoError(t, err)
	objList = make([]string, 0, 1)
	err = json.Unmarshal(data, &objList)
	require.NoError(t, err)
	require.Len(t, objList, 1)
	require.Equal(t, objects[2], objList[0])
}
