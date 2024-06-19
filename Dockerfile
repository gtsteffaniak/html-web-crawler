FROM golang:1.22-alpine as builder
WORKDIR /app/crawler
COPY [ "./", "./" ]
RUN go build -ldflags='-w -s' .

FROM alpine
# chromium needed for javascript enabled querying
RUN apk add --no-cache chromium
ENV CHROME_EXECUTABLE="/usr/bin/chromium"
WORKDIR /app
COPY --from=builder [ "/app/crawler", "./" ]
CMD ["/app/crawler"]