FROM google/golang
ADD *.go /gopath/botfounder/
ENV GIT_SSL_NO_VERIFY 1
ENV GOBIN /bin
WORKDIR /gopath/botfounder/
RUN go get
EXPOSE 3000
CMD ["/bin/botfounder"]
