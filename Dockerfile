FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod init robrysuje
RUN go build -o main .

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/main .
COPY static/ ./static/
COPY templates/ ./templates/

EXPOSE 8080

CMD ["./main"]