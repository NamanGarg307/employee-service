FROM golang:latest as builder


# # create a working directory
COPY . /employee-records-service
WORKDIR /employee-records-service

RUN go mod tidy

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o employee-records-service

# use a minimal alpine image
FROM alpine:latest

# add ca-certificates in case you need them
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
# set working directory

WORKDIR /root

# copy the binary from builder
COPY --from=builder /employee-records-service .
COPY app.rsa.pub .

EXPOSE 9876:9876
# run the binary

CMD ["./employee-records-service"]