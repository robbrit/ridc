package main

import (
  "os"
  "os/signal"
  "syscall"
  "runtime/pprof"
  "./riddick"
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

  database := riddick.CreateDatabase()

  go riddick.StartWebInterface(database)
  go riddick.StartSocketInterface(database, "3001")

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
