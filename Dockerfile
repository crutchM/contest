FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download && \
    go build -o migrations ./cmd/migrator/main.go && \
    go build -o contest-app ./cmd/main.go

CMD ["./migrations && ./contest-app"]