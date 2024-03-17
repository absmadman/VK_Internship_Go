FROM golang:alpine

ADD go.mod .

COPY . .

EXPOSE 8080

CMD ["go", "run", "main.go"]
