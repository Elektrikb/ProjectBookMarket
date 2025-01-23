FROM golang:1.23.1-alpine

WORKDIR /Projectmugen

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o Projectmugen ./

EXPOSE 8080

CMD ["./Projectmugen"]