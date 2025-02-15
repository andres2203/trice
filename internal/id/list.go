// Copyright 2020 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

// Package id List is responsible for id List managing
package id

// List management

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/rokath/trice/pkg/msg"
)

// ScZero does replace all ID's in source tree with 0
func ScZero(SrcZ string, cmd *flag.FlagSet) error {
	if SrcZ == "" {
		cmd.PrintDefaults()
		return errors.New("no source tree root specified")
	}
	ZeroSourceTreeIds(SrcZ, !DryRun)
	return nil
}

// SubCmdReNewList renews the trice id list parsing the source tree without changing any source file.
// It creates a new FnJSON and tries to add id:tf pairs from the source tree.
// If equal tf are found with different ids they are all added.
// If the same id is found with different tf only one is added. The others are reported as warning.
// If any TRICE* is found without Id(n) or with Id(0) it is ignored.
// SubCmdUpdate needs to know which IDs are used in the source tree to reliable add new IDs.
func SubCmdReNewList() (err error) {
	lu := make(TriceIDLookUp)
	// Do not perform lu.AddFmtCount() here.
	return updateList(lu)
}

// SubCmdRefreshList refreshes the trice id list parsing the source tree without changing any source file.
// It only reads FnJSON and tries to add id:tf pairs from the source tree.
// If equal tf are found with different ids they are all added.
// If the same id is found with different tf only one is added. The others are reported as warning.
// If any TRICE* is found without Id(n) or with Id(0) it is ignored.
// SubCmdUpdate needs to know which IDs are used in the source tree to reliable add new IDs.
func SubCmdRefreshList() (err error) {
	lu := NewLut(FnJSON)
	// Do not perform lu.AddFmtCount() here.
	return updateList(lu)
}

func refreshListAdapter(root string, lu TriceIDLookUp, tflu TriceFmtLookUp, _ *bool) {
	refreshList(root, lu, tflu)
}

func updateList(lu TriceIDLookUp) error {
	tflu := lu.reverse()

	// keep a copy
	lu0 := make(TriceIDLookUp)
	for k, v := range lu {
		lu0[k] = v
	}
	var listModified bool
	walkSrcs(refreshListAdapter, lu, tflu, &listModified)

	// listModified does not help here, because it indicates that some sources are updated and therefore the list needs an update too.
	// But here we are only scanning the source tree, so if there would be some changes they are not relevant because sources are not changed here.
	// And if all
	eq := reflect.DeepEqual(lu0, lu)

	if Verbose {
		fmt.Println(len(lu0), " -> ", len(lu), "ID's in List", FnJSON)
	}
	if !eq && !DryRun {
		msg.FatalOnErr(lu.toFile(FnJSON))
	}

	return nil // SubCmdUpdate() // to do
}

// SubCmdUpdate is subcommand update
func SubCmdUpdate() error {
	lu := NewLut(FnJSON)
	tflu := lu.reverse()
	var listModified bool
	walkSrcs(IDsUpdate, lu, tflu, &listModified)
	if Verbose {
		fmt.Println(len(lu), "ID's in List", FnJSON, "listModified=", listModified)
	}
	if listModified && !DryRun {
		msg.FatalOnErr(lu.toFile(FnJSON))
	}
	return nil
}

func walkSrcs(f func(root string, lu TriceIDLookUp, tflu TriceFmtLookUp, pListModified *bool), lu TriceIDLookUp, tflu TriceFmtLookUp, pListModified *bool) {
	if 0 == len(Srcs) {
		Srcs = append(Srcs, "./") // default value
	}
	for i := range Srcs {
		s := Srcs[i]
		srcU := ConditionalFilePath(s)
		if _, err := os.Stat(srcU); err == nil { // path exists
			f(srcU, lu, tflu, pListModified)
		} else if os.IsNotExist(err) { // path does *not* exist
			fmt.Println(s, " -> ", srcU, "does not exist!")
		} else {
			fmt.Println(s, "Schrodinger: file may or may not exist. See err for details.")
			// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
			// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
		}
	}
}
