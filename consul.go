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
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/command/agent"

)

const (
	DEFAULT_WAIT_TIME = (time.Duration(120) * time.Second)
)

type ConsulDistroStore struct {
	sync.RWMutex
	// the consul agent
	agent *agent.Agent
	// the cluster context config
	context *Context
	// the configuration for consul
	config *agent.Config
	// the client to the consul service
	client *api.Client
	// the consul http interface
	http_api []*agent.HTTPServer
	// the dns api
	dns_api []*agent.DNSServer
	// a map of those listening to key events
	key_listeners map[chan *KeyAPIEvent]bool
	// a map of those listening to node events
	node_listeners map[chan *NodeAPIEvent]bool
}

// Created a new node in the cluster
//  cfg:          configuration used to pass to the provider
func NewConsulDistributedStore(cfg *Context) (DistroStore, error) {
	var err error
	service := new(ConsulDistroStore)
	service.context = cfg
	service.key_listeners = make(map[chan *KeyAPIEvent]bool, 0)
	service.node_listeners = make(map[chan *NodeAPIEvent]bool, 0)

	// step: create the agent for the service
	if service.agent, err = service.createConsulAgent(cfg); err != nil {
		return nil, err
	}

	// step: create the client
	if service.client, err = service.createConsulClient(cfg); err != nil {
		return nil, err
	}

	return service, nil
}

func (r *ConsulDistroStore) Config() *Context {
	return r.context
}

func (r *ConsulDistroStore) createConsulAgent(cfg *Context) (*agent.Agent, error) {
	var err error
	var scadaList net.Listener
	var output io.Writer
	// step: set the logging
	output = os.Stderr
	if cfg.LogLevel == "NONE" {
		output = ioutil.Discard
	}

	// step: parse the context and fill in a config
	if r.config, err = r.parseContext(cfg); err != nil {
		return nil, nil
	}
	// step: create the actual agent
	service, err := agent.Create(r.config, output)
	if err != nil {
		return nil, err
	}

	// step: join other members
	if cfg.Members != nil && len(cfg.Members) > 0 {
		_, err := service.JoinLAN(cfg.Members)
		if err != nil {
			return service, nil
		}
	}

	// step: wait for the service to be available

	if cfg.EnableHTTP {
		r.http_api, err = agent.NewHTTPServers(service, r.config, scadaList, output)
		if err != nil {
			return nil, err
		}
	}

	service.StartSync()
	return service, nil
}

func (r *ConsulDistroStore) createConsulClient(cfg *Context) (*api.Client, error) {
	address := r.config.ClientAddr
	if address == "0.0.0.0" {
		address = "127.0.0.1"
	}
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", address, r.config.Ports.HTTP)
	config.Datacenter = r.config.Datacenter
	cli, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// Add in new member to the cluster
//  member: 	the endpoint address i.e. the <IPADDRESS>:<PORT>
func (r *ConsulDistroStore) Join(member string) error {
	if !isEndpoint(member) {
		return ErrInvalidMemberAddress
	}
	_, err := r.agent.JoinLAN([]string{member})
	if err != nil {
		return err
	}
	return nil
}

func (r *ConsulDistroStore) Close() error {
	if err := r.agent.Leave(); err != nil {
		return err
	}
	if err := r.agent.Shutdown(); err != nil {
		return err
	}
	return nil
}

// Retrieve a list of node presently in the cluster
func (r *ConsulDistroStore) Nodes() ([]*Node, error) {
	list := make([]*Node, 0)
	members := r.agent.LANMembers()
	for _, member := range members {
		node := &Node{
			ID:      member.Name,
			Address: member.Addr.String(),
			Port:    int(member.Port),
		}
		list = append(list, node)
	}
	return list, nil
}

// Get the value from the consul key/value store
//  key:		the key we are interested in
func (r *ConsulDistroStore) Get(key string) (string, bool, error) {
	pair, _, err := r.kv().Get(key, nil)
	if err != nil {
		return "", false, err
	}
	if pair == nil {
		return "", false, nil
	}
	return string(pair.Value), true, nil
}

// Check to see if a key exists in the store
// key:		the key you are looking for
func (r *ConsulDistroStore) Exists(key string) (bool, error) {
	_, found, err := r.Get(key)
	return found, err
}

// Set a key/pair in the consul k/v store
//  key: 	the key you wish to set
//  data:	the value of the key
func (r *ConsulDistroStore) Set(key, data string) error {
	keypair := &api.KVPair{
		Key:   key,
		Value: []byte(data),
	}
	if _, err := r.kv().Put(keypair, nil); err != nil {
		return err
	}
	return nil
}

// Add a listener for node membership events
//  channel: 	the channel to pass the events upon
func (r *ConsulDistroStore) AddNodeListener(channel chan *NodeAPIEvent) {
	if _, found := r.node_listeners[channel]; !found {
		r.Lock()
		defer r.Unlock()
		r.node_listeners[channel] = true
	}
}

// Add a listener for key events
//  channel: 	the channel to pass the events upon
func (r *ConsulDistroStore) AddKeyListener(channel chan *KeyAPIEvent) {
	if _, found := r.key_listeners[channel]; !found {
		r.Lock()
		defer r.Unlock()
		r.key_listeners[channel] = true
	}
}

func (r *ConsulDistroStore) waitIndex() (uint64, error) {
	if _, meta, err := r.kv().Get("/", nil); err != nil {
		return 0, err
	} else {
		return meta.LastIndex, nil
	}
}

func (r *ConsulDistroStore) watchKeys() error {
	// the wait index for consul
	var wait_index uint64

	for {
		// get the wait index if not set
		if idx, err := r.waitIndex(); err != nil {

		} else {
			wait_index = idx
		}

		// wait for any changes on in the keys
		_, meta, err := r.kv().List("/", &api.QueryOptions{WaitIndex: wait_index,
			WaitTime: DEFAULT_WAIT_TIME})
		if err != nil {
			// we need to backoff and wait for a bit

			continue
		}
		// update the index
		wait_index = meta.LastIndex
	}
}

func (r *ConsulDistroStore) kv() *api.KV {
	return r.client.KV()
}

func (r *ConsulDistroStore) parseContext(cfg *Context) (*agent.Config, error) {
	config := agent.DefaultConfig()
	config.Server = true
	config.DataDir = cfg.DataDir
	config.EnableDebug = cfg.EnableDebug
	config.Bootstrap = cfg.Bootstrap
	config.EncryptKey = cfg.EncryptKey
	config.NodeName = cfg.NodeName
	config.LogLevel = "NONE"
	config.Datacenter = cfg.Datacenter
	config.VerifyIncoming = false
	config.VerifyOutgoing = false
	config.BindAddr = cfg.BindAddress
	config.AdvertiseAddr = cfg.BindAdvertised
	config.ClientAddr = cfg.ClientAddress
	config.Ports = agent.PortConfig{
		DNS:     cfg.PortsConfig.DNS,
		HTTP:    cfg.PortsConfig.HTTP,
		HTTPS:   cfg.PortsConfig.HTTPS,
		RPC:     cfg.PortsConfig.RPC,
		SerfLan: cfg.PortsConfig.SerfLan,
		SerfWan: cfg.PortsConfig.SerfWan,
		Server:  cfg.PortsConfig.Server,
	}
	config.StartJoin = cfg.Members
	//config.RetryInterval = time.Duration(2) * time.Second
	//config.RetryMaxAttempts = 3
	return config, nil
}
