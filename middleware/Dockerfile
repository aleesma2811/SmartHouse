FROM golang:1.27rc2-alpine3.24 AS build

WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN go build -o middleware .

FROM alpine:3.20

WORKDIR /app
COPY --from=build /app/middleware .

EXPOSE 8000

CMD ["./middleware"]
