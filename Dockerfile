# Build stage
FROM golang:1.21.1-alpine3.17 AS builder
ENV GOPROXY=https://goproxy.io,direct
WORKDIR /app
COPY . .
RUN go build -o main main.go
# Run stage
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder  /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration


EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]