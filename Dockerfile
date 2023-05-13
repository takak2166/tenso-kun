FROM golang:latest as builder
WORKDIR /opt/
COPY ./ /opt/
ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build -o tenso-kun main.go line.go send.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /opt/tenso-kun ./
CMD ["./tenso-kun"]