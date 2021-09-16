FROM golang:1.17-alpine as builder

ARG VERSION

RUN apk add --no-cache git

RUN git clone --branch "$VERSION" --single-branch --depth 1 \
    https://github.com/korylprince/http-file-upload.git  /go/src/github.com/korylprince/http-file-upload

RUN cd /go/src/github.com/korylprince/http-file-upload && \
    go install -mod=vendor github.com/korylprince/http-file-upload


FROM alpine:3.13

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/http-file-upload /

CMD ["/http-file-upload"]
