FROM golang:alpine as builder
RUN apk --no-cache add git ca-certificates
RUN mkdir -p /go/src/github.com/walmartdigital/katalog
ADD . /go/src/github.com/walmartdigital/katalog
WORKDIR /go/src/github.com/walmartdigital/katalog
WORKDIR /go/src/github.com/walmartdigital/katalog
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM alpine
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/walmartdigital/katalog/main /app/
COPY --from=builder /go/src/github.com/walmartdigital/katalog/health.sh /app/
ENTRYPOINT ["/app/main"]