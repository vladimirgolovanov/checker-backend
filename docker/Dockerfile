FROM golang:latest AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o checker .

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=build /app/checker .
CMD ["/checker"]
