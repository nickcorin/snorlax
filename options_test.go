package snorlax

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type OptionsTestSuite struct {
	suite.Suite
	client *Client
	server *httptest.Server
}

func (suite *OptionsTestSuite) SetupSuite() {
	suite.server = httptest.NewServer(http.HandlerFunc(echoHandler))
	suite.client = New(&ClientOptions{
		BaseURL: suite.server.URL,
	})
}

func (suite *OptionsTestSuite) TestWithBaseURL() {
	url := "https://www.example.com"
	suite.client = New(&ClientOptions{
		BaseURL: url,
	})

	suite.Require().NotNil(suite.client)
	suite.Require().Equal(suite.client.opts.BaseURL, url)
}

func (suite *OptionsTestSuite) TestWithCallOptions() {
	username, password := "test", "12345"
	suite.client = New(nil)

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
		username, password)))

	res, err := suite.client.Get(context.TODO(), "https://www.example.com", nil,
		WithBasicAuth(username, password),
		WithHeader("custom-header", "test"))
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Contains(res.Request.Header.Get("Authorization"),
		fmt.Sprintf("Basic %s", auth))
}

func (suite *OptionsTestSuite) TearDownSuite() {
	suite.server.Close()
}

func TestOptionsTestSuite(t *testing.T) {
	suite.Run(t, new(OptionsTestSuite))
}
