# example-webserver

The plan was for this example webserver to use the [fasthttp](https://github.com/valyala/fasthttp) golang library due to its lower memory footprint over the native net/http golang library, as well as other [significant](https://github.com/smallnest/go-web-framework-benchmark#basic-test) benchmarking results. However, the external APIs we're sending requests to are not able to handle an incredible amount of load such that the faster HTTP golang library shouldn't matter _that much_.

We will include our "benchmark" results by running tests using [wrk](https://github.com/wg/wrk)

We will be sending requests to 2 external services:
- http://uinames.com/api/
- http://api.icndb.com/jokes/random?firstName=__PAUL__&lastName=__BLART__&limitTo=[nerdy]

After profiling both of the above URLs (results in benchmarks directory), it is apparent the first one (uinames) does not allow as much traffic as the second by a significant amount. In order to work around this, we will be caching 500 random names.

## Benchmarks

### Setup
I noticed during benchmarking that I was seeing socket errors relating to too many files being opened.
```
{
  "file": "v1.go:48",
  "func": "randomNameJoke",
  "level": "info",
  "msg": "request to joke either returned an error or status code != 200dial tcp4 104.24.114.160:80: socket: too many open files"
}
```

I have seen this before, so I checked the limit on my machine by running `uname -n`. 256 on MacOS (argh).
Even on an Ubuntu 18 VM, the default is 1024. So I increased it:
```
$ echo "fs.file-max=400000" >> /etc/sysctl.conf
$ sysctl -p

# And for good measure
$ reboot
```

### Results
I ran the server on a machine with the following specs:
```
 OS: Ubuntu 18.04 bionic
 Kernel: x86_64 Linux 4.15.0-51-generic
 Uptime: 9m
 Packages: 557
 Shell: su
 CPU: Intel Xeon E312xx (Sandy Bridge, IBRS update) @ 4x 2GHz
 GPU:
 RAM: 179MiB / 15032MiB
```

```
Running 2m test @ http://localhost:5000/v1/
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   176.51ms   32.00ms   1.67s    95.90%
    Req/Sec   185.53     42.64   333.00     74.00%
  Latency Distribution
     50%  170.86ms
     75%  173.76ms
     90%  178.00ms
     99%  328.44ms
  264893 requests in 2.00m, 48.80MB read
  Socket errors: connect 0, read 85, write 0, timeout 0
Requests/sec:   2206.00
Transfer/sec:    416.12KB
```

## Building
There are 2 options for building and running this example webserver  
### 1 Local dev env is setup
```
# Clone repo
$ git clone git@github.com:CamelCaseNotation/example-webserver.git

# Download latest Go version
# https://github.com/travis-ci/gimme#installation--usage
$ brew install gimme
$ gimme stable

# Download all dependencies
$ go mod tidy

# Build the binary
# This creates the binary at $(pwd)/bin/server
$ go build -o bin/server

# Run the server
$ ./bin/server
```
### 2 Docker is installed
```
# You can also just build the docker image and run the container

$ docker build .

# Previous command will end with "Successfully built __ID__"
# To start a container using that image as a daemon and exposing the default port of 5000 such that you can send requests to it via localhost:5000, run command below

$ docker run -d --name example-webserver -p 5000:5000 aff5bc597b5c

# See the container running
$ docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED              STATUS              PORTS                    NAMES
30709672cdd1        aff5bc597b5c        "/server"           About a minute ago   Up About a minute   0.0.0.0:5000->5000/tcp   example-webserver

# Stop and remove container when done
$ docker stop 30709672cdd1 && docker rm 30709672cdd1
```