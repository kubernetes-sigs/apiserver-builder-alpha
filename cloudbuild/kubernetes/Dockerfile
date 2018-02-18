FROM golang:1.9-stretch

RUN mkdir -p src/k8s.io
WORKDIR $GOPATH/src/k8s.io
RUN go get github.com/jteeuwen/go-bindata/go-bindata
RUN git clone https://github.com/kubernetes/kubernetes --depth=1 -b release-1.9

WORKDIR $GOPATH/src/k8s.io/kubernetes
RUN apt update
RUN apt install rsync -y

ENV CGO 0

RUN make clean
ENV WHAT cmd/kube-apiserver
RUN make
ENV WHAT cmd/kube-controller-manager
RUN make
ENV WHAT cmd/kubectl
RUN make

ENV KUBE_BUILD_PLATFORMS darwin/amd64
# RUN make clean
# RUN mkdir -p _output/bin
ENV WHAT cmd/kube-apiserver
RUN make
ENV WHAT cmd/kube-controller-manager
RUN make
ENV WHAT cmd/kubectl
RUN make

RUN mkdir $GOPATH/bin/darwin
RUN mkdir $GOPATH/bin/linux

RUN cp _output/local/bin/darwin/kube-apiserver $GOPATH/bin/darwin
RUN cp _output/local/bin/darwin/kube-controller-manager $GOPATH/bin/darwin
RUN cp _output/local/bin/darwin/kubectl $GOPATH/bin/darwin

RUN cp _output/local/bin/linux/kube-apiserver $GOPATH/bin/linux
RUN cp _output/local/bin/linux/kube-controller-manager $GOPATH/bin/linux
RUN cp _output/local/bin/linux/kubectl $GOPATH/bin/linux