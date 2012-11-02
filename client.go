/**
 *  Copyright 2012 Paul Querna
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package keystone

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httputil"
)

var debugprint = true

type ClientOptions struct {
	Username  string
	Password  string
	APIKey    string
	BaseURL   string
	Version   string
	UserAgent string
}

var defaultOptions = ClientOptions{
	Version:   "v2.0",
	BaseURL:   "https://identity.api.rackspacecloud.com/",
	UserAgent: "keystone-client (golang; https://github.com/pquerna/go-keystone-client)",
}

type authenticateWithPassword struct {
	Auth struct {
		Credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"passwordCredentials"`
	} `json:"auth"`
}

type authenticateWithAPIKey struct {
	Auth struct {
		Credentials struct {
			Username string `json:"username"`
			APIKey   string `json:"apikey"`
		} `json:"RAX-KSKEY:apiKeyCredentials"`
	} `json:"auth"`
}

type ServiceCatalog struct {
	TenantId string
	/* TOOD: map of services */
}

type KeystoneClient struct {
	opts   ClientOptions
	client *http.Client
}

func (kc *KeystoneClient) baseURL() string {
	rv := kc.opts.BaseURL + kc.opts.Version
	return rv
}

func (kc *KeystoneClient) authReqBody() interface{} {
	if len(kc.opts.APIKey) > 0 {
		data := authenticateWithAPIKey{}

		data.Auth.Credentials.Username = kc.opts.Username
		data.Auth.Credentials.APIKey = kc.opts.APIKey
		return data

	} else if len(kc.opts.Password) > 0 {
		data := authenticateWithPassword{}

		data.Auth.Credentials.Username = kc.opts.Username
		data.Auth.Credentials.Password = kc.opts.Password
		return data
	}

	panic("opts must include APIKey or Password")
}

func (kc *KeystoneClient) runReq(req *http.Request) (*http.Response, error) {
	if debugprint {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			println(err.Error())
		}
		println("")
		print(string(dump))
	}
	return kc.client.Do(req)
}

func (kc *KeystoneClient) prepReq(method string, url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", kc.opts.UserAgent)
	req.ContentLength = int64(len(body))

	return req, nil
}

func (kc *KeystoneClient) ServiceCatalog() (*ServiceCatalog, error) {
	data := kc.authReqBody()
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := kc.prepReq("POST", kc.baseURL(), body)

	if err != nil {
		return nil, err
	}

	resp, err := kc.runReq(req)

	if debugprint {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			println(err.Error())
		}
		println("")
		print(string(dump))
	}

	return nil, nil
}

func NewKeystoneClient() *KeystoneClient {
	return &KeystoneClient{
		client: &http.Client{},
	}
}

func Dial(opts ClientOptions) (*KeystoneClient, error) {

	kc := NewKeystoneClient()

	if len(opts.Username) == 0 {
		return nil, errors.New("Username must be set on client options.")
	}

	kc.opts.Username = opts.Username

	if len(opts.Password) == 0 && len(opts.APIKey) == 0 {
		return nil, errors.New("APIKey or Password must be set.")
	}
	kc.opts.Password = opts.Password
	kc.opts.APIKey = opts.APIKey

	kc.opts.Version = defaultOptions.Version
	if len(opts.Version) != 0 {
		kc.opts.Version = opts.Version
	}

	kc.opts.BaseURL = defaultOptions.BaseURL
	if len(opts.BaseURL) != 0 {
		kc.opts.BaseURL = opts.BaseURL
	}

	kc.opts.UserAgent = defaultOptions.UserAgent
	if len(opts.UserAgent) != 0 {
		kc.opts.UserAgent = opts.UserAgent
	}

	return kc, nil
}
