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
	"errors"
)

var (
	ErrInvalidConfig = errors.New("Invalid configuration supplied")
	// an invalid member / endpoint
	ErrInvalidMemberAddress = errors.New("Invalid members / endpoint address")
)

type DistroStore interface {
	Config() *Context
	// shutdown and release resources
	Close() error
	// join a new member to the cluster
	Join(member string) error
	// get a list of the nodes in the cluster
	Nodes() ([]*Node, error)
	// check if a key exists
	Exists(key string) (bool, error)
	// set a value in the store
	Set(key string, data string) error
	// get the value from the store
	Get(key string) (string, bool, error)
	// add a node listener for the cluster
	AddNodeListener(channel chan *NodeAPIEvent)
	// watch for changes in the store
	AddKeyListener(channel chan *KeyAPIEvent)
}

func New(cfg *Context) (DistroStore, error) {
	if cfg == nil {
		return nil, errors.New("You have not specified any configuration")
	}
	return NewConsulDistributedStore(cfg)
}
