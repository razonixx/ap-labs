Lab 2.5 - Depth-limiting Web Crawler
========================

Description
--------------------------
This is a web crawler with user-defined depth written in go.

Compilation
--------------------
To build the binary, run `go build crawl3.go`.

The program recieves 2 arguments, the depth desired and the URL to crawl.

Usage: `./crawl -depth=n link` where n is an integer and link is the URL.

Examples: 

- `./crawl3 -depth=2 https://google.com`
- `./crawl3 -depth=3 http://www.gopl.io/`
- `./crawl3 -depth=1  http://www.gopl.io/`
