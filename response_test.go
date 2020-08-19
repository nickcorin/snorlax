package snorlax

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ResponseTestSuite struct {
	suite.Suite
	client *Client
	server *httptest.Server
}

func (suite *ResponseTestSuite) SetupSuite() {
	suite.server = httptest.NewServer(http.HandlerFunc(echoHandler))
	suite.client = New(&ClientOptions{
		BaseURL: suite.server.URL,
	})
}

func (suite *ResponseTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *ResponseTestSuite) TestIsSuccess() {
	successResponse := Response{http.Response{
		StatusCode: http.StatusOK,
	}}

	failedResponse := Response{http.Response{
		StatusCode: http.StatusInternalServerError,
	}}

	suite.Require().True(successResponse.IsSuccess())
	suite.Require().False(failedResponse.IsSuccess())
}

func (suite *ResponseTestSuite) TestJSON() {
	type Pokemon struct {
		Name   string `json:"name"`
		Number int    `json:"number"`
	}

	body := []byte(`{"name": "snorlax", "number": 143}`)
	res, err := suite.client.Post(context.TODO(), "/example", nil,
		bytes.NewBuffer(body))
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	var pokemon Pokemon
	err = res.JSON(&pokemon)
	suite.Require().NoError(err)
	suite.Require().Equal("snorlax", pokemon.Name)
	suite.Require().Equal(143, pokemon.Number)
}

func (suite *ResponseTestSuite) TestRawBody() {
	type Pokemon struct {
		Name   string `json:"name"`
		Number int    `json:"number"`
	}

	body := []byte(`{"name": "snorlax", "number": 143}`)
	res, err := suite.client.Post(context.TODO(), "/example", nil,
		bytes.NewBuffer(body))
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	responseReader, err := res.RawBody()
	suite.Require().NoError(err)
	suite.Require().NotNil(responseReader)

	responseBody, err := ioutil.ReadAll(responseReader)
	suite.Require().NoError(err)
	suite.Require().NotNil(responseBody)
	suite.Require().EqualValues(body, responseBody)
}

func TestResponseTestSuite(t *testing.T) {
	suite.Run(t, new(ResponseTestSuite))
}
