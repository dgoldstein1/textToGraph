FROM golang:1.9

# setup go
ENV GOBIN $GOPATH/bin
ENV PATH $GOBIN:/usr/local/go/bin:$PATH

COPY build $GOBIN
RUN crawler --version

COPY LICENSE /LICENSE
COPY VERSION /VERSION

ENV COMMAND "--help"
CMD textToGraph $COMMAND
