package test

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"testing"

	gwclient "github.com/peng225/orochi/internal/gateway/api/client"
	"github.com/stretchr/testify/require"
)

func TestObject_CreateAndGet(t *testing.T) {
	testCases := []struct {
		ecConfig string
	}{
		{
			ecConfig: "2D1P",
		},
		{
			ecConfig: "2D2P",
		},
		{
			ecConfig: "3D1P",
		},
	}
	c, err := gwclient.NewClient(gatewayBaseURL)
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.ecConfig, func(t *testing.T) {
			bucket := prepareBucket(t, tc.ecConfig)
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
		})
	}
}

func TestObject_Delete(t *testing.T) {
	testCases := []struct {
		ecConfig string
	}{
		{
			ecConfig: "2D1P",
		},
		{
			ecConfig: "2D2P",
		},
		{
			ecConfig: "3D1P",
		},
	}
	c, err := gwclient.NewClient(gatewayBaseURL)
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.ecConfig, func(t *testing.T) {
			bucket := prepareBucket(t, tc.ecConfig)
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
		})
	}
}

func TestObject_List(t *testing.T) {
	c, err := gwclient.NewClient(gatewayBaseURL)
	require.NoError(t, err)

	bucket := prepareBucket(t, "2D1P")
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

	// Without query parameters.
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
		Limit: &limit,
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

	// With startFrom parameter.
	res3, err := c.ListObjects(t.Context(), bucket, &gwclient.ListObjectsParams{
		StartFrom: &nextObjectID,
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

func TestObject_Create_BucketDoesNotExist(t *testing.T) {
	c, err := gwclient.NewClient(gatewayBaseURL)
	require.NoError(t, err)

	bucket := "dummy"
	object := "test-object"
	resp, err := c.CreateObjectWithBody(t.Context(), bucket, object,
		"application/octet-stream", strings.NewReader("test-data"))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestObject_Create_InvalidParameter(t *testing.T) {
	c, err := gwclient.NewClient(gatewayBaseURL)
	require.NoError(t, err)

	testCases := []struct {
		name   string
		object string
	}{
		{
			name:   "contains _",
			object: "test_object",
		},
		{
			name:   "contains /",
			object: "foo/bar",
		},
		{
			name:   "contains .",
			object: "invalid.name",
		},
	}

	bucket := prepareBucket(t, "2D1P")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp, err := c.CreateObjectWithBody(t.Context(), bucket, tc.object,
				"application/octet-stream", strings.NewReader("test-data"))
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestObject_Get_BucketDoesNotExist(t *testing.T) {
	c, err := gwclient.NewClient(gatewayBaseURL)
	require.NoError(t, err)

	bucket := "dummy"
	object := "test-object"
	resp, err := c.GetObject(t.Context(), bucket, object)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestObject_Get_InvalidParameter(t *testing.T) {
	c, err := gwclient.NewClient(gatewayBaseURL)
	require.NoError(t, err)

	testCases := []struct {
		name   string
		object string
	}{
		{
			name:   "contains _",
			object: "test_object",
		},
		{
			name:   "contains /",
			object: "foo/bar",
		},
		{
			name:   "contains .",
			object: "invalid.name",
		},
	}

	bucket := prepareBucket(t, "2D1P")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp, err := c.GetObject(t.Context(), bucket, tc.object)
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestObject_Get_NotFound(t *testing.T) {
	c, err := gwclient.NewClient(gatewayBaseURL)
	require.NoError(t, err)

	bucket := prepareBucket(t, "2D1P")
	object := "test-object"

	resp, err := c.GetObject(t.Context(), bucket, object)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestObject_Delete_BucketDoesNotExist(t *testing.T) {
	c, err := gwclient.NewClient(gatewayBaseURL)
	require.NoError(t, err)

	bucket := "dummy"
	object := "test-object"
	resp, err := c.DeleteObject(t.Context(), bucket, object)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestObject_Delete_InvalidParameter(t *testing.T) {
	c, err := gwclient.NewClient(gatewayBaseURL)
	require.NoError(t, err)

	testCases := []struct {
		name   string
		object string
	}{
		{
			name:   "contains _",
			object: "test_object",
		},
		{
			name:   "contains /",
			object: "foo/bar",
		},
		{
			name:   "contains .",
			object: "invalid.name",
		},
	}

	bucket := prepareBucket(t, "2D1P")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp, err := c.DeleteObject(t.Context(), bucket, tc.object)
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestObject_List_BucketDoesNotExist(t *testing.T) {
	c, err := gwclient.NewClient(gatewayBaseURL)
	require.NoError(t, err)

	bucket := "dummy"
	resp, err := c.ListObjects(t.Context(), bucket, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestObject_List_TooLargeLimitParam(t *testing.T) {
	c, err := gwclient.NewClient(gatewayBaseURL)
	require.NoError(t, err)

	bucket := prepareBucket(t, "2D1P")
	limit := 1001
	resp, err := c.ListObjects(t.Context(), bucket, &gwclient.ListObjectsParams{
		Limit: &limit,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
