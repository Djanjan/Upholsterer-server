FROM golang:latest
ADD . /go/src/go_service
WORKDIR /go/src/go_service
RUN go get -v -d
RUN go install go_service
ENTRYPOINT /go/bin/go_service
EXPOSE 8080