FROM golang:1.21

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build .
EXPOSE 8080
EXPOSE 50051

CMD ["./storage"]