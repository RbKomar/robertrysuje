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

ENV PORT=3000
EXPOSE 3000

CMD ["./main"]