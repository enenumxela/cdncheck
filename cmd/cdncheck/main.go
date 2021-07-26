package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/projectdiscovery/cdncheck"
)

var (
	inputList   string
	concurrency int
)

func init() {
	flag.StringVar(&inputList, "input-list", "", "")
	flag.StringVar(&inputList, "iL", "", "")
	flag.IntVar(&concurrency, "concurrency", 10, "")
	flag.IntVar(&concurrency, "c", 10, "")

	flag.Usage = func() {
		h := "USAGE:\n"
		h += "  cdncheck [OPTIONS]\n"

		h += "\nOPTIONS:\n"
		h += "  -iL, --input-list   input IP list (use `iL -` to read from stdin)\n"
		h += "   -c, --concurrency  number of concurrent threads (default: 10)\n"

		fmt.Fprint(os.Stderr, h)
	}

	flag.Parse()
}

func main() {
	IPs := make(chan string, concurrency)

	go func() {
		defer close(IPs)

		var scanner *bufio.Scanner

		if inputList == "-" {
			stat, err := os.Stdin.Stat()
			if err != nil {
				log.Fatalln(errors.New("no stdin"))
			}

			if stat.Mode()&os.ModeNamedPipe == 0 {
				log.Fatalln(errors.New("no stdin"))
			}

			scanner = bufio.NewScanner(os.Stdin)
		} else {
			openedFile, err := os.Open(inputList)
			if err != nil {
				log.Fatalln(err)
			}
			defer openedFile.Close()

			scanner = bufio.NewScanner(openedFile)
		}

		for scanner.Scan() {
			if scanner.Text() != "" {
				IPs <- scanner.Text()
			}
		}

		if scanner.Err() != nil {
			log.Fatalln(scanner.Err())
		}
	}()

	wg := &sync.WaitGroup{}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			client, err := cdncheck.NewWithCache()
			if err != nil {
				log.Fatal(err)
			}

			for IP := range IPs {
				if found, err := client.Check(net.ParseIP(IP)); found && err == nil {
					// log.Println("ip is part of cdn")
				} else {
					fmt.Println(IP)
				}
			}
		}()
	}

	wg.Wait()
}
