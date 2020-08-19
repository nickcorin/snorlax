package snorlax_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nickcorin/snorlax"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
	client *snorlax.Client
	server *httptest.Server
}

func (suite *ClientTestSuite) SetupSuite() {
	h := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	}

	suite.server = httptest.NewServer(http.HandlerFunc(h))
	suite.client = snorlax.NewClient(&snorlax.ClientOptions{
		BaseURL: suite.server.URL,
	})
}

func (suite *ClientTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *ClientTestSuite) TestClient_Delete() {
	res, err := suite.client.Delete(context.TODO(), "/delete", nil, nil)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), res)
	require.Equal(suite.T(), http.StatusOK, res.StatusCode)
}

func (suite *ClientTestSuite) TestClient_Get() {
	res, err := suite.client.Get(context.TODO(), "/get", nil)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), res)
	require.Equal(suite.T(), http.StatusOK, res.StatusCode)
}

func (suite *ClientTestSuite) TestClient_Post() {
	res, err := suite.client.Post(context.TODO(), "/post", nil, nil)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), res)
	require.Equal(suite.T(), http.StatusOK, res.StatusCode)
}

func (suite *ClientTestSuite) TestClient_Put() {
	res, err := suite.client.Put(context.TODO(), "/put", nil, nil)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), res)
	require.Equal(suite.T(), http.StatusOK, res.StatusCode)
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
