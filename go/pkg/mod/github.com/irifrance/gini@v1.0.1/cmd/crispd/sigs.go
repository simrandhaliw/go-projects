// Copyright 2016 The Gini Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var sigs = make(chan os.Signal, 1)

func init() {
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	go func() {
		for {
			sig := <-sigs
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				fmt.Printf("\ninterrupted\n")
				if CrispD != nil {
					CrispD.Shutdown()
				}
				os.Exit(1)
			case syscall.SIGUSR1:
				fmt.Printf("\ncause SIGUSR1\n")
			}
		}
	}()
}
