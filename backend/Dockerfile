FROM golang:1.20-alpine AS build
WORKDIR /app
COPY go.mod .
COPY . .
RUN go build -o hn-gateway .

FROM alpine:3.18
WORKDIR /app
COPY --from=build /app/hn-gateway /app/hn-gateway
EXPOSE 8080
CMD ["/app/hn-gateway"]
