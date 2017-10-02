FROM astronomerio/docker-golang-alpine
MAINTAINER Ken Herner <ken@astronomer.io>

WORKDIR /go/src/github.com/astronomerio/s3proxy
COPY . .

RUN make build

FROM alpine
MAINTAINER Ken Herner <ken@astronomer.io>

COPY --from=0 /go/src/github.com/astronomerio/s3proxy/s3proxy /usr/local/bin/s3proxy

RUN apk --no-cache add ca-certificates

ENV GIN_MODE=release
EXPOSE 4041

ENTRYPOINT ["s3proxy"]