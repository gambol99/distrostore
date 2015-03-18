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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEndpoint(t *testing.T) {
	assert.Equal(t, false, isEndpoint("172.0.0.11:2222222"), "the method should have been false")
	assert.Equal(t, true, isEndpoint("172.0.0.11:222"), "the method should have been true")
}
