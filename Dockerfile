# build stage
FROM golang:1.21.3-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o bank-sys main.go

# run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/bank-sys .
COPY app.env .
EXPOSE 8080
CMD [ "/app/bank-sys" ]