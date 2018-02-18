FROM golang:stretch

RUN pwd
RUN ls -l -R
RUN ls -l -R /workspace

#RUN mkdir -p src/googlecloudplatform
#WORKDIR $GOPATH/src/googlecloudplatform
#RUN git clone https://github.com/kubernetes/code-generator --depth=1 -b release-1.9
#
#WORKDIR $GOPATH
#ENV CGO 0
#ENV GOARCH amd64
#
#ENV GOOS darwin
#RUN mkdir -p bin/darwin
#RUN go build -o bin/darwin/client-gen k8s.io/code-generator/cmd/client-gen
#RUN go build -o bin/darwin/conversion-gen k8s.io/code-generator/cmd/conversion-gen
#
#ENV GOOS linux
#RUN mkdir -p bin/linux
#RUN go build -o bin/linux/client-gen k8s.io/code-generator/cmd/client-gen
#RUN go build -o bin/linux/conversion-gen k8s.io/code-generator/cmd/conversion-gen
