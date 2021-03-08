package main

import (
	"github.com/tomasen/rollover"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// There are few things this example did:
// 1. call rollover.Wait() to kill and wait parent exit if there is any
// 2, show the command line arguments, so we can see is the child inherited parent's arguments
// 3. initiating rollover(restart) when signal HUP received
//    or exit when signal INT, TERM, KILL received
// 4. Always wait for 3 seconds before exit
func main() {
	p, err := rollover.Wait()
	if err != nil {
		log.Println(err)
	}
	pid := os.Getpid()
	log.Println("Current PID:", pid)
	if p == nil {
		log.Println(pid, ": this is a parent")
	} else {
		log.Println(pid, ": this is a child")
	}

	// simulate clean up before exit
	defer func() {
		time.Sleep(3 * time.Second)
		log.Println(pid, ": safely exit after clean up")
	}()

	log.Println(pid, ": have arguments", os.Args)

	log.Println("kill -HUP", pid, " -> to rollover(restart).")
	log.Println("kill", pid, " -> to quit")

	sigKill := make(chan os.Signal, 1)
	signal.Notify(sigKill, os.Interrupt, syscall.SIGTERM, os.Kill)

	sigRestart := make(chan os.Signal, 1)
	signal.Notify(sigRestart, syscall.SIGHUP)

	select {
	case <-sigKill:
	case <-sigRestart:
		log.Println(pid, ": initiating rollover")
		p, err := rollover.Restart()
		if err != nil {
			log.Println(pid, ": error rollover child", err)
		}
		if p == nil {
			log.Println(pid, ": error starting child")
		} else {
			log.Println(pid, ": child started runing, pid:", p.Pid)
		}
	}
}
