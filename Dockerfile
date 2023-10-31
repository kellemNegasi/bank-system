# build stage
FROM golang:1.21.3-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o bank-sys main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
# run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/bank-sys .
COPY --from=builder /app/migrate ./migrate
COPY start.sh .
COPY app.env .
COPY db/migration ./migration
EXPOSE 8080
CMD [ "/app/bank-sys" ]
ENTRYPOINT [ "/app/start.sh" ]