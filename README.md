# Description
Web crawler that generates a site map

## Quick start

Build
```shell script
make build
```

Print usage
```shell script
./crawler -help
```

Generate the site map of a site (change the url and output file accordingly)
```shell script
./crawler -url http://site.com -out sitemap.txt
```

Run tests
```shell script
make test
```

## About
This is a simple web crawler that generates a site map.
It will fetch the provided URL and parse its links saving those links in an 
in-memory map, for each new link it will do the same mentioned steps (i.e.: fetch it and parse 
its links), each fetch and parse will be run in a concurrent Go routine.

This is a diagram of the program where the `Rx`s are routines:

```
R1 - url: "/", parsed links: ["/a", "/b", "/c"] -+
                                                 |
R2 - url: "/b", parsed links: ["/c", "/d"] ------| ----> R0 - save in in-memory map 
                                                 |            & for each new path
R3 - url: "/a", parsed links: ["/", "/d"] -------+            fetch its links
                                                                  |
R4 - url: "/d" fetching links...                                  |
 ^                                                                |
 |                                                                |
 +---------------------------- start routine ---------------------+
```
`R0` is a routine that listens to the routines that are 
fetching and parsing the URLs, `R0` will save in-memory the new links and it will launch a new 
routine for each unseen link in order fetch it and parse its links.

When we don't see any new links and all the go routines have finished parsing their pages we print
the results to the output file.

The file format is similar to this:
```
/
 /a
 /b
  /d
 /c
```
The root is at the top with no spaces and then its children in new lines with one space, 
we do not print a path more than once so it might be that "/a" has a link to "/" but "/" 
has already been printed in an upper level
