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
	"io"
	"io/ioutil"

	"github.com/hashicorp/consul/consul"
	"fmt"
)

type PortConfig struct {
	DNS     int // DNS Query interface
	HTTP    int // HTTP API
	HTTPS   int // HTTPS API
	RPC     int // CLI RPC
	SerfLan int // LAN gossip (Client + Server)
	SerfWan int // WAN gossip (Server onlyg)
	Server  int // Server internal RPC
}

func (p *PortConfig) ApplyIndex(index int) {
	p.DNS += index
	p.HTTP += index
	p.HTTPS += index
	p.RPC += index
	p.SerfLan += index
	p.SerfWan += index
	p.Server += index
}

func (p PortConfig) String() string {
	ports := `
	DNS: %d
	HTTP: %d
	HTTPS: %d
	RPC: %d
	SerfLan: %d
	SerfWan: %d
	`
	return fmt.Sprintf(ports, p.DNS, p.HTTP, p.HTTPS, p.RPC, p.SerfLan, p.SerfWan)
}

// the context is a stripped down version of configuration for the Consul
type Context struct {
	// feature http
	EnableHTTP bool
	// feature dns
	EnableDNS bool
	// whether this node is the bootstrap node
	Bootstrap bool
	// a list of members to connect to
	Members []string
	// our node name
	NodeName string
	// enable debug
	EnableDebug bool
	// the loglevel
	LogOutput io.Writer
	// the key used to encrypt the traffic
	EncryptKey string
	// the datacenter
	Datacenter string
	// the datadir
	DataDir string
	// the agent address
	ClientAddress string
	// the address to bind
	BindAddress string
	// the address to advertise
	BindAdvertised string
	// the port configuration for the above
	PortsConfig PortConfig
}

func DefaultContext() *Context {
	return &Context{
		EnableHTTP:    true,
		EnableDNS:     false,
		Members:       make([]string, 0),
		EnableDebug:   false,
		LogOutput:     ioutil.Discard,
		Datacenter:    "dc1",
		ClientAddress: "0.0.0.0",
		BindAddress:   "0.0.0.0",
		PortsConfig: PortConfig{
			DNS:     8600,
			HTTP:    8500,
			HTTPS:   8501,
			RPC:     8400,
			SerfLan: consul.DefaultLANSerfPort,
			SerfWan: consul.DefaultWANSerfPort,
			Server:  8300,
		},
	}
}
