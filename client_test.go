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
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestInvalidParams(c *C) {

	opt := ClientOptions{
		Username: "foo",
	}

	client, err := Dial(opt)

	c.Assert(client, IsNil)
	c.Assert(err, ErrorMatches, "APIKey or Password(.+)")

	opt = ClientOptions{
		APIKey: "XXXXX",
	}

	client, err = Dial(opt)

	c.Assert(client, IsNil)
	c.Assert(err, ErrorMatches, "Username must be(.+)")

}

func (s *MySuite) TestServiceCatalog(c *C) {
	opt := ClientOptions{
		Username: "foo",
		APIKey:   "XXXXX",
	}

	client, err := Dial(opt)

	sc, err := client.ServiceCatalog()

	c.Assert(err, IsNil)

	c.Assert(sc, IsNil)

}

func (s *MySuite) TestServiceCatalog2(c *C) {
	opt := ClientOptions{
		Username: "foo",
		Password: "XXXXX",
	}

	client, err := Dial(opt)

	sc, err := client.ServiceCatalog()

	c.Assert(err, IsNil)

	c.Assert(sc, IsNil)

}
