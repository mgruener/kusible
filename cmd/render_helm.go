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

	helmutil "github.com/bedag/kusible/pkg/wrapper/helm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	helmcli "helm.sh/helm/v3/pkg/cli"
)

func newRenderHelmCmd(c *Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "helm [playbook]",
		Short:                 "Use helm to render manifests for an inventory entry",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  c.wrap(runRenderHelm),
	}
	addRenderFlags(cmd)

	return cmd
}

func runRenderHelm(c *Cli, cmd *cobra.Command, args []string) error {
	playbookFile := args[0]

	playbookSet, err := loadPlaybooks(c, playbookFile)
	if err != nil {
		return err
	}

	settings := helmcli.New()

	for name, playbook := range playbookSet {
		for _, play := range playbook.Config.Plays {
			for _, repo := range play.Repos {
				if err := helmutil.RepoAdd(repo.Name, repo.URL, settings); err != nil {
					log.WithFields(log.Fields{
						"play":  play.Name,
						"repo":  repo.Name,
						"entry": name,
						"error": err.Error(),
					}).Error("Failed to add helm repo for play.")
					return err
				}
			}
			manifests, err := helmutil.TemplatePlay(play, settings)
			if err != nil {
				log.WithFields(log.Fields{
					"play":  play.Name,
					"entry": name,
					"error": err.Error(),
				}).Error("Failed to render play manifests with helm.")
				return err
			}
			fmt.Printf(manifests)
		}
	}
	return nil
}