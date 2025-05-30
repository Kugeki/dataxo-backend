FROM golang:1.24.2-alpine3.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o server ./cmd/server


FROM alpine:latest as runner

WORKDIR /app/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]