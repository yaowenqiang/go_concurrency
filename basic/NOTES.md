# concurrency in go

each goroutine use only 2kb  memory
each Thread in operation system use 2Mb

1000 goroutines = 1 Native Thread

## Go is Build for concurrency

+ goroutines are easy to create and manage
+ Channels make moving data around easily
+ Mutexes and locks are intergrated so you get good wranings
+ the whole language is build around these concepts
+ recall the select example, where we used a special syntax for channels

Go's designers chose M:N concurrency for a reasy

it's easy to understand ,has low overhead, and avoids callback hell
