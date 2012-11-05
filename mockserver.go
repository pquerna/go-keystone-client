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
	"bufio"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
)

func mockHandler(sourceDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.Method + r.RequestURI

		filename = strings.Replace(filename, "/", "_", -1)
		filename = strings.Replace(filename, ".", "_", -1)
		filename = strings.Replace(filename, "-", "_", -1)
		filename = filename + ".asis"

		filename = path.Join(sourceDir, filename)

		f, err := os.Open(filename)

		if err != nil {
			http.Error(w, "Invalid path: "+err.Error(), 599)
			return
		}

		defer f.Close()

		br := bufio.NewReader(f)

		resp, err := http.ReadResponse(br, r)
		if err != nil {
			http.Error(w, "Invalid HTTP Response in "+filename+": "+err.Error(), 599)
			return
		}
		rbody, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			http.Error(w, "Problem reading http response "+filename+": "+err.Error(), 599)
			return
		}

		for k := range resp.Header {
			for _, v := range resp.Header[k] {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(resp.StatusCode)

		w.Write(rbody)
	}
}

type MockHTTPServer struct {
	server   http.Server
	listener *net.TCPListener
}

func (mts *MockHTTPServer) ListenAndServe() error {

	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0})

	mts.listener = l

	if err != nil {
		return err
	}

	go mts.server.Serve(l)

	return nil
}

func (mts *MockHTTPServer) URL() string {
	return "http://" + mts.listener.Addr().String() + "/"
}

func (mts *MockHTTPServer) Close() {
	mts.listener.Close()
}

func NewMockHTTPServer(sourceDir string) *MockHTTPServer {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(mockHandler(sourceDir)))
	return &MockHTTPServer{
		server: http.Server{
			Handler: mux,
		},
	}
}
