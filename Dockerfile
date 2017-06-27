FROM golang:1.8.3 as builder
RUN mkdir -p /go/src/github.com/puppetlabs/lumogon/
COPY . /go/src/github.com/puppetlabs/lumogon/
WORKDIR /go/src/github.com/puppetlabs/lumogon/
RUN make build
# RUN go build -buildmode=plugin -o lumogon_diff_plugin.so dummyplugins/dummyplugin.go


FROM debian:jessie
ENV LUMOGON_ENDPOINT=https://consumer.app.lumogon.com/api/v1/
COPY --from=builder /go/src/github.com/puppetlabs/lumogon/bin/lumogon /
# COPY --from=builder /go/src/github.com/puppetlabs/lumogon/lumogon_diff_plugin.so /plugins/lumogon_diff_plugin.so
ENTRYPOINT ["/lumogon"]
