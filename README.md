WARNING
-------

**The library is KNOWN TO HAVE BUGS!!!** Unfortunately, currently I don't have resources to investigate them thoroughly enough and in timely fashion. In case somebody is interested in taking ownership of the library, I'm open to ceding it. That said, the issues totally haunt me and occasionally I stubbornly try to come back to them and pick the fight up again. In particular:

- #3 was confirmed to be **an omission in the original paper/algorithm**. As far as I understand, it surfaces when one of the polygons used has self-overlapping edges (e.g. when an edge (0,0)-(1,1) is used twice in the same polygon). I believe it should be possible to fix, but it requires thorough analysis of the algorithm and good testing. One attempt I made at a fix which seemed OK initially was later found to break the library even more and thus I reverted it.
- #8 was reported recently and I haven't yet had time to even start investigating it.

About
-----

[![Build Status on Travis-CI.](https://travis-ci.org/akavel/polyclip-go.svg?branch=master)](https://travis-ci.org/akavel/polyclip-go)
![License: MIT.](https://img.shields.io/badge/license-MIT-orange.svg)
[![Documentation on godoc.org.](https://godoc.org/github.com/akavel/polyclip-go?status.svg)](https://godoc.org/github.com/akavel/polyclip-go)

Library polyclip-go is a pure Go, MIT-licensed implementation of an [algorithm for Boolean operations on 2D polygons] [fmartin] (invented by F. Martínez, A.J. Rueda, F.R. Feito) -- that is, for calculation of polygon intersection, union, difference and xor.

The original paper describes the algorithm as performing in time _O((n+k) log n)_, where _n_ is number of all edges of all polygons in operation, and _k_ is number of intersections of all polygon edges.

  [fmartin]: http://wwwdi.ujaen.es/~fmartin/bool_op.html

  ![](http://img684.imageshack.us/img684/5296/drawqk.png "Polygons intersection example, calculated using polyclip-go")

Example
-------

Simplest Go program using polyclip-go for calculating intersection of a square and triangle:

    // example.go
    package main
    
    import (
        "fmt"
        "github.com/akavel/polyclip-go" // or: bitbucket.org/...
    )
    
    func main() {
        subject := polyclip.Polygon{{{1, 1}, {1, 2}, {2, 2}, {2, 1}}} // small square
        clipping := polyclip.Polygon{{{0, 0}, {0, 3}, {3, 0}}}        // overlapping triangle
        result := subject.Construct(polyclip.INTERSECTION, clipping)

        // will print triangle: [[{1 1} {1 2} {2 1}]]
        fmt.Println(result)
    }

To compile and run the program above, execute the usual sequence of commands:

    go get github.com/akavel/polyclip-go  # or: bitbucket.org/...
    go build example.go
    ./example      # Windows: example.exe

For full package documentation, run locally `godoc github.com/akavel/polyclip-go`, or visit [online documentation for polyclip-go][godoc].
    
  [godoc]: http://godoc.org/github.com/akavel/polyclip-go
    
See also
--------
  * [Online docs for polyclip-go][godoc].
  * Microsite about [the original algorithm][fmartin], from its authors (with PDF, and public-domain code in C++).
  * The [as3polyclip] library -- a MIT-licensed ActionScript3 library implementing this same algorithm (it actually served as a base for polyclip-go). The page also contains some thoughts with regards to speed of the algorithm.
  
  [as3polyclip]: http://code.google.com/p/as3polyclip/
