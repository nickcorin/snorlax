package snorlax_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nickcorin/snorlax"
	"github.com/stretchr/testify/suite"
)

type ResponseTestSuite struct {
	suite.Suite
	client snorlax.Client
	server *httptest.Server
}

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	for headerKey, headerValues := range r.Header {
		for _, headerValue := range headerValues {
			w.Header().Add(headerKey, headerValue)
		}
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (suite *ResponseTestSuite) SetupSuite() {
	suite.server = httptest.NewServer(http.HandlerFunc(EchoHandler))
	suite.client = snorlax.Client{
		BaseURL: suite.server.URL,
	}
}

func (suite *ResponseTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *ResponseTestSuite) TestIsSuccess() {
	successResponse := snorlax.Response{http.Response{
		StatusCode: http.StatusOK,
	}}

	failedResponse := snorlax.Response{http.Response{
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
