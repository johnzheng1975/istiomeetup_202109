FROM golang:1.14.9 as builder

ADD . /go/src/servicetest
WORKDIR /go/src/servicetest

# build the source
RUN make build

FROM alpine

# set working directory
WORKDIR /go/bin

# copy the binary from builder
COPY --from=builder /go/src/servicetest/main /go/bin/servicetest

CMD ["/go/bin/servicetest"]

EXPOSE 8080
