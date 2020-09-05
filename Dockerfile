FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
# RUN go get -u github.com/cosmtrek/air # Air doesn'r work since inotify is buggy from docker
RUN go get github.com/markbates/refresh
# RUN go build src/main.go