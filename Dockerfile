FROM golang:1.22-alpine as builder
WORKDIR /app/crawler
COPY [ "./", "./" ]
RUN go build -ldflags='-w -s' .

FROM alpine
WORKDIR /app
COPY --from=builder [ "/app/crawler", "./" ]
CMD ["/app/crawler"]