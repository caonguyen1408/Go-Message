FROM ubuntu
RUN apt-get update -y
RUN apt-get -y install git golang
RUN mkdir -p ${GOPATH}/go-rest-api; 
WORKDIR ${GOPATH}/go-rest-api/
COPY main.go ${GOPATH}/go-rest-api/
RUN go get -d -v ./...
RUN go build -o go-rest-api .
CMD [ "./go-rest-api", "daemon off;" ]
