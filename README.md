About
-----

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

> **Note:** the polyclip-go library is developed for a fairly modern (at the time of writing) version of Go compiler, known as Go 1 RC1. If you have an older "release" version of Go, you may have problems compiling, and you are advised to switch to a newer "weekly" version. On the other hand, if you do have a newer version and encounter problems, please try using the `go fix` tool to update the polyclip-go library. I'll also be grateful if you could contact me about that.
    
For full package documentation, run locally `godoc github.com/akavel/polyclip-go`, or visit [online documentation for polyclip-go][gopkgdoc].
    
  [gopkgdoc]: http://godoc.org/github.com/akavel/polyclip-go
    
See also
--------
  * [Online godoc for polyclip-go][gopkgdoc] (courtesy of [http://gopkgdoc.appspot.com](http://gopkgdoc.appspot.com)).
  * Microsite about [the original algorithm][fmartin], from its authors (with PDF, and public-domain code in C++).
  * The [as3polyclip] library -- a MIT-licensed ActionScript3 library implementing this same algorithm (it actually served as a base for polyclip-go). The page also contains some thoughts with regards to speed of the algorithm.
  
  [as3polyclip]: http://code.google.com/p/as3polyclip/
