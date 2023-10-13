package main

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"os"
	"syscall"
)

func killPipe() error {
	pid := os.Getpid()
	pgid, err := syscall.Getpgid(pid)
	if err != nil {
		fmt.Println("Failed to get PGID:", err)
		return err
	}

	// Get the list of all PIDs on the system using gopsutil.
	allPids, err := process.Pids()
	if err != nil {
		fmt.Println("Failed to get PIDs:", err)
		return err
	}
	for _, otherPid := range allPids {
		if int(otherPid) == pid {
			continue
		}
		otherPgid, err := syscall.Getpgid(int(otherPid))
		if err != nil {
			continue
		}

		// Check if the other process is in the same group as the current process.
		if otherPgid == pgid {
			err = syscall.Kill(int(otherPid), syscall.SIGINT)
			if err != nil {
				fmt.Println("Error sending signal to process", err)
				return err
			}
		}
	}
	return nil
}
