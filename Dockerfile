FROM golang:1.15.2
WORKDIR /go/src/checker/
COPY *.go .
RUN go get -d -v -t
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o checker .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /checker/
COPY --from=0 /go/src/checker/checker .
CMD ["./checker"]