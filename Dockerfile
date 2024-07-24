FROM docker.io/rockylinux/rockylinux:8-minimal

RUN microdnf update -y && microdnf install -y make rpmdevtools rpmlint dnf 'dnf-command(builddep)'
RUN microdnf install -y go-srpm-macros go{lang,lang-bin,-toolset}

# This isn't the real build but it will hopefully cache at least some of the build
RUN mkdir -p /root/v2/ && cd /root/v2/ && go mod init temp && go get github.com/gwelch-contegix/glauth/v2@master && go mod download && go build -ldflags "-s -w" -v -buildvcs=false -trimpath github.com/glauth/glauth/v2 || true
COPY / /root/rpmbuild
WORKDIR /root/rpmbuild/

RUN dnf builddep -y SPECS/glauth.spec
RUN spectool -g -R SPECS/glauth.spec
RUN rpmbuild -bb SPECS/glauth.spec
