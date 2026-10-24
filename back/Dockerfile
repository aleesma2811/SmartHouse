# Etapa 1: compilar el binario
FROM golang:1.27rc2-alpine3.24 AS build

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /bin/app .

# Etapa 2: imagen final solo con el binario
FROM alpine:3.24

COPY --from=build /bin/app /app

EXPOSE 4000

# El binario es PID 1 y recibe SIGTERM directamente
CMD ["/app"]
