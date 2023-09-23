
FROM alpine:latest

RUN apk update && apk add --no-cache libc6-compat
COPY ./shop /go/src/
COPY ./config.yaml /

CMD ["/go/src/shop"]