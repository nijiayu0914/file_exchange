FROM golang:1.14.4 as build

ENV GO111MODULE on
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /go/cache

ADD go.mod .
ADD go.sum .
RUN go mod download
WORKDIR /go/release
ADD . .
RUN go build -ldflags="-s -w" -o compile
FROM scratch as prod
COPY --from=build /go/release/compile /
EXPOSE 8080
ENTRYPOINT ["/compile", "-config_type=http", "-config_path=http://192.168.0.130:7000/file_exchange/prod/config.yml"]
