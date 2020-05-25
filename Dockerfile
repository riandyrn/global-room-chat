FROM golang:1.14.3-alpine3.11 as build
WORKDIR /go/src/github.com/riandyrn/global-room-chat

COPY . .
RUN go build -mod=vendor -o app

FROM alpine:3.11
RUN apk add ca-certificates &&\
    apk add -U tzdata

COPY --from=build /go/src/github.com/riandyrn/global-room-chat/ .
ENTRYPOINT [ "./app" ]