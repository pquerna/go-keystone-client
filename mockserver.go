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
	"fmt"
	"net"
	"net/http"
)

func mockHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside handler")
	fmt.Println("Inside handler")
	fmt.Fprintf(w, "Hello world from my Go program!")
}

type MockHTTPServer struct {
	server   http.Server
	listener *net.TCPListener
}

func (mts *MockHTTPServer) ListenAndServe() error {

	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0, })

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
	mux.Handle("/", http.HandlerFunc(mockHandler))
	return &MockHTTPServer{
		server: http.Server{
			Handler: mux,
		},
	}
}
