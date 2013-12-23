/*
Copyright (C) 2013 Rob Britton

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
  "os"
  "os/signal"
  "syscall"
  "runtime/pprof"
  "./ridc"
  "flag"
  "log"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpuprofile to file")

func main() {
  flag.Parse()

  if *cpuprofile != "" {
    f, err := os.Create(*cpuprofile)
    if err != nil {
      log.Printf("Could not open profile file: %s", err)
    } else {
      log.Printf("Logging CPU profile to %s", *cpuprofile)
      pprof.StartCPUProfile(f)
    }
  }

  database := ridc.CreateDatabase()

  go ridc.StartWebInterface(database)
  go ridc.StartSocketInterface(database, "3001")

  cleanup := func() {
    pprof.StopCPUProfile()
  }

  signalHandler := make(chan os.Signal, 1)
  signal.Notify(signalHandler, os.Interrupt)
  signal.Notify(signalHandler, syscall.SIGTERM)
  go func() {
    <-signalHandler
    log.Printf("Shutting down...")
    cleanup()
    os.Exit(1)
  }()

  select {}
}
