Golang HTML5 SSE BY NSQ AND REDIS
========================

This is an minimalistic example of how to do
[HTML5 Server Sent Events](http://en.wikipedia.org/wiki/Server-sent_events)
with [Go (golang)](http://golang.org/).  

The main advantage of HTML5 SSE over long polling is that there is a nice
API for it in modern browsers, so that you need not use iframes and such.
SSE is easier than Websockets in the sense that it communicates exclusively
over HTTP and therefore does not require a separate server.  Websockets,
however, supports two-way real-time communication between the client and
the server.

## Installing

Check out the repository from GitHub

	git clone https://github.com/zhangjunMaster/go-nsq-sse

## Running

First do
   install nsq

Then

    start nsq https://nsq.io/overview/quick_start.html

To run the server, do 

	go run ./server.go

Then point your web browser to `http://localhost:8001`.

## Thanks

This code is based on 

* Golang HTML5 SSE Example (https://github.com/kljensen/golang-html5-sse-example); 
* logger https://github.com/fang2329/logger

## License (the Unlicense)

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

For more information, please refer to <http://unlicense.org/>

