/* Acksin Autotune - Kernel Autotuning
 * Copyright (C) 2015 Acksin <hey@acksin.com>
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package stats

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/acksin/autotune/stats/fd"
	"github.com/acksin/autotune/stats/memory"
	"github.com/acksin/procfs"
)

type Stats struct {
	CmdName string

	Duration int
	Every    int
}

func (n *Stats) Synopsis() string {
	return ""
}

func (n *Stats) Help() string {
	return ""
}

func (n *Stats) Run(args []string) int {
	flags := flag.NewFlagSet(n.CmdName, flag.ContinueOnError)
	flags.IntVar(&n.Every, "every", -1, "Run stats [every] seconds and give average.")
	flags.IntVar(&n.Duration, "duration", -1, "Run command for [duration] seconds.")

	if err := flags.Parse(args); err != nil {
		return -1
	}

	s := Response{}

	s.System.Memory = memory.New()

	for _, proc := range n.processes() {
		exe, err := proc.Executable()
		if err != nil || exe == "" {
			status, _ := proc.NewStatus()
			exe = status.Name
		}

		p := Process{
			Exe:    exe,
			PID:    proc.PID,
			Memory: memory.NewProcess(proc),
			FD:     fd.NewProcess(proc),
		}

		s.Processes = append(s.Processes, p)
	}

	printJson(s)

	return 0
}

func (n *Stats) processes() procfs.Procs {
	fs, err := procfs.NewFS(procfs.DefaultMountPoint)
	if err != nil {
		return nil
	}

	procs, err := fs.AllProcs()
	if err != nil {
		return nil
	}

	return procs
}

func printJson(s interface{}) error {
	js, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(js))

	return nil
}

func New(cmdName string) *Stats {
	return &Stats{
		CmdName: cmdName,
	}
}
