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
	"net"
	"testing"
	"time"

	"github.com/hashicorp/consul/consul"
)

var nextPort = 15000

func getPort() int {
	p := nextPort
	nextPort++
	return p
}

func testServerConfig(t *testing.T, NodeName string) (string, *consul.Config) {
	dir := tmpDir(t)
	config := consul.DefaultConfig()
	config.NodeName = NodeName
	config.Bootstrap = true
	config.Datacenter = "dc1"
	config.DataDir = dir
	config.RPCAddr = &net.TCPAddr{IP: []byte{127, 0, 0, 1}, Port: getPort()}
	config.SerfLANConfig.MemberlistConfig.BindAddr = "127.0.0.1"
	config.SerfLANConfig.MemberlistConfig.BindPort = getPort()
	config.SerfLANConfig.MemberlistConfig.SuspicionMult = 2
	config.SerfLANConfig.MemberlistConfig.ProbeTimeout = 50 * time.Millisecond
	config.SerfLANConfig.MemberlistConfig.ProbeInterval = 100 * time.Millisecond
	config.SerfLANConfig.MemberlistConfig.GossipInterval = 100 * time.Millisecond
	config.SerfWANConfig.MemberlistConfig.BindAddr = "127.0.0.1"
	config.SerfWANConfig.MemberlistConfig.BindPort = getPort()
	config.SerfWANConfig.MemberlistConfig.SuspicionMult = 2
	config.SerfWANConfig.MemberlistConfig.ProbeTimeout = 50 * time.Millisecond
	config.SerfWANConfig.MemberlistConfig.ProbeInterval = 100 * time.Millisecond
	config.SerfWANConfig.MemberlistConfig.GossipInterval = 100 * time.Millisecond
	config.RaftConfig.LeaderLeaseTimeout = 20 * time.Millisecond
	config.RaftConfig.HeartbeatTimeout = 40 * time.Millisecond
	config.RaftConfig.ElectionTimeout = 40 * time.Millisecond
	config.ReconcileInterval = 100 * time.Millisecond
	return dir, config
}

func testServer(t *testing.T) (string, *consul.Server) {
	return testServerDC(t, "dc1")
}

func testServerDC(t *testing.T, dc string) (string, *consul.Server) {
	return testServerDCBootstrap(t, dc, true)
}

func testServerDCBootstrap(t *testing.T, dc string, bootstrap bool) (string, *consul.Server) {
	name := fmt.Sprintf("Node %d", getPort())
	dir, config := testServerConfig(t, name)
	config.Datacenter = dc
	config.Bootstrap = bootstrap
	server, err := consul.NewServer(config)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	return dir, server
}
