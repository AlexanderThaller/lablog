// Copyright © 2016 Alexander Thaller <alexander@thaller.ws>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	"github.com/AlexanderThaller/lablog/src/helper"
	"github.com/juju/errgo"

	"github.com/spf13/cobra"
)

func init() {
	cmdShow.AddCommand(cmdShowProjects)
}

var cmdShowProjects = &cobra.Command{
	Use:   "projects",
	Short: "Show current projects",
	Long:  `Show all projects which currently have any type of entry.`,
	RunE:  runCmdShowProjects,
}

func runCmdShowProjects(cmd *cobra.Command, args []string) error {
	store, err := helper.DefaultStore(flagDataDir)
	if err != nil {
		return errgo.Notef(err, "can not get data store")
	}

	projects, err := helper.ProjectNamesFromArgs(store, args, flagShowArchive)
	if err != nil {
		return errgo.Notef(err, "can not get list of projects")
	}

	for _, project := range projects.List() {
		fmt.Println(project.Name)
	}

	return nil
}
