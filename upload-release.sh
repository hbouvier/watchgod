#!/bin/sh
if [ -z ${GITHUB_TOKEN} ] ; then
  echo "You must create a GITHUB TOKEN to upload the release"
  exit 0
fi

ID=$$
VERSION=$(cat VERSION.txt)
REGISTRY=hbouvier

GIT_COMMIT_REV:=$(shell git rev-parse HEAD)
GIT_COMMIT_REV_SHORT=$(git rev-parse --short HEAD)
GIT_COMMIT_MESSAGE_RAW=$(git log --oneline | head -1 | cut -f2- -d' ')
GIT_COMMIT_MESSAGE=$(echo ${GIT_COMMIT_MESSAGE_RAW} | tr -dc '[:print:]' | sed 's/"//g' | sed "s/'//g")

## Upload to GITHUB
echo "Preparing metadata for watchgod version ${VERSION} rev ${GIT_COMMIT_REV_SHORT} title: '${GIT_COMMIT_MESSAGE}'"
echo '{"tag_name":"v'${VERSION}'-'${GIT_COMMIT_REV_SHORT}'","name":"v'${VERSION}'-'${GIT_COMMIT_REV_SHORT}'","target_commitish":"'${GIT_COMMIT_REV}'","body":"'${GIT_COMMIT_MESSAGE}'","prerelease":true,"draft":false}' > /tmp/watchgod.${ID}.json
RELEASE=$(curl "https://api.github.com/repos/hbouvier/watchgod/releases" \
               -sXPOST \
               --header "Authorization: token ${GITHUB_TOKEN}" \
               --header "Content-Type: application/json; charset=UTF-8" \
               -d @/tmp/watchgod.${ID}.json)
RELEASE_ID=$(echo ${RELEASE} | jq '.id')
echo "\tRELEASE_ID=${RELEASE_ID}"

cd release/bin
for ARCH in $(find . -type d -print | grep -vE '^.$') ; do
  cd ${ARCH}
  echo "\tcompressing watchgod-${ARCH}-v${VERSION}-${GIT_COMMIT_REV_SHORT}.zip"
  zip -9 ../../watchgod-${ARCH}-v${VERSION}-${GIT_COMMIT_REV_SHORT}.zip watchgod
  cd ..
  ASSET=$(curl "https://uploads.github.com/repos/hbouvier/watchgod/releases/${RELEASE_ID}/assets?name=watchgod-${ARCH}-v${VERSION}-${GIT_COMMIT_REV_SHORT}.zip" \
              -sXPOST \
              --header "Authorization: token ${GITHUB_TOKEN}" \
              --header "Content-Type: application/zip" \
              --form "file=@./watchgod-${ARCH}-v${VERSION}-${GIT_COMMIT_REV_SHORT}.zip")
  echo "\tuploading release ${ASSET}"
done
cd ..
echo "\tOK"
