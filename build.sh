#!/bin/bash
cd $(
	cd "$(dirname "$BASH_SOURCE")"
	cd -P "$(dirname "$(readlink "$BASH_SOURCE" || echo .)")"
	pwd
) || exit 1

rm -rf tmp
mkdir -p tmp
gcp -r -t ./tmp v2
gcp -r -t ./tmp/v2 LICENSE
gcp -r -t ./tmp/v2/pkg/plugins ../glauth-keycloak
rm -f rpmbuild/SOURCES/glauth.tar.gz
rm -f rpmbuild/SOURCES/go.mod
gcp -r v2/go.mod rpmbuild/SOURCES/go.mod
tar czf rpmbuild/SOURCES/glauth.tar.gz -C tmp v2

docker build -t glauth_build -f Dockerfile rpmbuild/ || exit 1
docker create --name glauth_build glauth_build || exit 1
docker cp glauth_build:/root/rpmbuild/RPMS/ ./ || exit 1
docker container rm glauth_build || exit 1
docker image rm glauth_build || exit 1
