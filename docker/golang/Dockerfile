FROM golang:1.9

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR $GOPATH

RUN go get github.com/golang/dep/... \
    && go install github.com/golang/dep/...

RUN go get github.com/DATA-DOG/godog/cmd/godog \
    && go install github.com/DATA-DOG/godog/cmd/godog

RUN go get gopkg.in/alecthomas/gometalinter.v1 \
    && go install gopkg.in/alecthomas/gometalinter.v1 \
    && gometalinter.v1 --install

