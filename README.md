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
dev_appserver.py app.yaml
```

`dev_appserver.py` emulates App Engine locally. Check out the log, and go to its admin page. This can be very helpful for debugging issues with how mail is stored, i.e in Cloud Datastore.

We don't particularly mind you using our domain for testing, but if you are planning on forking us we do mind. Check out `config.sample.json`: our CircleCI configuration copies itself. What you're seeing as sample is what we do in production. If you're testing locally, you __will__ need to copy it over for yourself!

Also in `config.sample.json` is `mailinterval`: this is how many minutes we tell your Wii to space out the check-in for more mail.