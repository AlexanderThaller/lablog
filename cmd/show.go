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

import "github.com/spf13/cobra"

var flagShowArchive bool

func init() {
	cmdShow.PersistentFlags().BoolVarP(&flagShowArchive, "archive", "a",
		false, "Determines if entries from the archive will be shown. (default is false)")

	cmdShow.AddCommand(cmdShowTodos)

	RootCmd.AddCommand(cmdShow)
}

var cmdShow = &cobra.Command{
	Use:   "show [command]",
	Short: "Show current projects and entries",
	Long:  `Show a list of currently available projects or their entries like notes, todos, tracks, etc., see help for all options.`,
	Run:   runCmdShow,
}

func runCmdShow(cmd *cobra.Command, args []string) {
	cmd.Help()
}
