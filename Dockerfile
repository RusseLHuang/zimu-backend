FROM golang:1.14.0-alpine3.11 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go build -o main .
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Second stage run the binary file
FROM alpine
EXPOSE 8777
RUN mkdir /app

# Copy from first stage
COPY --from=builder /app/main /app/main
COPY config.* /

ENTRYPOINT ["/app/main"]
