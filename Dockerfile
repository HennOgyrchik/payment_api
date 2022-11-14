FROM golang:1.18.4 as builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o carnival ./cmd/main.go

FROM alpine:latest
COPY --from=builder /src/carnival /
CMD /carnival