FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /car-pooling-challenge

# This Dockerfile is optimized for go binaries, change it as much as necessary
# for your language of choice.

# RUN apk --no-cache add ca-certificates=20190108-r0 libc6-compat=1.1.19-r10

EXPOSE 9091

ENTRYPOINT [ "/car-pooling-challenge" ]