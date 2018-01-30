# disguise-api
a modern solution to an age-old situation

[![CircleCI](https://circleci.com/gh/Disconnect24/disguise-api/tree/master.svg?style=svg)](https://circleci.com/gh/Disconnect24/disguise-api/tree/master)


`disguise-api` is the nickname of the mail scripts. It will heavily be based off the authenticated (at the time of writing, master) fork of [Mail-Go](https://github.com/RiiConnect24/Mail-Go).

Instead of using MySQL, we'll be using Google Cloud Datastore.

## Awesome, how can I try it out?
You'll need the [Google Cloud SDK](https://cloud.google.com/sdk/) installed.
```terminal-session
git clone git@github.com:Disconnect24/disguise-api
cd disguise-api
dev_appserver.py app.yaml --default_gcs_bucket_name=<devel-bucket-name>
```

`dev_appserver.py` emulates App Engine locally. Check out the log, and go to its admin page.
This can be very helpful for debugging issues with how mail is stored, i.e in Cloud Datastore.
Feel free to omit the `default_gcs_bucket_name` argument if you're not planning on testing inbound parse -- it's required otherwise, as that's where inbound mail is stored.

We don't particularly mind you using our domain for testing, but if you are planning on forking us we do mind. Check out `config.sample.json`: our CircleCI configuration copies it and edits it. What you're seeing as sample is what we do in production. If you're testing locally, you __will__ need to copy it over for yourself!

Also in `config.sample.json` is `mailinterval`: this is how many minutes we tell your Wii to space out the check-in for more mail.


# Inner mechanics
The `frontend` directory contains the main web patcher, which generates credentials and modifies the given `nwc24msg.cfg` accordingly.

`inbound_parse.go` is in charge of handling inbound mail from SendGrid. It formulates a Wii Mail with attachments if needed, and stores in Google Cloud Storage (due to limited stored mail size). Because of this, and the fact that `send.go` has a specific stored format, "Bucketed" mail will have their `Body` as the path to get the mail's contents. `receive.go` handles this format as well, reading from Cloud Storage and serving right back out.