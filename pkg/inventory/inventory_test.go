// Copyright © 2019 Michael Gruener
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inventory

import (
	"fmt"
	"sort"
	"testing"

	"github.com/bedag/kusible/pkg/values"
	"github.com/go-test/deep"
	"gotest.tools/assert"
)

func basicInventoryTest(path string, filter string, skip bool, expected []string) (*Inventory, error) {
	ejsonSettings := values.EjsonSettings{
		PrivKey:     "",
		KeyDir:      "",
		SkipDecrypt: false,
	}

	inventory, err := NewInventory(path, ejsonSettings, skip)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory: %s", err)
	}

	result, err := inventory.EntryNames(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve inventory entries: %s", err)
	}

	// we only want to compare the elements and not
	// the order of the elements
	sort.Strings(expected)
	sort.Strings(result)

	if diff := deep.Equal(result, expected); diff != nil {
		return nil, fmt.Errorf("unexpected list of inventory entries: %#v", diff)
	}
	return inventory, nil
}

func TestInventoryEntriesFull(t *testing.T) {
	inventoryPath := "testdata/clusters_default.yaml"
	skipKubeconfig := true
	filter := ".*"
	expected := []string{
		"cluster-test-01-preflight",
		"cluster-dev-01",
		"cluster-test-01",
		"cluster-stage-01",
		"cluster-stage-02",
		"cluster-stage-03",
		"cluster-prod-01",
		"cluster-prod-02",
		"cluster-prod-03",
		"cluster-prod-04",
	}

	_, err := basicInventoryTest(inventoryPath, filter, skipKubeconfig, expected)
	assert.NilError(t, err)
}

func TestInventoryEntriesSingle(t *testing.T) {
	inventoryPath := "testdata/clusters_default.yaml"
	skipKubeconfig := true
	expected := []string{
		"cluster-dev-01",
	}
	filter := expected[0]

	_, err := basicInventoryTest(inventoryPath, filter, skipKubeconfig, expected)
	assert.NilError(t, err)
}

func TestInventoryLoadFile(t *testing.T) {
	inventoryPath := "testdata/clusters_file.yaml"
	skipKubeconfig := false
	filter := ".*"
	expected := []string{
		"cluster-test-01",
		"cluster-test-02",
		"cluster-test-03",
	}
	_, err := basicInventoryTest(inventoryPath, filter, skipKubeconfig, expected)
	assert.NilError(t, err)
}
