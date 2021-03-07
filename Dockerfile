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
ADD /config/prod/config.yml /config/
RUN go build -ldflags="-s -w" -o compile
FROM scratch as prod
COPY --from=build /go/release/compile /
COPY --from=build /go/release/config/prod/config.yml ./config/config.yml
EXPOSE 8080
CMD ["/compile"]
