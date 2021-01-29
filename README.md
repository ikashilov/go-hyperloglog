# go-hyperloglog

The Golang HyperLogLog implementation from the [original paper](http://algo.inria.fr/flajolet/Publications/FlFuGaMe07.pdf)

**HyperLogLog** is an algorithm for the count-distinct problem, approximating the number of distinct elements in a multise

The usage is quite simple:
```go
    var someData []string
    // fill or obtain someData

    hll, err := hyperloglog.New(.001)
    if err != nil {
        log.Fatalln(err)
    }

    for _, x := range someData {
        hll.Add(x)
    }

    fmt.Println("Estimated cardinality:", hll.Count())
```

### Quick start
```
$ go get github.com/ikashilov/go-hyperloglog
$ cd $GOPATH/src/github.com/ikashilov/go-hyperloglog
$ go test -test.v
$ go test -bench=.
```
