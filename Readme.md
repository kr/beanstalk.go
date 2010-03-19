# beanstalk.go

beanstalk.go is a client library for the protocol used by [beanstalkd][].

## Installation

No need for that. Just put `import "github.com/kr/beanstalk.go.git"` at the
top of your go package, and install your own package with [goinstall][].

## Example

A producer:

    package main

    import "github.com/kr/beanstalk.go.git"

    func main() {
        t := beanstalk.Open("localhost:11300").Tube("default")
        c.put("hello")
    }

And a worker:

    package main

    import "github.com/kr/beanstalk.go.git"

    func main() {
        ts := beanstalk.Open("localhost:11300").Tubes([]string{"default"})
        for {
            j, err := ts.Reserve()
            fmt.Println(j.Body) // prints "hello"
            j.Delete()
        }
    }

## Credit Where It's Due

 * [spymemcached][] for the idea of making optimizing transformations on the
   command stream.

 * Go's standard libraries, especially net and http, for inspiration and
   guidance.

[beanstalkd]: http://kr.github.com/beanstalkd/
[spymemcached]: http://code.google.com/p/spymemcached/
[goinstall]: http://golang.org/cmd/goinstall/
