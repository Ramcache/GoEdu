FROM golang:1.23.2 AS builder

WORKDIR /GoEdu

COPY . .

RUN go mod tidy
RUN go build -o main ./cmd/server/main.go

FROM debian:bookworm-slim

WORKDIR /GoEdu

COPY --from=builder /GoEdu/main .

ENV POSTGRES_HOST=${POSTGRES_HOST}
ENV POSTGRES_PORT=${POSTGRES_PORT}
ENV POSTGRES_USER=${POSTGRES_USER}
ENV POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
ENV POSTGRES_DB=${POSTGRES_DB}

EXPOSE 8080

CMD ["./main"]
