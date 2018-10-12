FROM golang:latest
ADD . /go/src/github.com/ggriffiths/gofun/rproxy
RUN chmod a+x  /go/src/github.com/ggriffiths/gofun/rproxy/test/run_tests.sh
CMD ["/go/src/github.com/ggriffiths/gofun/rproxy/test/run_tests.sh"]