package snorlax_test

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nickcorin/snorlax"
	"github.com/stretchr/testify/suite"
)

type HooksTestSuite struct {
	suite.Suite
	client snorlax.Client
}

func (suite *HooksTestSuite) SetupSuite() {
	suite.client = snorlax.DefaultClient
}

func (suite *HooksTestSuite) TestWithBasicAuth() {
	username, password := "snorlax", "s3cr3t"
	headerKey := http.CanonicalHeaderKey("Authorization")

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	suite.Require().Empty(r.Header.Get(headerKey))

	snorlax.WithBasicAuth(username, password)(suite.client, r)
	headerValue := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
		username, password)))
	suite.Require().Equal(fmt.Sprintf("Basic %s", headerValue),
		r.Header.Get(headerKey))
}

func (suite *HooksTestSuite) TestWithHeader() {
	headerKey, headerValue := http.CanonicalHeaderKey("pokemon"), "snorlax"

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	suite.Require().Empty(r.Header.Get(headerKey))

	snorlax.WithHeader(headerKey, headerValue)(suite.client, r)
	suite.Require().Equal(headerValue, r.Header.Get(headerKey))
}

func TestHooksTestSuite(t *testing.T) {
	suite.Run(t, new(HooksTestSuite))
}
