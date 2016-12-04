// Copyright (c) 2016, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"errors"
	"io"
	"math/rand"
	"strings"
)

// lookback specifies how many words are used for a prefix.
const lookback = 2

const loopSize = 6

// Suffix represents the occurences of a following word in a markov chain.
type Suffix struct {
	Suffix      string
	Occurences  int
	Probability float32
}

// MarkovMap represents a markov chain with a multi word prefix as key,
// and a list of suffixes as value of the map.
type MarkovMap map[string][]Suffix

// BuildMarkovChain reads titles separated by newlines from a given input,
// and builds a markov chain, aswell as a list of normalized (no unecessary
// whitespace) original titles. This list can be used to discard generated
// titles which closely resemble original ones.
func BuildMarkovChain(input io.Reader) (MarkovMap, []string) {
	markov := make(MarkovMap)
	originals := make([]string, 0)
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		title := scanner.Text()
		orig := markov.addTitle(title)
		originals = append(originals, orig)
	}
	markov.calculateProbabilities()
	return markov, originals
}

func (m MarkovMap) addTitle(title string) string {
	words := strings.Fields(title)
	if len(words) <= lookback {
		// ignore titles which are too short
		return join(words)
	}

	for i, word := range words {
		prefix := join(words[max(0, i-lookback):i])
		var suffixes []Suffix
		if existingSuffixes, ok := m[prefix]; ok {
			suffixes = existingSuffixes
		} else {
			suffixes = make([]Suffix, 0)
		}

		var suffix Suffix
		suffix.Suffix = word
		new := true
		for j, existingSuffix := range suffixes {
			if existingSuffix.Suffix == suffix.Suffix {
				new = false
				suffixes[j].Occurences++
				break
			}
		}

		if new {
			suffix.Occurences++
			suffixes = append(suffixes, suffix)
			m[prefix] = suffixes
		}
	}

	// return normalized original title
	return join(words)
}

func join(words []string) string {
	return strings.Join(words, " ")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m MarkovMap) calculateProbabilities() {
	for _, suffixes := range m {
		var total int
		for _, suffix := range suffixes {
			total += suffix.Occurences
		}
		for i, suffix := range suffixes {
			suffix.Probability = float32(suffix.Occurences) / float32(total)
			suffixes[i] = suffix
		}
	}
}

func sample(suffixes []Suffix, rnd *rand.Rand) Suffix {
	for {
		i := rnd.Intn(len(suffixes))
		if rnd.Float32() < suffixes[i].Probability {
			return suffixes[i]
		}
	}
}

// GenerateTitle generates a title based on a markov chain.
// An error is returned if the title contains a loop of a certain length.
func (m MarkovMap) GenerateTitle(rnd *rand.Rand) (string, error) {
	title := make([]string, 0)
	next := sample(m[""], rnd)
	for next.Suffix != "" {
		title = append(title, next.Suffix)
		candidates, ok := m[join(title[max(len(title)-lookback, 0):])]
		if !ok {
			break
		}
		next = sample(candidates, rnd)
		if hasLoop(title) {
			return join(title), errors.New("loop detected")
		}

	}
	return join(title), nil
}

func hasLoop(words []string) bool {
	if len(words) <= loopSize {
		return false
	}
	loopStr := join(words[len(words)-loopSize:])
	if strings.Contains(join(words[:len(words)-loopSize]), loopStr) {
		return true
	}
	return false
}
