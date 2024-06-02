FROM golang:1.21.5-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o main passVault/cmd

FROM alpine:latest
COPY --from=build /app/main /app/main
COPY --from=build /app/config.yaml /app/config.yaml
COPY --from=build /app/secrets.yaml /app/secrets.yaml
WORKDIR /app
CMD ["./main", "server"]
