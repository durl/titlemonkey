// Copyright (c) 2016, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var verbose bool

func main() {
	verbose = true // TODO: use cli flag
	if len(os.Args) < 2 {
		exitUsage("no command specified")
	}
	cmd := os.Args[1]
	switch cmd {
	case "fetch":
		cmdFetch()
	case "gen":
		cmdGen()
	default:
		exitUsage("unknown command")
	}
}

const usage = `usage: monkeytitles fetch|gen
	fetch <url>: fetch titles from url and print to stdout.
	gen <num>:   generate <num> random titles based on the titles provided to stdin.
`

func exitUsage(msg string) {
	fmt.Fprintln(os.Stderr, usage)
	fmt.Fprintln(os.Stderr, "error: "+msg)
	os.Exit(2)
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func cmdFetch() {
	if len(os.Args) != 3 {
		exitUsage("illegal number of arguments")
	}
	url := os.Args[2]

	titles, err := fetchRssTitles(url)
	if err != nil {
		exit(err)
	}
	for _, title := range titles {
		fmt.Println(title)
	}
}

func cmdGen() {
	// params
	if len(os.Args) != 3 {
		exitUsage("illegal number of arguments")
	}
	num, err := strconv.ParseUint(os.Args[2], 10, 64)
	if err != nil {
		exit(err)
	}

	// build markov chain
	beginTime := time.Now()
	markov, originals := BuildMarkovChain(os.Stdin)
	buildCompleteTime := time.Now()

	// generate titles
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	var loops int
	var origMatches int
	for i := uint64(0); i < num; {
		title, err := markov.GenerateTitle(rnd)
		if err != nil {
			loops++
			//fmt.Fprintln(os.Stderr, err.Error())
		}
		if isKnown(originals, title) {
			origMatches++
			continue
		}
		fmt.Println(title)
		i++
	}
	endTime := time.Now()

	// stats
	if verbose {
		printStats(originals, beginTime, buildCompleteTime, endTime, loops, origMatches)
	}
}

const statsMsg = `[Parameter Overview]
	input titles:  %d
	longest input: %d words
	lookback:      %d words
	loop size:     %d words
[Performance]
	input analyzing:  %s
	title generation: %s
[Discarded Titles]
	%d contained loops
	%d matched original titles
`

func printStats(originals []string, beginTime, buildCompleteTime, endTime time.Time, loops, originalMatches int) {

	maxInWords := 0
	for _, title := range originals {
		maxInWords = max(maxInWords, len(strings.Fields(title)))
	}

	buildTime := buildCompleteTime.Sub(beginTime).String()
	genTime := endTime.Sub(buildCompleteTime).String()

	fmt.Fprintf(os.Stderr, statsMsg,
		len(originals), maxInWords, lookback, loopSize, // overview
		buildTime, genTime, // performance
		loops, originalMatches) // discarded
}

func isKnown(originals []string, title string) bool {
	for _, orig := range originals {
		if strings.Contains(orig, title) {
			return true
		}
	}
	return false
}
