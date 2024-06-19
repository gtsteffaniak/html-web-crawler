FROM golang:1.22-alpine as builder
WORKDIR /app/crawler
COPY [ "./", "./" ]
RUN go build -ldflags='-w -s' .

FROM node:alpine
RUN apk add --no-cache musl-dev
RUN npm init playwright@latest
RUN npx npx playwright install
WORKDIR /app
COPY --from=builder [ "/app/crawler", "./" ]
CMD ["/app/crawler"]