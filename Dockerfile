FROM golang:1.20 as base

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY . .
RUN git config --global --add safe.directory /app && \
    go mod download && \
    go mod verify && \
    go build -v -o sponty .

FROM golang:1.20 as run

COPY --from=base /app/sponty /usr/local/bin/sponty

CMD ["sponty"]