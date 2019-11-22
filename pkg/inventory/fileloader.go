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
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func NewKubeconfigFileLoaderFromParams(params map[string]string) *kubeconfigFileLoader {
	result := map[string]string{
		"decryptkey": os.Getenv("EJSON_PRIVKEY"),
		"path":       "kubeconfig",
	}

	for k, v := range params {
		result[strings.ToLower(k)] = v
	}

	return NewKubeconfigFileLoader(
		result["decryptkey"],
		result["path"])
}

func NewKubeconfigFileLoader(path string, decryptKey string) *kubeconfigFileLoader {
	return &kubeconfigFileLoader{Path: path, DecryptKey: decryptKey}
}

func (loader *kubeconfigFileLoader) Load() ([]byte, error) {
	_, err := os.Stat(loader.Path)
	if err != nil {
		return nil, err
	}

	mime, _, err := mimetype.DetectFile(loader.Path)
	if err != nil {
		return nil, err
	}

	var rawKubeconfig []byte
	switch mime {
	case "text/plain":
		rawKubeconfig, err = ioutil.ReadFile(loader.Path)
		if err != nil {
			return nil, err
		}
	case "application/x-7z-compressed":
		rawKubeconfig, err = extractSingleTar7ZipFile(loader.Path, loader.DecryptKey)
		if err != nil {
			return nil, err
		}
	case "application/octet-stream":
		rawKubeconfig, err = decryptOpensslSymmetricFile(loader.Path, loader.DecryptKey)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Unknown kubeconfig source file type: " + mime)
	}

	return rawKubeconfig, nil
}

func (loader *kubeconfigFileLoader) Type() string {
	return "file"
}

func (loader *kubeconfigFileLoader) Config() []byte {
	// TODO file loader config dump
	var result []byte
	return result
}
