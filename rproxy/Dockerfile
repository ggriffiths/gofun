FROM golang:latest
ADD . /go/src/github.com/ggriffiths/gofun/rproxy
RUN go build /go/src/github.com/ggriffiths/gofun/rproxy/cmd/rproxy
CMD ["./rproxy"]