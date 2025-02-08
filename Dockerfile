FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY Makefile /app/Makefile
COPY ./static/css/styles.css /app/static/css/styles.css

# RUN apk add --no-cache make curl
# RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 && \
#     chmod +x tailwindcss-linux-x64 && \
#     mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss

# RUN make tailwind-build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/static/css/output.css ./static/css/

EXPOSE 8080

CMD ["./server"]