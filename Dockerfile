FROM golang:1.9.1 as builder
RUN mkdir -p /go/src/github.com/puppetlabs/lumogon/
COPY . /go/src/github.com/puppetlabs/lumogon/
WORKDIR /go/src/github.com/puppetlabs/lumogon/
RUN make build

FROM scratch
COPY --from=builder /go/src/github.com/puppetlabs/lumogon/bin/lumogon /
COPY --from=builder /go/src/github.com/puppetlabs/lumogon/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/lumogon"]
