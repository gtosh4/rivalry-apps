FROM golang:1.14 AS build
WORKDIR /app
COPY . .
RUN go build -mod=vendor -ldflags="-s -w" -o /app/rivalry-apps ./cmd/rivalry-apps

FROM debian
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean
WORKDIR /app
COPY --from=build /app/rivalry-apps ./
ENTRYPOINT ["/app/rivalry-apps"]
