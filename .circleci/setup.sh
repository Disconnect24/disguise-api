curl -sSL https://sdk.cloud.google.com > /tmp/gcl && bash /tmp/gcl --install-dir=$HOME --disable-prompts
# Set up client secret for GAE later on.
echo "$CLIENT_SECRET" | base64 --decode > ${HOME}/client-secret.json
$HOME/google-cloud-sdk/bin/gcloud --quiet components update
$HOME/google-cloud-sdk/bin/gcloud auth activate-service-account --key-file $HOME/client-secret.json
$HOME/google-cloud-sdk/bin/gcloud config set project $GCLOUD_PROJECT
