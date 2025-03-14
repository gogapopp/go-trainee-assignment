FROM golang:1.24

WORKDIR ${GOPATH}/avito-shop/
COPY . ${GOPATH}/avito-shop/

RUN go build -o /build ./cmd/avito-shop-service \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]