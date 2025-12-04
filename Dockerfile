FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o kanban_board_server .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/kanban_board_server .

EXPOSE 3000

CMD ["./kanban_board_server"]
