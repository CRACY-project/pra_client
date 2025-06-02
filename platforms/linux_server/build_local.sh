#!/bin/bash
set -e

ENVIRONMENT=$1
VERSION="v1.14.0"
LINUX_BUILD_DIR=platforms/linux_server/build
GIT_SHA="bd00e948fd51b58e2b2e63a2322c4ca5e0f62516"
TIMESTAMP=$(date +'%a, %b %d, %Y  %r')
FILE=go.mod

if [ ! -f "$FILE" ]; then
    echo "Please run script from main directory"
    exit 1
fi

# Default to testing environment if no argument is provided
if [ -z "$ENVIRONMENT" ]; then
    ENVIRONMENT="testing"
fi
echo "Environment is set to: $ENVIRONMENT"

# Replace any `/` in the version with `-`
VERSION=$(echo "$VERSION" | sed 's/\//-/g')

# Get the version number from the tag
INSTALLER_VERSION=$(echo $VERSION | cut -d "-" -f1 | cut -d "v" -f2)

# Check if INSTALLER_VERSION starts with a digit and remove all characters except digits and periods
if ! echo "$INSTALLER_VERSION" | grep -qE '^[0-9]'; then
    INSTALLER_VERSION="0.${INSTALLER_VERSION}"
fi
INSTALLER_VERSION=$(echo "$INSTALLER_VERSION" | sed 's/[^0-9.]//g')

# Build the executable
CGO_ENABLED=0 go build -tags "server $ENVIRONMENT" -o jimberfw_launcher launcher/main.go
CGO_ENABLED=0 go build -o jimberfw -tags "server $ENVIRONMENT" -ldflags="-X 'github.com/jimbersoftware/jimberfw_server/environment.environment=$ENVIRONMENT' -X 'github.com/jimbersoftware/jimberfw_server/environment.version=$VERSION' -X 'github.com/jimbersoftware/jimberfw_server/environment.buildDate=$TIMESTAMP' -X 'github.com/jimbersoftware/jimberfw_server/environment.shaCommitHash=$GIT_SHA'" cmd/server/main.go

echo "$LINUX_BUILD_DIR exists."
rm -Rf $LINUX_BUILD_DIR/*
mkdir -p $LINUX_BUILD_DIR/DEBIAN
sed "s/VERSION_PLACEHOLDER/${INSTALLER_VERSION}/" platforms/linux_server/deb/control >platforms/linux_server/build/DEBIAN/control
chmod +x platforms/linux_server/deb/postinst
cp platforms/linux_server/deb/postinst platforms/linux_server/build/DEBIAN

chmod +x platforms/linux_server/deb/preinst
cp platforms/linux_server/deb/preinst platforms/linux_server/build/DEBIAN

mkdir -p $LINUX_BUILD_DIR/etc/systemd/system/

cp platforms/linux_server/deb/jimbernetworkisolation.service $LINUX_BUILD_DIR/etc/systemd/system/
mkdir -p $LINUX_BUILD_DIR/usr/local/bin
cp jimberfw $LINUX_BUILD_DIR/usr/local/bin/jimberfw
cp jimberfw_launcher $LINUX_BUILD_DIR/usr/local/bin/jimberfw_launcher

mkdir -p installer/
dpkg-deb -Zxz --build $LINUX_BUILD_DIR
mv platforms/linux_server/build.deb installer/"${ENVIRONMENT}"_JimberNetworkIsolation.deb

SHA1_HASH=$(sha1sum ${LINUX_BUILD_DIR}/usr/local/bin/jimberfw | cut -d " " -f1)
echo -n "${SHA1_HASH}" >$LINUX_BUILD_DIR/jimberfw.sha1

touch platforms/linux_server/build-info.txt
cat >platforms/linux_server/build-info.txt <<EOF
Version: ${VERSION}
Tag: ${VERSION}
Git SHA hash: ${GIT_SHA}
Binary SHA1 hash: ${SHA1_HASH}
Timestamp: ${TIMESTAMP}
EOF

cat platforms/linux_server/build-info.txt
