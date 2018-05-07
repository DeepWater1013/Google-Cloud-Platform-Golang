#!/bin/bash

set -e

export GOOGLE_APPLICATION_CREDENTIALS=$KOKORO_KEYSTORE_DIR/71386_golang-samples-kokoro-service-account
export GOLANG_SAMPLES_KMS_KEYRING=ring1
export GOLANG_SAMPLES_KMS_CRYPTOKEY=key1

TIMEOUT=25m

curl https://storage.googleapis.com/gimme-proj/linux_amd64/gimmeproj > /bin/gimmeproj && chmod +x /bin/gimmeproj;
gimmeproj version;
export GOLANG_SAMPLES_PROJECT_ID=$(gimmeproj -project golang-samples-tests lease $TIMEOUT);
if [ -z "$GOLANG_SAMPLES_PROJECT_ID" ]; then
  echo "Lease failed."
  exit 1
fi
echo "Running tests in project $GOLANG_SAMPLES_PROJECT_ID";
trap "gimmeproj -project golang-samples-tests done $GOLANG_SAMPLES_PROJECT_ID" EXIT

set -x

export GOLANG_SAMPLES_SPANNER=projects/golang-samples-tests/instances/golang-samples-tests

go version
date

if [[ -d /cache ]]; then
  time mv /cache/* .
  echo 'Uncached'
fi

# Re-organize files
export GOPATH=$PWD/gopath
target=$GOPATH/src/github.com/GoogleCloudPlatform
mkdir -p $target
mv github/golang-samples $target
cd $target/golang-samples

# Check use of Go 1.7 context package
! grep -R '"context"$' * || { echo "Use golang.org/x/net/context"; false; }

if [[ $KOKORO_BUILD_ARTIFACTS_SUBDIR = *"system-tests"* && -n $GOLANG_SAMPLES_GO_VET ]]; then
  echo "This test run will run end-to-end tests.";
  export GOLANG_SAMPLES_E2E_TEST=1
  export PATH="$PATH:/tmp/google-cloud-sdk/bin";
  ./testing/kokoro/configure_gcloud.bash;
fi

# Download imports.
GO_IMPORTS=$(go list -f '{{join .Imports "\n"}}{{"\n"}}{{join .TestImports "\n"}}' ./... | sort | uniq | grep -v golang-samples | grep -v golang.org/x/tools/imports)
time go get -u -v -d $GO_IMPORTS

# Manually clone golang.org/x/tools since it's incompatible with Go 1.6.
mkdir -p $GOPATH/src/golang.org/x;
pushd $GOPATH/src/golang.org/x;
if [ ! -d tools ]; then
  git clone https://go.googlesource.com/tools;
fi
popd;

# Pin golang.org/x/tools and go-sql-driver/mysql to support Go 1.6.
if go version | grep go1\.6\.; then
  pushd $GOPATH/src/github.com/go-sql-driver/mysql;
  git checkout v1.3;
  popd;

  pushd $GOPATH/src/golang.org/x/tools;
  git checkout 8e070db38e5c55da6a85c81878ab769bf5667848;
  popd;
fi

go get github.com/tebeka/go2xunit
go install golang.org/x/tools/imports;
go install -v $GO_IMPORTS

# Do the easy stuff before running tests. Fail fast!
if [ $GOLANG_SAMPLES_GO_VET ]; then
  diff -u <(echo -n) <(gofmt -d -s .)
  go vet ./...
fi

date

# Run all of the tests
OUTFILE=gotest.out
2>&1 go test -timeout $TIMEOUT -v ./... | tee $OUTFILE
$GOPATH/bin/go2xunit -fail -input $OUTFILE -output sponge_log.xml
