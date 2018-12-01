FROM golang:1.9 as builder

# Set go bin which doesn't appear to be set already.
ENV GOBIN /go/bin

RUN go get -u github.com/golang/dep/...

COPY . ${GOPATH}/src/myapp
WORKDIR ${GOPATH}/src/myapp

# Go dep!
RUN dep ensure -vendor-only && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /myapp

RUN go test -v

# STEP 2 build a small image
FROM scratch

COPY --from=builder /myapp /myapp
EXPOSE 8000

CMD ["/myapp"]