About
-----

Library polyclip.go is a pure Go, MIT-licensed implementation of an [algorithm for Boolean operations on 2D polygons] [fmartin] (invented by F. Martínez, A.J. Rueda, F.R. Feito) &mdash; that is, for calculation of polygon intersection, union, difference and xor.

The original paper describes the algorithm as performing in time _O((n+k) log n)_, where _n_ is number of all edges of all polygons in operation, and _k_ is number of intersections of all polygon edges.

  [fmartin]: http://wwwdi.ujaen.es/~fmartin/bool_op.html

  ![](http://img684.imageshack.us/img684/5296/drawqk.png "Polygons intersection example, calculated using polyclip.go")

Example
-------

Simplest Go program using polyclip.go for calculating intersection of a square and triangle:

    // main.go
    package main
    
    import (
        "fmt"
        "github.com/akavel/polyclip.go" // or: bitbucket.org/...
    )
    
    func main() {
        subject := polyclip.Polygon{{{1, 1}, {1, 2}, {2, 2}, {2, 1}}} // small square
        clipping := polyclip.Polygon{{{0, 0}, {0, 3}, {3, 0}}}        // overlapping triangle
        result := subject.Construct(polyclip.INTERSECTION, clipping)

        // will print triangle: [[{1 1} {1 2} {2 1}]]
        fmt.Println(out)
    }

To compile and run the program above, execute the usual sequence of commands:

    goinstall -make=false github.com/akavel/polyclip.go  # or: bitbucket.org/...
    8g main.go   # or 6g, 5g, depending on your system
    8l main.8    # or: 6l main.6, 5l main.5
    ./a.out      # Windows: a.out.exe
    
For full package documentation, run locally `godoc github.com/akavel/polyclip.go`, or visit [online documentation for polyclip.go][gopkgdoc].
    
  [gopkgdoc]: http://gopkgdoc.appspot.com/pkg/github.com/akavel/polyclip.go
    
See also
--------
  * [Online godoc for polyclip.go][gopkgdoc] (courtesy of [http://gopkgdoc.appspot.com](http://gopkgdoc.appspot.com)).
  * Microsite about [the original algorithm][fmartin], from its authors (with PDF, and public-domain code in C++).
  * The [as3polyclip] library &mdash; a MIT-licensed ActionScript3 library implementing this same algorithm (it actually served as a base for polyclip.go). The page also contains some thoughts with regards to speed of the algorithm.
  
  [as3polyclip]: http://code.google.com/p/as3polyclip/