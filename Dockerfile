FROM golang:alpine as builder
RUN apk --no-cache add git dep ca-certificates
RUN mkdir -p /go/src/katalog
ADD . /go/src/katalog
WORKDIR /go/src/katalog
RUN dep ensure
WORKDIR /go/src/katalog/src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/katalog/src /app/
WORKDIR /app
CMD ./main