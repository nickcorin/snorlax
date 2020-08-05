package transit

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
	client *client
	server *httptest.Server
}

func (suite *ClientTestSuite) SetupSuite() {
	suite.server = httptest.NewServer(http.HandlerFunc(echoHandler))
	suite.client = NewClient(WithBaseURL(suite.server.URL)).(*client)
}

func (suite *ClientTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *ClientTestSuite) TestNewClient() {
	client := NewClient().(*client)
	assert.ObjectsAreEqualValues(client.opts, defaultOptions)
}

func (suite *ClientTestSuite) TestClient_Delete() {
	body := []byte("test")
	res, err := suite.client.Delete("/", nil, bytes.NewBuffer(body))
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), res)

	defer res.Body.Close()
	responseBody, err := ioutil.ReadAll(res.Body)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), []byte("test"), responseBody)
}

func (suite *ClientTestSuite) TestClient_Get() {
	res, err := suite.client.Get("/", nil)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), res)

	defer res.Body.Close()
	require.NoError(suite.T(), err)
}

func (suite *ClientTestSuite) TestClient_Post() {
	body := []byte("test")
	res, err := suite.client.Post("/", nil, bytes.NewBuffer(body))
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), res)

	defer res.Body.Close()
	responseBody, err := ioutil.ReadAll(res.Body)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), []byte("test"), responseBody)
}

func (suite *ClientTestSuite) TestClient_Put() {
	body := []byte("test")
	res, err := suite.client.Put("/", nil, bytes.NewBuffer(body))
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), res)

	defer res.Body.Close()
	responseBody, err := ioutil.ReadAll(res.Body)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), []byte("test"), responseBody)
}

func (suite *ClientTestSuite) TestClient_Do() {
	body := []byte("test")
	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), req)

	res, err := suite.client.Do(req)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), res)
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), responseBody)
	require.Equal(suite.T(), http.StatusOK, res.StatusCode)
	require.Equal(suite.T(), body, responseBody)
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}