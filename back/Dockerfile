FROM golang:1.27rc2-alpine3.24

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum* ./

RUN go mod download

COPY . .

EXPOSE 4000

CMD ["air"]