#!/bin/bash
set -e

FILE=go.mod

if [ ! -f "$FILE" ]; then
    echo "Please run script from main directory"
    exit 1
fi

LINUX_BUILD_DIR=platforms/linux_server/build

ENVIRONMENT=$1

VERSION=$2
GIT_SHA=`git rev-parse HEAD`

TIMESTAMP=$(date +'%a, %b %d, %Y  %r')

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
CGO_ENABLED=0 go build -tags "$1" -o jimberpra_launcher launcher/main.go

echo "ENVIRONMENT=$ENVIRONMENT"
echo "VERSION=$VERSION"
echo "TIMESTAMP=$TIMESTAMP"
echo "GIT_SHA=$GIT_SHA"

CGO_ENABLED=0 go build -o jimberpra -tags "$1" \
    -ldflags="-X github.com/jimbersoftware/pra_client/environment.environment=${ENVIRONMENT} \
              -X github.com/jimbersoftware/pra_client/environment.version=${VERSION} \
              -X 'github.com/jimbersoftware/pra_client/environment.buildDate=${TIMESTAMP}' \
              -X github.com/jimbersoftware/pra_client/environment.shaCommitHash=${GIT_SHA}" \
    main.go

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
cp jimberpra $LINUX_BUILD_DIR/usr/local/bin/jimberfw
cp jimberpra_launcher $LINUX_BUILD_DIR/usr/local/bin/jimberfw_launcher

dpkg-deb -Zxz --build $LINUX_BUILD_DIR
mv platforms/linux_server/build.deb platforms/linux_server/JimberPraClient.deb

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

if [ "$DEPLOY" = "no" ]; then
    exit
fi

PLATFORM="server-linux"
# because bad code creates bad code, I have no time to rework this now
CLIENTPLATFORM="linux-server"
# PLATFORMPATH="linux_server"

# if [ "$ENVIRONMENT" = "dircon" ]; then
#     IP="104.248.202.67"
#     USER="jimber"
#     LABEL="dircon-"
# fi

# if [ "$ENVIRONMENT" = "testing" ]; then
#     IP="185.69.165.101"
#     USER="jimber"
#     LABEL="testing-"
# fi

# if [ "$ENVIRONMENT" = "staging" ]; then
#     IP="159.65.201.221"
#     USER="jimber"
#     LABEL="staging-"
# fi

# if [ "$ENVIRONMENT" = "production" ]; then
#     IP="signal.jimber.io"
#     USER="builduser"
#     LABEL=""
# # fi

# scp jimberfw ${USER}@${IP}:/var/www/binaries/${VERSION}_${PLATFORM}
# scp $LINUX_BUILD_DIR/jimberfw.sha1 ${USER}@${IP}:/var/www/binaries/${VERSION}_${PLATFORM}.sha1
# scp platforms/${PLATFORMPATH}/JimberNetworkIsolation.deb ${USER}@${IP}:/var/www/clients/${LABEL}${CLIENTPLATFORM}-${VERSION}.deb

# if [ ! "$ENVIRONMENT" = "production" ] && [ "$dircon" != "true" ]; then
#     scp platforms/${PLATFORMPATH}/JimberNetworkIsolation.deb ${USER}@${IP}:/var/www/clients/${LABEL}${CLIENTPLATFORM}-latest.deb
# fi
