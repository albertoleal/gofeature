// Copyright 2015 Features authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api_test

import (
	"fmt"
	"net/http"

	"github.com/albertoleal/features/engine"
	"github.com/apihub/apihub/requests"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateFeature(c *C) {
	featureKey := "login_via_email"
	defer func() {
		ffk := engine.FeatureFlagKey{Key: featureKey}
		s.ng.DeleteFeatureFlag(ffk)
	}()

	headers, code, body, _ := httpClient.MakeRequest(requests.Args{
		AcceptableCode: http.StatusCreated,
		Method:         "POST",
		Path:           "/features",
		Body:           fmt.Sprintf(`{"key": "%s", "percentage": 20, "enabled": true}`, featureKey),
	})

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, "{\"enabled\":true,\"key\":\"login_via_email\",\"percentage\":20}\n")
}

func (s *S) TestCreateFeatureMissingRequiredFields(c *C) {
	headers, code, body, _ := httpClient.MakeRequest(requests.Args{
		AcceptableCode: http.StatusBadRequest,
		Method:         "POST",
		Path:           "/features",
		Body:           `{"percentage": 20, "enabled": true}`,
	})

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, "{\"error\":\"bad_request\",\"error_description\":\"Key cannot be empty.\"}")
}

func (s *S) TestCreateFeatureInvalidJSON(c *C) {
	headers, code, body, _ := httpClient.MakeRequest(requests.Args{
		AcceptableCode: http.StatusBadRequest,
		Method:         "POST",
		Path:           "/features",
		Body:           `{"percentage": 2: "enabled": true}`,
	})

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, "{\"error\":\"bad_request\",\"error_description\":\"invalid character ':' after object key:value pair\"}")
}

func (s *S) TestCreateFeatureWithExistingKey(c *C) {
	featureKey := "login_via_email"
	feature := engine.FeatureFlag{
		Key:     featureKey,
		Enabled: false,
	}
	s.ng.UpsertFeatureFlag(feature)

	defer func() {
		ffk := engine.FeatureFlagKey{Key: featureKey}
		s.ng.DeleteFeatureFlag(ffk)
	}()

	headers, code, body, _ := httpClient.MakeRequest(requests.Args{
		AcceptableCode: http.StatusBadRequest,
		Method:         "POST",
		Path:           "/features",
		Body:           fmt.Sprintf(`{"key": "%s", "percentage": 20, "enabled": true}`, featureKey),
	})

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, "{\"error\":\"bad_request\",\"error_description\":\"There's another feature for the same key value.\"}\n")
}