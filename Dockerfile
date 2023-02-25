FROM golang:1.20-alpine as builder
RUN mkdir /build
COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build -o hello-app

FROM alpine:3.17
COPY --from=builder /build/hello-app .
ENTRYPOINT ["./hello-app"]
