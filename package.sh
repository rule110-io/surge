#!/bin/bash

set -e

set -o allexport
source .env set +o allexport



# functions
requeststatus() { # $1: requestUUID
    requestUUID=${1?:"need a request UUID"}
    req_status=$(xcrun altool --notarization-info "$requestUUID" \
                              --username "${AC_USERNAME}" \
                              --password "${AC_PASSWORD}" 2>&1 \
                 | awk -F ': ' '/Status:/ { print $2; }' )
    echo "$req_status"
}

notarizefile() { # $1: path to file to notarize, $2: identifier
    filepath=${1:?"need a filepath"}
    identifier=${2:?"need an identifier"}
    
    # upload file
    echo "## uploading $filepath for notarization"
    requestUUID=$(xcrun altool --notarize-app \
                               --primary-bundle-id "$identifier" \
                               --username "${AC_USERNAME}" \
                               --password "${AC_PASSWORD}" \
                               --file "$filepath" 2>&1 \
                  | awk '/RequestUUID/ { print $NF; }')
                               
    echo "Notarization RequestUUID: $requestUUID"
    
    if [[ $requestUUID == "" ]]; then 
        echo "could not upload for notarization"
        exit 1
    fi
        
    # wait for status to be not "in progress" any more
    request_status="in progress"
    while [[ "$request_status" == "in progress" ]]; do
        echo -n "waiting... "
        sleep 10
        request_status=$(requeststatus "$requestUUID")
        echo "$request_status"
    done
    
    # print status information
    xcrun altool --notarization-info "$requestUUID" \
                 --username "${AC_USERNAME}" \
                 --password "${AC_PASSWORD}"
    echo 
    
    if [[ $request_status != "success" ]]; then
        echo "## could not notarize $filepath"
        exit 1
    fi
    
}



rm -rf ./build/bin

#sed "s/0.0.0/${VERSION}/" ./build/darwin/Info.plist.src > ./build/darwin/Info.plist
#CGO_LDFLAGS=-mmacosx-version-min=10.13 wails build -platform darwin/amd64 -o surge
#CGO_LDFLAGS=-mmacosx-version-min=10.13 wails build -platform darwin/arm64 -o surge
CGO_LDFLAGS=-mmacosx-version-min=10.13 wails build -platform darwin/universal -o surge

cd ./build/bin/

echo "Signing the binary..."
codesign -s "${SIGNING_IDENTITY}" -o runtime -v "./Surge.app/Contents/MacOS/Surge"

echo "Creating DMG..."
create-dmg surge.dmg ./Surge.app --overwrite --identity="${SIGNING_IDENTITY}" --dmg-title "Install Surge"
mv surge.dmg "surge.${VERSION}.dmg"

echo "TARing..."
tar -czvf surge.${VERSION}.tar.gz ./Surge.app

echo "Zipping..."
zip -r surge.zip ./Surge.app
mv surge.zip "surge.${VERSION}.zip"


echo "Notarizing..."

notarizefile "surge.${VERSION}.zip" "com.rule110.surge"
notarizefile "surge.${VERSION}.dmg" "com.rule110.surge"
xcrun stapler staple "surge.${VERSION}.dmg"

rm -rf ./build/bin/Surge.app

open .