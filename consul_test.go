/*
Copyright 2014 Rohith All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package distrostore

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func tmpDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "consul")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	return dir
}

var (
	current_index = 100
	test_server   DistroStore
	lock          sync.Once
)

func createTestServerBootstrap(t *testing.T) DistroStore {
	config := DefaultContext()
	config.Bootstrap = true
	config.BindAddress = "127.0.0.1"
	config.DataDir = tmpDir(t)
	config.LogOutput = os.Stdout
	config.PortsConfig.ApplyIndex(current_index)
	return createServer(config, t)
}

func createTestServerMembers(t *testing.T, members []string) DistroStore {
	config := DefaultContext()
	config.BindAddress = "127.0.0.1"
	config.Members = members
	config.DataDir = tmpDir(t)
	config.LogOutput = os.Stdout
	config.PortsConfig.ApplyIndex(current_index)
	return createServer(config, t)
}

func createServer(config *Context, t *testing.T) DistroStore {
	if server, err := New(config); err != nil {
		t.Fatalf("Unable to create the fake consul cluster, error: %s", err)
	} else {
		time.Sleep(time.Duration(3) * time.Second)
		return server
	}
	return nil
}

func createFixedService(t *testing.T) DistroStore {
	lock.Do(func() {
		cfg := DefaultContext()
		cfg.Datacenter = "dc1"
		cfg.LogOutput = os.Stdout
		cfg.NodeName = "test1"
		cfg.Bootstrap = true
		cfg.BindAddress = "127.0.0.1"
		cfg.BindAdvertised = "127.0.0.1"
		cfg.DataDir = tmpDir(t)
		cfg.EnableDebug = false
		test_server = createServer(cfg, t)
	})
	return test_server
}

func TestGet(t *testing.T) {
	server := createFixedService(t)
	err := server.Set("test", "hello")
	assert.Nil(t, err, "we should not recieve an error here %s", err)
	value, found, err := server.Get("test")
	assert.Nil(t, err, "we should not generate an error")
	assert.True(t, found, "the found boolean should be true")
	assert.NotEmpty(t, value, "the retrieve value for the key should not be empty")
	assert.Equal(t, "hello", value)
}

func TestSet(t *testing.T) {
	server := createFixedService(t)
	err := server.Set("test", "hello")
	assert.Nil(t, err, "we should not recieve an error here %s", err)
}

func TestNodes(t *testing.T) {
	server := createFixedService(t)
	members, err := server.Nodes()
	assert.Nil(t, err, "we should not recieve an error here %s", err)
	assert.NotNil(t, members, "the members list should not be nil")
	assert.Equal(t, 1, len(members), "the size of the members should be one")
}

func TestNodesList(t *testing.T) {
	server := createFixedService(t)
	members, err := server.Nodes()
	assert.Nil(t, err, "we should not recieve an error here %s", err)
	assert.NotNil(t, members, "the members list should not be nil")
	assert.Equal(t, 1, len(members), "the size of the members should be one")
	member := members[0]
	assert.Equal(t, "127.0.0.1", member.Address, "the member address is incorrect")
	assert.Equal(t, 8301, member.Port, "the member port is incorrect")
	assert.Equal(t, "test1", member.ID, "the member ID is incorrect")
}

func TestExists(t *testing.T) {
	server := createFixedService(t)
	key := "exist_flag"
	val := "125423613127632"
	found, err := server.Exists(key)
	assert.Nil(t, err, "we should not recieve an error: %s", err)
	assert.False(t, found, "the found flag should have been false")
	err = server.Set(key, val)
	assert.Nil(t, err, "we should not recieve an error: %s", err)
	found, err = server.Exists(key)
	assert.Nil(t, err, "we should not recieve an error: %s", err)
	assert.True(t, found, "the found flag should have been true")
}

func TestConfig(t *testing.T) {
	server := createFixedService(t)
	config := server.Config()
	assert.NotNil(t, config, "we have not recieved the cluster config")
}

func TestJoining(t *testing.T) {
	server := createFixedService(t)
	config := DefaultContext()
	config.BindAddress = "127.0.0.1"
	config.DataDir = tmpDir(t)
	config.NodeName = "secon"
	config.PortsConfig.ApplyIndex(current_index)
	endpoint := fmt.Sprintf("127.0.0.1:%d", server.Config().PortsConfig.SerfLan)
	config.Members = make([]string, 0)
	config.Members = append(config.Members, endpoint)
	secondary := createServer(config, t)
	nodes, err := server.Nodes()
	assert.Nil(t, err, "unable to get a list of the node, error: %s", err)
	assert.NotNil(t, nodes, "the list of node is nil")
	time.Sleep(time.Duration(2) * time.Second)
	assert.Equal(t, 2, len(nodes), "the nodes size should be two")
	secondary.Close()
}
