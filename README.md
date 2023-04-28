# beginng Go lang source directory

there are some advantages using Golang

(cross platform)

cross platform build example(in case of the Raspberry PI)
$ GOOS=linux GOARCH=arm GOARM=6 go build hello_go.go

(compiler)

faster than script languages like Python and Java Script

(integrated web server)

"net/http" library works first enough

(secure)

GC function help memory leak issue

(applicable to IoT and Robotics)

Gobot can help these functionality in Raspberry PI

(easy to use multi core processor)

Native support for the multi core processor, don't request complicated description to make it multi threads

(easy multi thread synchronization)

there are two alternatives, "WaitGroup" or "channel"

To add mifare card check in application @serial using Node.js

(directries)

concpara:pararell processing

myfare: myfare card application folder

serial: read myfare card info and display it using Node.js(now it was updated to myfare folder golang app)

webserver: very primitive webserver function

go_chat: simple chat app using gin/melody/sqlite3
