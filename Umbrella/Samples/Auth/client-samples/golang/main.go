/*
Copyright (c) 2021 Cisco and/or its affiliates.
This software is licensed to you under the terms of the Cisco Sample
Code License, Version 1.1 (the "License"). You may obtain a copy of the
License at

https://developer.cisco.com/docs/licenses

All use of the material herein must be in accordance with the terms of
the License. All rights not expressly granted by the License are
reserved. Unless required by applicable law or agreed to separately in
writing, software distributed under the License is distributed on an "AS
IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
or implied.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	envKey_ApiKey    = "API_KEY"
	envKey_ApiSecret = "API_SECRET"

	tokenUri          = "https://api.umbrella.com/auth/v2/token"
	reportingUri      = "https://api.umbrella.com/reports/v2/summary?from=-5days&to=now"
	apiReqCallWaitSec = 1
	apiRepeatCount    = 3
)

func umbrellaReports(httpClient *http.Client) {
	startTime := time.Now()
	res, err := httpClient.Get(reportingUri)
	if err != nil {
		fmt.Printf("Error calling API %s, %s\n", reportingUri, err.Error())
	} else {
		defer res.Body.Close()
		b, _ := ioutil.ReadAll(res.Body)
		latency := (time.Since(startTime).Nanoseconds()) / 1000000
		fmt.Printf("Code: %d: RspTime: %v(ms): %s\n%s\n", res.StatusCode, latency, reportingUri, b)
	}
}

func main() {
	// Parse envars for client credentials and orgId
	clientId, clientIdFound := os.LookupEnv(envKey_ApiKey)
	if !clientIdFound {
		fmt.Printf("Mandatory environment variable %s not set, ", envKey_ApiKey)
		os.Exit(1)
	}
	clientSecret, clientSecretFound := os.LookupEnv(envKey_ApiSecret)
	if !clientSecretFound {
		fmt.Printf("Mandatory environment variable %s not set, ", envKey_ApiSecret)
		os.Exit(1)
	}

	// Create oauth2 client
	config := &clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     tokenUri,
	}
	httpClient := config.Client(context.Background())
	if httpClient == nil {
		fmt.Printf("Error creating Oauth2 http client for %s ", config.TokenURL)
		return
	}

	// Call the api multiple times with client credentials.
	// The access token is acquired and refreshed upon expiry automatically.
	for i := 0; i < apiRepeatCount; i++ {
		// Call reporting api
		umbrellaReports(httpClient)

		// Sleep for some time
		time.Sleep(apiReqCallWaitSec * time.Second)
	}
}
