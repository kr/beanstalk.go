# beanstalk.go

beanstalk.go is a client library for the protocol used by [beanstalkd][].

## Installation

Use the `goinstall` command, and there is no need for installation. Just use
`import "github.com/kr/beanstalk.go.git"` at the top of your go package.

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

[beanstalkd]: http://kr.github.com/beanstalkd/
