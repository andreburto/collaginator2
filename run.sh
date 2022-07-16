#!/usr/bin/env bash
export RANDOM_URL="https://afoyu5tqu4.execute-api.us-east-1.amazonaws.com/yvonne/random"
export IMAGE_FILE="yvonne.png"

cd "$(dirname $0)"

export BASE_DIR="$(pwd)"

cd $BASE_DIR

go build

# Get screen DIMENTIONS
OUTPUT=$(osascript -e 'tell application "Finder" to get bounds of window of desktop' | sed -e 's/,//g')
WIDTH=$(echo $OUTPUT | awk '{print $3}')
HEIGHT=$(echo $OUTPUT | awk '{print $4}')

# Create the image
$BASE_DIR/collaginator "$@" --width $WIDTH --height $HEIGHT --file $IMAGE_FILE

if [ $? != 0 ]; then exit $?; fi

# Set the image as wallpaper
osascript << EOL
tell application "Finder" to set desktop picture to "$BASE_DIR/$IMAGE_FILE" as POSIX file as alias
EOL

# Reset the dock to refresh the background.
killall Dock