# beanstalk.go

beanstalk.go is a client library for the protocol used by [beanstalkd][].

## Installation

No need for that. Just put `import "github.com/kr/beanstalk.go.git"` at the
top of your go package, and install your own package with [goinstall][].

## Overview

To open a connection, do

    conn, err := beanstalk.Dial("localhost:11300")

This package provides a simple, blocking interface. Go makes it easy to add
asynchrony if you want.

**Common Case:** To submit a job and get its id, do

    id, err := conn.Put(...)

**Fire and Forget:** If you don't care about the id, no need to wait around:

    go conn.Put(...)

**Full Asynchrony:** If you don't want to wait but still need the id, it's
still easy:

    go func() {
      id, err := conn.Put(...)
    }()

## Concurrency and Optimizations

You can perform operations on the queue at will from all goroutines. This
package will issue commands to beanstalkd and return results to where they
belong. It'll also optimize the commands a bit to reduce network traffic.

For example, the use, watch, and ignore commands are managed entirely by this
library, and issued only as necessary.

It also collects as many outbound commands as possible into a single network
packet, for efficiency.

In the future, it will manage the order of outgoing commands to further reduce
the need for use, watch, and ignore commands, and to prevent reserve commands
from stalling other commands unnecessarily.

## Complete Example

A producer:

    package main

    import "github.com/kr/beanstalk.go.git"

    func main() {
        conn, err := beanstalk.Dial("localhost:11300")
        conn.Put("hello", 0, 0, 10)
    }

And a worker:

    package main

    import "github.com/kr/beanstalk.go.git"

    func main() {
        conn, err := beanstalk.Dial("localhost:11300")
        for {
            j, err := conn.Reserve()
            fmt.Println(j.Body) // prints "hello"
            j.Delete()
        }
    }

## Contributing

 1. Fix a bug or implement a new feature.
 1. Add tests verifying your change.
 1. Publish your changes in a public git repository. At most one feature or
    bug fix per commit. Topic branches are preferred, especially for larger
    changes.
 1. Send [me email][] with the URL of your repository and the branch(es) you
    want pulled.

## Credit Where It's Due

 * [spymemcached][] for the idea of making optimizing transformations on the
   command stream.

 * Go's standard libraries, especially net and http, for inspiration and
   guidance.

[beanstalkd]: http://kr.github.com/beanstalkd/
[spymemcached]: http://code.google.com/p/spymemcached/
[goinstall]: http://golang.org/cmd/goinstall/
[me email]: mailto:kr@xph.us
