// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tree CLI applicatiom - command line interface application allow create and
// process any tree data. For example it may be decision tree.
package main

import "fmt"

const (
	appShort   = "ctree"
	appName    = "Tree CLI application"
	appVersion = "v0.0.1"
)

func main() {

	// Application logo
	logo(appName, appVersion)

	// Create Tree CLI interface
	cli, err := NewTreeCli(appShort)
	if err != nil {
		fmt.Println("can't create Tree CLI interface, err:", err)
	}

	// TODO: Run batch files
	cli.batch.run(aliasBatchFile)
	// cli.batch.run(connectBatchFile)

	// Run Teonet CLI commands menu
	fmt.Print(
		"\n",
		"Usage:	<command> [arguments]\n",
		"use help command to get commands list\n\n",
	)
	cli.menu.Run()
}

// logo prints teonet logo
func logo(appName, version string) { fmt.Println(logoString(appName, version)) }

// logo returns string with application logo
func logoString(appName, version string) string {
	return fmt.Sprint("" +
		"  ____ _____              \n" +
		" / ___|_   _| __ ___  ___  tm\n" +
		"| |     | || '__/ _ \\/ _ \\\n" +
		"| |___  | || | |  __/  __/\n" +
		" \\____| |_||_|  \\___|\\___|\n" +
		"\n" +
		appName + " ver " + version + "\n",
	)
}
