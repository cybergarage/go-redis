// Copyright (C) 2022 The go-redis Authors All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
 go-redisd is an example of implementing a compatible Redis server using go-mysql.
	NAME
	 go-redisd

	SYNOPSIS
	 go-redisd [OPTIONS]

	OPTIONS
	-v      : Enable verbose output.
	-p      : Enable profiling.

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
*/
package main

import (
	"flag"
	"log"

	"os"
	"os/signal"
	"syscall"

	"github.com/cybergarage/go-redis/examples/go-redisd/server"
)

const (
	programName = " go-redisd"
)

func main() {
	// isDebugEnabled := flag.Bool("debug", false, "enable debugging log output")
	// isProfileEnabled := flag.Bool("profile", false, "enable profiling server")
	flag.Parse()

	server := server.NewServer()
	err := server.Start()
	if err != nil {
		log.Printf("%s couldn't be started (%s)", programName, err.Error())
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM)

	exitCh := make(chan int)

	go func() {
		for {
			s := <-sigCh
			switch s {
			case syscall.SIGHUP:
				log.Printf("caught SIGHUP, restarting...")
				err = server.Restart()
				if err != nil {
					log.Printf("%s couldn't be restarted (%s)", programName, err.Error())
					os.Exit(1)
				}
			case syscall.SIGINT, syscall.SIGTERM:
				log.Printf("Caught %s, stopping...", s.String())
				err = server.Stop()
				if err != nil {
					log.Printf("%s couldn't be stopped (%s)", programName, err.Error())
					os.Exit(1)
				}
				exitCh <- 0
			}
		}
	}()

	code := <-exitCh

	os.Exit(code)
}
