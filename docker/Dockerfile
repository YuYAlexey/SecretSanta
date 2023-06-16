FROM golang:alpine

WORKDIR /SecretSanta
COPY . .

RUN go build -ldflags '-w -s' -a -o SecretSanta main.go

CMD ["./SecretSanta"]
EXPOSE 8080