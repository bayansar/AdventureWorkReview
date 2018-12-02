FROM golang:1.9 as builder

# Set go bin which doesn't appear to be set already.
ENV GOBIN /go/bin

RUN go get -u github.com/golang/dep/...

COPY . ${GOPATH}/src/github.com/bayansar/AdventureWorkReview/
WORKDIR ${GOPATH}/src/github.com/bayansar/AdventureWorkReview/

# Go dep!
RUN dep ensure -vendor-only && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /review-app

RUN go test -v

# STEP 2 build a small image
FROM scratch

ENV RABBIT_URI amqp://guest:guest@rabbitmq:5672
ENV MYSQL_USER root
ENV MYSQL_PASSWORD 1234
ENV DB_NAME=adventureworks
ENV VALIDATE_QUEUE_NAME=validate
ENV NOTIFY_QUEUE_NAME=notify
ENV BAD_WORDS=fee,nee,cruul,leent
ENV DB_HOST=mysql

COPY --from=builder /review-app /review-app
EXPOSE 8888

CMD ["/review-app"]