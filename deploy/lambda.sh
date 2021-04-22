#!/usr/bin/env bash

echo "Deploying update to $FUNCTION_NAME"
echo "New code artifact is here: ${BUCKET}/${PROJECT_KEY}/${DRONE_COMMIT_SHA:0:7}.zip"

# Send stdout to /dev/null
# We'll check the exit code
# If there's an error, we'll see it
aws lambda update-function-code --function-name $FUNCTION_NAME --s3-bucket $BUCKET --s3-key $PROJECT_KEY/${DRONE_COMMIT_SHA:0:7}.zip --publish > /dev/null

if [ $? -ne 0 ]
then
 echo "Failed to update Lambda"
 exit 1
fi