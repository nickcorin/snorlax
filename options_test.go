package transit

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type OptionsTestSuite struct {
	suite.Suite
	client        *client
	server        *httptest.Server
	options       clientOptions
	targetOptions clientOptions
}

func (suite *OptionsTestSuite) SetupSuite() {
	suite.server = httptest.NewServer(http.HandlerFunc(echoHandler))
	suite.client = NewClient(WithBaseURL(suite.server.URL))
}

func (suite *OptionsTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *OptionsTestSuite) SetupTest() {
	suite.options = defaultOptions
	suite.targetOptions = clientOptions{
		baseURL: "https://www.example.com",
		headers: make(http.Header),
	}
	suite.targetOptions.headers.Set("TestKey", "TestValue")
}

type headerAdapter struct {
	Key string
	Value string
}

func (adapter headerAdapter) Adapt(r *http.Request) {
	r.Header.Set(adapter.Key, adapter.Value)
}

func (suite *OptionsTestSuite) TestWithAdapter() {
	adapter := headerAdapter{Key: "TestKey", Value: "TestValue"}
	suite.client = NewClient(WithAdapter(adapter), WithBaseURL(suite.server.URL))

	res, err := suite.client.Get("/", nil)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), res)
	require.Contains(suite.T(), res.Header.Values(adapter.Key), adapter.Value)
}

func (suite *OptionsTestSuite) TestWithBaseURL() {
	WithBaseURL(suite.targetOptions.baseURL)(&suite.options)
	require.Equal(suite.T(), suite.targetOptions.baseURL, suite.options.baseURL)
}

func (suite *OptionsTestSuite) TestWithHeader() {
	for k, v := range suite.targetOptions.headers {
		for _, header := range v {
			WithHeader(k, header)(&suite.options)
		}
	}

	for k := range suite.targetOptions.headers {
		assert.Equal(suite.T(), suite.targetOptions.headers.Get(k),
			suite.options.headers.Get(k))
	}
}

func TestOptionsTestSuite(t *testing.T) {
	suite.Run(t, new(OptionsTestSuite))
}
