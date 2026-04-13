package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/edsrzf/mmap-go"
	"github.com/zeebo/xxh3"
)

var cache = make(map[uint64]struct{})

func main() {
	if len(os.Args) == 1 || os.Args[1] == "--help" {
		_, _ = fmt.Fprintln(os.Stderr, "usage: anew file")
		os.Exit(2)
	}

	// 1. cache the files hashed lines
	f, err := os.Open(os.Args[1])

	if err != nil {
		log.Fatalln(err)
	}

	m, err := mmap.Map(f, mmap.RDONLY, 0)

	if err != nil {
		log.Fatalln(err)
	}

	sc := bufio.NewScanner(bytes.NewReader(m))

	for sc.Scan() {
		cache[xxh3.HashString(sc.Text())] = struct{}{}
	}

	_ = m.Unmap()
	_ = f.Close()

	// 2. reopen the file for writing
	f, err = os.OpenFile(os.Args[1], os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		log.Fatalln(err)
	}

	// 3. write unique new lines to the file
	sc = bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		val := sc.Text()
		key := xxh3.HashString(val)

		if _, ok := cache[key]; !ok {
			cache[key] = struct{}{}
			_, _ = fmt.Fprintf(f, "%s\n", val)
		}
	}

	_ = f.Close()
}
