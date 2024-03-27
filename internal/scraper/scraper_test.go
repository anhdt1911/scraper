package scraper

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildUrl(t *testing.T) {
	tests := []struct {
		name     string
		keyword  string
		expected string
	}{
		{
			name:     "keyword with leading and trailing spaces",
			keyword:  "  buy domain   ",
			expected: "https://www.google.com/search?q=buy+domain",
		},
		{
			name:     "keyword with multiple spaces",
			keyword:  "sell   your   domain",
			expected: "https://www.google.com/search?q=sell+your+domain",
		},
		{
			name:     "keyword without spaces",
			keyword:  "test.id",
			expected: "https://www.google.com/search?q=test.id",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := buildUrl(test.keyword)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestScraper_Request(t *testing.T) {
	// Create a test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a mock response
		w.WriteHeader(http.StatusOK)
	}))

	defer testServer.Close()

	// Create a new instance of the Scraper
	scraper := &Scraper{
		client: testServer.Client(),
	}

	// Test cases
	tests := []struct {
		name           string
		url            string
		expectedStatus int
		isError        assert.ErrorAssertionFunc
	}{
		{
			name:           "valid request",
			url:            testServer.URL,
			expectedStatus: http.StatusOK,
			isError:        assert.NoError,
		},
		// {
		// 	name:           "invalid URL",
		// 	url:            "invalid-url",
		// 	expectedStatus: 0,
		// 	isError:        assert.Error,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := scraper.request(test.url)

			assert.Equal(t, test.expectedStatus, res.StatusCode)

			test.isError(t, err)
		})
	}
}

// func TestScraper_Scrape(t *testing.T) {
// 	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		htmlContent := `
// 			<html>
// 				<body>
// 					<a href="http://example.com">Link 1</a>
// 					<a href="http://example.com">Link 2</a>
// 					<a href="https://example.com">Link 3</a>
// 					<div class="result-stats">10 results</div>
// 					<span>Sponsored</span>
// 					<span>Ads</span>
// 					<span>Other</span>
// 				</body>
// 			</html>
// 		`
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(htmlContent))
// 	}))

// 	defer testServer.Close()

// 	scraper := &Scraper{
// 		client: testServer.Client(),
// 	}

// 	tests := []struct {
// 		name           string
// 		keyword        string
// 		expectedResult *SearchResult
// 		expectedError  assert.ErrorAssertionFunc
// 	}{
// 		{
// 			name:    "valid keyword",
// 			keyword: "example",
// 			expectedResult: &SearchResult{
// 				Keyword:           "example",
// 				TotalSearchResult: "Khoảng 10.570.000.000 kết quả",
// 				LinkAmount:        18,
// 				AdwordAMount:      0,
// 			},
// 			expectedError: assert.NoError,
// 		},
// 		{
// 			name:           "invalid keyword",
// 			keyword:        "",
// 			expectedResult: nil,
// 			expectedError:  assert.Error,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			res, err := scraper.Scrape(test.keyword)

// 			assert.Equal(t, test.expectedError, err)

// 			if err == nil {
// 				assert.Equal(t, test.expectedResult.Keyword, res.Keyword)
// 				assert.Equal(t, test.expectedResult.TotalSearchResult, res.TotalSearchResult)
// 				assert.Equal(t, test.expectedResult.LinkAmount, res.LinkAmount)
// 				assert.Equal(t, test.expectedResult.AdwordAMount, res.AdwordAMount)
// 			}
// 		})
// 	}
// }
