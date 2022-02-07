FROM golang:1.17.6

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY . .
RUN go mod download && \
    go mod verify && \
    go test ./... && \
    go build -v -o /usr/local/bin/sponty ./...

CMD ["sponty"]