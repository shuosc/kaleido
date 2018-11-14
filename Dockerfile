FROM golang:1.11-alpine
RUN apk add git
COPY ./common /go/src/kaleido/common
COPY ./master /go/src/kaleido/master
COPY ./node /go/src/kaleido/node
WORKDIR /go/src/kaleido/node
RUN go get && go build
WORKDIR /go/src/kaleido/master
RUN go get && go build

FROM alpine
MAINTAINER longfangsong@icloud.com
COPY --from=0 /go/src/kaleido/node/node /
WORKDIR /
CMD ./node
EXPOSE 8080

FROM alpine
MAINTAINER longfangsong@icloud.com
RUN apk update && apk add ca-certificates
COPY --from=0 /go/src/kaleido/master/master /
COPY ./*.sql /
WORKDIR /
CMD ./master
EXPOSE 8086

FROM nginx:1.15.6-alpine
MAINTAINER longfangsong@icloud.com
COPY ./scope/dist /usr/share/nginx/html
