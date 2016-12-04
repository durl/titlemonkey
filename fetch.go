// Copyright (c) 2016, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/mmcdole/gofeed"

func fetchRssTitles(url string) ([]string, error) {

	// fetch rss
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url)
	if err != nil {
		return nil, err
	}

	// extract titles
	titles := make([]string, len(feed.Items))
	for i, item := range feed.Items {
		titles[i] = item.Title
	}
	return titles, nil

}
