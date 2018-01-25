// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio" // for version 1.
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
)
    
var (
	input_chan chan int
)

func cmd_input() {
	// version 1
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadBytes('\n')
	cmd_string := string(input[0 : len(input)-1])

	switch cmd_string {
	case "login":
		
		break
	}

	input_chan <- 1
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	input_chan = make(chan int, 0)
	go cmd_input()

	for {

		select {
		case <-input_chan:
			go cmd_input()
		case <-interrupt:
			log.Println("interrupt")
			select {
			case <-time.After(1 * time.Second):
			}
			return
		}
	}
}
