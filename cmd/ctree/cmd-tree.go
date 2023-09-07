package main

import (
	"fmt"

	"github.com/teonet-go/teonet/cmd/teonet/menu"
)

// Create CmdTree commands
func (cli *Treecli) newCmdTree() menu.Item {
	return CmdTree{TreecliCommand: TreecliCommand{cli}}
}

// CmdTree connect to peer command ----------------------------------------
type CmdTree struct {
	TreecliCommand
}

func (c CmdTree) Name() string  { return cmdTree }
func (c CmdTree) Usage() string { return "[flag] <name>" }
func (c CmdTree) Help() string  { return "select or create tree ('choose' flag runs by default)" }
func (c CmdTree) Exec(line string) (err error) {
	var list, save, new, choose bool
	flags := c.NewFlagSet(c.Name(), c.Usage(), c.Help())
	flags.BoolVar(&new, "new", list, "create new tree")
	flags.BoolVar(&list, "list", list, "show list of trees")
	flags.BoolVar(&choose, "choose", list, "choose tree to use")
	flags.BoolVar(&save, "save", list, "save list of trees")
	err = flags.Parse(c.menu.SplitSpace(line))
	if err != nil {
		return
	}
	args := flags.Args()

	switch {
	// Check help
	case len(args) > 0 && args[0] == cmdHelp:
		flags.Usage()
		return

	// Check length of arguments
	case len(args) == 0:
		flags.Usage()
		err = ErrWrongNumArguments
		return

	// Check -new flag
	case new:
		fmt.Printf("new tree `%s` created\n", args[0])
		return

		// Check -list flag
		// case list:
		// 	aliases := c.alias.list()
		// 	for i := range aliases {
		// 		fmt.Printf("%s\n", aliases[i])
		// 	}
		// 	return

		// Check -save flag
		// case save:
		// 	aliases := c.alias.list()
		// 	c.batch.Save(aliasBatchFile, CmdTree, aliases)
		// 	return

	}

	// Add alias
	// c.alias.add(args[0], args[1])

	return
}
func (c CmdTree) Compliter() (cmpl []menu.Compliter) {
	return c.menu.MakeCompliterFromString([]string{"-list", "-save"})
}
