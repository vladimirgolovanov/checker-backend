FROM golang:latest AS build
WORKDIR /app
COPY . .
#RUN go get github.com/golang/dep/cmd/dep && \
#      dep ensure && \
#      rm -f schema/bindata.go && \
#      go generate ./schema
RUN CGO_ENABLED=0 GOOS=linux go build -a -o checker .

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=build /app/checker .
CMD ["/checker"]
