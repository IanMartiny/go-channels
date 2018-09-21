package main

import (
	"bytes"
	"channels/list"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createFile() {
	count := 5
	IPs := list.GenList(count)
	var data []byte
	for _, ip := range IPs {
		data = append(data, []byte(ip)...)
		data = append(data, []byte("\n")...)
	}
	ioutil.WriteFile("IPList", data, 0644)
}

func parseIPs(data []byte) []string {
	var lastInd int
	var stringIPs []string
	for i, b := range data {
		if b == '\n' {
			stringIPs = append(stringIPs, string(data[lastInd:i]))
			lastInd = i + 1
		}
	}

	return stringIPs
}

func sendIP(c chan string, v string) {
	c <- v
}

func receiveIP(c chan string) {
	v := <-c
	fmt.Printf("Received: %s\n", v)
	split := strings.Split(v, ".")
	if len(split) == 4 {
		var updated bytes.Buffer
		fmt.Printf("Once I've worked with it, I can send something back\n")
		updated.WriteString(v)
		updated.WriteString(".did work")

		// need to send completed "work", if it is just sent in this thread the
		// thread will halt, either the receive or the send (or both) need to be
		// in a new thread

		go func() { c <- updated.String() }()
	}
}

func fileSetup() []string {
	data, err := ioutil.ReadFile("IPList")
	var IPs []string
	if err != nil {
		// if it doesn't exist, I create one of fake IPs
		createFile()
	} else {
		// If it does exist, then parse that file
		IPs = parseIPs(data)
	}

	return IPs
}

func main() {
	// open file of IP addresses
	IPs := fileSetup()

	// make a wait group, without the waitGroup the main thread will finish
	// immediately, killing all children processes
	var wg sync.WaitGroup
	for _, v := range IPs {
		// make a channel to communicate between go routines
		var c = make(chan string)
		// add one to the waitGroup, at the end of the program we tell the main
		// thread to wait until the waitGroup is 0
		wg.Add(1)
		// go <function> creates a new child routine to run in parallel, here I
		// am using an anonymous function
		go func(ch chan string) {
			// call a function with the channel
			receiveIP(ch)

			// reading from a channel will halt execution of a routine (that's
			// why sending and receiving are done in new threads) this will wait
			// till someone sends something over the channel.
			retVal := <-ch
			fmt.Printf("returned from receiveIP: %s\n", retVal)
			// Done() subtracts one from the waitGroup, potentially allowing
			// main to complete defer waits for the encapsulating function to
			// complete before executing. Thus I want receive to finish
			// completely before I tell the waitGroup to count this process as
			// done
			defer wg.Done()
		}(c)
		// call another function in parallel
		go func(ch chan string, val string) {

			// send a value over the channel to be received. This one takes an
			// input, because this anonymous function will probably not actually
			// run until the main loop has completed meaning if I used 'v' directly
			// it would always use the last value in sendIP, as an input to a function
			// it is saved each iteration.
			sendIP(ch, val)
		}(c, v)
	}
	// tell the main thread to wait until all go routines are done
	wg.Wait()
}
