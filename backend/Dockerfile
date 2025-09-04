FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o server .

FROM alpine:3.20
RUN apk add --no-cache bash curl ca-certificates docker-cli
WORKDIR /app
COPY --from=build /app/server /app/server
ENV PORT=8080
EXPOSE 8080
CMD ["/app/server"]
