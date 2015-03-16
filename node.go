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

// a structure for defining a node
type Node struct {
	// a unique id for the host
	ID string
	// the ip address of the node
	Address string
	// the port the node is running on
	Port int
}

func (n Node) String() string {
	return fmt.Sprintf("id: %s, node: %s:%d", n.ID, n.Address, n.Port)
}
