FROM golang:1.23.1-alpine
WORKDIR /myapp

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY docker.env .env

RUN go build -ldflags '-w -s' -a -o ./bin/api ./cmd/api \
    && go build -ldflags '-w -s' -a -o ./bin/migrate ./cmd/migrate

CMD ["/myapp/bin/api"]
EXPOSE 8080