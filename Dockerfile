FROM golang:1.11
WORKDIR $GOPATH/src/systemapi
ADD . /go/src/systemapi
RUN go get github.com/gin-gonic/gin
RUN go get github.com/mackerelio/go-osstat/cpu
RUN go get github.com/mackerelio/go-osstat/memory
RUN go get github.com/mackerelio/go-osstat/network
RUN go get github.com/mackerelio/go-osstat/disk
RUN go install /go/src/systemapi
ENTRYPOINT /go/bin/systemapi
EXPOSE 8080
