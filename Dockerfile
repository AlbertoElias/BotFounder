FROM google/golang
ADD *.go /gopath/botfounder/
ADD templates /gopath/botfounder/templates
ENV PORT 3000
ENV GIT_SSL_NO_VERIFY 1
ENV GOBIN /bin
WORKDIR /gopath/botfounder/
RUN go get
RUN go build
EXPOSE 3000
CMD ["/gopath/botfounder/botfounder"]
