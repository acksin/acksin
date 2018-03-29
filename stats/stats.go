/*
 * Copyright (C) 2017 Acksin, LLC <hi@opszero.com>
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package stats

import (
	"encoding/json"

	"github.com/opszero/procfs"

	"github.com/opszero/opszero/stats/cloud"
	"github.com/opszero/opszero/stats/container"
	"github.com/opszero/opszero/stats/disk"
	"github.com/opszero/opszero/stats/io"
	"github.com/opszero/opszero/stats/kernel"
	"github.com/opszero/opszero/stats/memory"
	"github.com/opszero/opszero/stats/network"
)

// Stats contains both the system and process statistics.
type Stats struct {
	// System specific information
	System System
	// Container specific information
	Container *container.Container
	// Cloud specific information
	Cloud *cloud.Cloud
	// Processes are the process information of the system
	Processes []Process

	config *shared.Config
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

func (n *Stats) containsPid(pids []int, proc procfs.Proc) bool {
	for _, pid := range pids {
		if proc.PID == pid {
			return true
		}
	}

	return false
}

// UnmarshalJSON converts a byte blob into the Stats object
// representation.
func (n *Stats) UnmarshalJSON(d []byte) error {
	return json.Unmarshal(d, n)
}

// New returns stats of the machine with pids filtering for
// processes. If pids are empty then it returns all process stats.
func New(c *shared.Config) (s *Stats) {
	s = &Stats{}

	s.System.Memory = memory.New()
	s.System.Disk = disk.New()
	s.System.Network = network.New()
	s.System.Kernel = kernel.New()

	s.Container = container.New()
	s.Cloud = cloud.New(c)

	for _, proc := range s.processes() {
		exe, err := proc.Executable()
		if err != nil || exe == "" {
			status, _ := proc.NewStatus()
			exe = status.Name
		}

		p := Process{
			Exe:    exe,
			PID:    proc.PID,
			Memory: memory.NewProcess(proc),
			IO:     io.NewProcess(proc),
		}

		s.Processes = append(s.Processes, p)
	}

	return s
}
