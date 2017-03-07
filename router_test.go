// Copyright 2015 Alex Browne and Soroush Pour.
// Allrights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package router

import (
	"reflect"
	"testing"
)

// routeTestCases is format for test case. Each test case takes possible route paths, the URL path provided
// and the expected route path and query parameters.
var routeTestCases = []struct {
	paths          []string
	path           string
	expectedPath   string
	expectedParams map[string]string
}{
	// Test route to /
	{
		paths:          []string{"/home", "/about", "/"},
		path:           "/",
		expectedPath:   "/",
		expectedParams: nil,
	},
	// Test basic literal
	{
		paths:          []string{"/home", "/about"},
		path:           "/home",
		expectedPath:   "/home",
		expectedParams: nil,
	},
	// Test trailing slash in path
	{
		paths:          []string{"/home", "/about"},
		path:           "/about/",
		expectedPath:   "/about",
		expectedParams: nil,
	},
	// Test query param
	{
		paths:        []string{"/home", "/home/{homeId}", "/about"},
		path:         "/home/55",
		expectedPath: "/home/{homeId}",
		expectedParams: map[string]string{
			"homeId": "55",
		},
	},
	// Test trailing slash after query param in path
	{
		paths:        []string{"/home", "/about", "/about/{aboutId}"},
		path:         "/about/55/",
		expectedPath: "/about/{aboutId}",
		expectedParams: map[string]string{
			"aboutId": "55",
		},
	},
	// Test multiple query params
	{
		paths:        []string{"/home", "/home/{homeId}", "/about", "/home/{homeId}/image/{imageSize}/{imageType}/jpg"},
		path:         "/home/55/image/800x600/panoramic/jpg",
		expectedPath: "/home/{homeId}/image/{imageSize}/{imageType}/jpg",
		expectedParams: map[string]string{
			"homeId":    "55",
			"imageSize": "800x600",
			"imageType": "panoramic",
		},
	},
	// Test tie breaker
	{
		paths:          []string{"/home/all", "/home/{homeId}", "/about"},
		path:           "/home/all",
		expectedPath:   "/home/all",
		expectedParams: nil,
	},
	// Test query param with weird characters
	{
		paths:        []string{"/home", "/home/{homeId}", "/about"},
		path:         "/home/@1.b$",
		expectedPath: "/home/{homeId}",
		expectedParams: map[string]string{
			"homeId": "@1.b$",
		},
	},
	// Test trailing slash after query param in path
	{
		paths:        []string{"/home", "/about", "/about/{aboutId}"},
		path:         "/about/@1.b$/",
		expectedPath: "/about/{aboutId}",
		expectedParams: map[string]string{
			"aboutId": "@1.b$",
		},
	},
	// Test multiple query params
	{
		paths:        []string{"/home", "/home/{homeId}", "/about", "/home/{homeId}/image/{imageSize}/{imageType}/jpg"},
		path:         "/home/@1.b$/image/800^600/%20a+/jpg",
		expectedPath: "/home/{homeId}/image/{imageSize}/{imageType}/jpg",
		expectedParams: map[string]string{
			"homeId":    "@1.b$",
			"imageSize": "800^600",
			"imageType": "%20a+",
		},
	},
}

func TestRouter(t *testing.T) {
	for _, tc := range routeTestCases {
		gotPath := ""
		gotParams := map[string]string{}
		r := New()
		for _, path := range tc.paths {
			handler := generateTestHandler(path, &gotPath, &gotParams)
			r.HandleFunc(path, handler)
		}
		r.pathChanged(tc.path, false)
		if gotPath != tc.expectedPath {
			t.Errorf("Failed for path=%s. Expected path: %s, Got path: %s", tc.path, tc.expectedPath, gotPath)
		}
		if tc.expectedParams != nil {
			if !reflect.DeepEqual(gotParams, tc.expectedParams) {
				t.Errorf("Failed for path=%s. Expected params: %s, Got params: %s", tc.path, tc.expectedParams, gotParams)
			}
		}
	}
}

func generateTestHandler(path string, gotPath *string, gotParams *map[string]string) Handler {
	return func(context *Context) {
		*gotPath = path
		*gotParams = context.Params
	}
}

func TestRemoveEmptyStrings(t *testing.T) {
	testCases := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			input:    []string{"a", "b", ""},
			expected: []string{"a", "b"},
		},
		{
			input:    []string{"a", "", "c"},
			expected: []string{"a", "c"},
		},
		{
			input:    []string{"", "b", "c"},
			expected: []string{"b", "c"},
		},
		{
			input:    []string{"", "", ""},
			expected: []string{},
		},
	}
	for i, tc := range testCases {
		got := removeEmptyStrings(tc.input)
		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("removeEmptyStrings failed for test case %d\nExpected: %v\nBut got  %v", i, tc.expected, got)
		}
	}
}
