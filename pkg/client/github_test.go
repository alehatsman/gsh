package client

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClientInitialization(t *testing.T) {
	gClient := New()
	if gClient == nil {
		t.Errorf("New should return new GithubClient")
	}
}

func TestSearch(t *testing.T) {
	testCases := []struct {
		query            string
		responseFilePath string
		expectedResult   *SearchResult
	}{
		{"", "../../test/search/empty_response.json", &SearchResult{TotalCount: 0, Items: []Repo{}}},
		{"simple", "../../test/search/simple_response.json", &SearchResult{TotalCount: 1, Items: []Repo{{Name: "react", URL: "https://github.com/facebook/react", Description: "A declarative, efficient, and flexible JavaScript library for building user interfaces.", Stars: 121666}}}},
		{"bad", "../../test/search/bad_response.json", nil},
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		query := req.URL.Query().Get("q")

		var filePath string
		for _, testCase := range testCases {
			if testCase.query == query {
				filePath = testCase.responseFilePath
			}
		}

		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		rw.Write(content)
	}))
	defer server.Close()

	gClient := NewWithParams(server.URL, server.URL)

	for _, testCase := range testCases {
		actualResult, _ := gClient.Search(testCase.query)

		if !reflect.DeepEqual(actualResult, testCase.expectedResult) {
			t.Error("Wrong response", actualResult, testCase.expectedResult)
		}
	}
}

func TestTrending(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		content, err := ioutil.ReadFile("../../test/trending/simple.html")
		if err != nil {
			panic(err)
		}

		rw.Write(content)
	}))
	defer server.Close()

	gClient := NewWithParams(server.URL, server.URL)

	sr, err := gClient.Trending("daily")
	if err != nil {
		t.Error(err)
	}

	if sr.TotalCount != 25 {
		t.Errorf("Expected to receive count %d, but got %d", 25, sr.TotalCount)
	}
}
