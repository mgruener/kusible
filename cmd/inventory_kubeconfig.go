/*
Copyright © 2019 Michael Gruener

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

package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newInventoryKubeconfigCmd(c *Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "kubeconfig [entry name]",
		Short:                 "Get the kubeconfig for a given inventory entry",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  c.wrap(runInventoryKubeconfig),
	}
	addInventoryFlags(cmd)

	return cmd
}

func runInventoryKubeconfig(c *Cli, cmd *cobra.Command, args []string) error {
	name := args[0]

	inv, err := getInventoryWithKubeconfig(c)
	if err != nil {
		return err
	}

	entry, ok := inv.Entries()[name]
	if !ok {
		log.WithFields(log.Fields{
			"entry": name,
		}).Error("Entry does not exist")
		return err
	}

	kubeconfig, err := entry.Kubeconfig().Yaml()
	if err != nil {
		log.WithFields(log.Fields{
			"entry": name,
			"error": err.Error(),
		}).Error("Failed to get kubeconfig")
		return err
	}
	fmt.Println(string(kubeconfig))
	return nil
}