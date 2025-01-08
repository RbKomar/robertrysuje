FROM golang:1.24rc1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates

RUN mkdir -p data

RUN adduser -D appuser
RUN chown -R appuser:appuser /app
USER appuser

EXPOSE 8080

CMD ["./main"]