FROM golang:1.20-alpine
WORKDIR /usr/gateway
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o ./bin
CMD ["./bin/gateway"]
