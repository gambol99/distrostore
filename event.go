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

import "fmt"

// an interface used for the join and leaving of nodes in the cluster
type NodeAPIEvent struct {
	// a link to the node info
	Node *Node
	// the status i.e. left, joined
	Status string
}

func (k NodeAPIEvent) String() string {
	return fmt.Sprintf("node: %s, status: %s", k.Node.ID, k.Status)
}

type KeyAPIEvent struct {
	// the key for this value
	Key string
	// the type of update (set, change, delete)
	Status string
}

func (k KeyAPIEvent) String() string {
	return fmt.Sprintf("key: %s, status: %s", k.Key, k.Status)
}
