FROM quay.io/centos/centos:stream8

RUN dnf update -y && dnf install -y make rpmdevtools rpmlint 'dnf-command(builddep)'
RUN dnf install -y go-srpm-macros go{lang-bin,lang,-toolset} --allowerasing

COPY /SOURCES/go.mod /root/v2/go.mod
# This isn't the real build but it will hopefully cache at least some of the build
RUN cd  /root/v2/ && go get github.com/glauth/glauth/v2@master && go mod download && go build -ldflags "-s -w" -v -buildvcs=false -trimpath github.com/glauth/glauth/v2 || true
COPY / /root/rpmbuild
WORKDIR /root/rpmbuild/
# RUN rm /root/v2
RUN dnf builddep -y SPECS/glauth.spec
RUN spectool -g -R SPECS/glauth.spec
RUN rpmbuild -bb SPECS/glauth.spec
