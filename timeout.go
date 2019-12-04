//+build ignore

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	flagDuration = flag.String("duration", "", "exit after this duration")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("No command specified")
	}
	cmd := exec.Command(args[0], args[1:]...)

	var wg sync.WaitGroup
	opr, opw := io.Pipe()
	cmd.Stdout = opw

	epr, epw := io.Pipe()
	cmd.Stderr = epw

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(os.Stdout, opr)
	}()
	go func() {
		defer wg.Done()
		io.Copy(os.Stderr, epr)
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)

	// Use a channel to signal completion so we can use a select statement
	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	// Start a timer
	var timeout <-chan time.Time

	if *flagDuration != "" {
		d, err := time.ParseDuration(*flagDuration)
		if err != nil {
			log.Fatalf("unable to parse duration: %v", err)
		}
		timeout = time.After(d)
	}

	// The select statement allows us to execute based on which channel
	// we get a message from first.
	defer wg.Wait()
wait:
	for {
		select {
		case <-timeout:
			// Timeout happened first, kill the process and print a message.
			fmt.Println("timeout reached...")
			err := cmd.Process.Signal(os.Kill)
			if err != nil {
				fmt.Println("Problem killing:", err)
				os.Exit(1)
			}
		case err := <-done:
			if err != nil {
				fmt.Println("Non-zero exit code:", err)
				os.Exit(1)
			}
			break wait
		case <-sigint:
			fmt.Println("Got SIGINT. Sending Kill to child")
			err := cmd.Process.Signal(os.Kill)
			if err != nil {
				fmt.Println("Problem killing:", err)
				os.Exit(1)
			}
		}
	}
	opw.Close()
	epw.Close()
}
