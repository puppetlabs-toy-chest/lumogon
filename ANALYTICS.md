# Anonymous Aggregate User Behaviour Analytics

Lumogon has begun gathering anonymous aggregate user behaviour analytics and reporting these to Google Analytics.

## Why?

TODO: Justification on why we're doing this

## What?
Lumogon's analytics record some shared information for every event:

- The Google Analytics version, i.e. `1` (https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters#v)
- The Lumogon analytics tracking ID, e.g. `UA-54263865-7` (https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters#tid)
- A Docker analytics user ID, e.g. `7TRN:IPZB:QYBB:VPBQ:UMPP:KARE:6ZNR:XE6T:7EWV:PKF4:ZOJD:TPYS`. This is generated by `docker`. This does not allow us to track individual users but does enable us to accurately measure user counts vs. event counts (https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters#cid)
- The Lumogon application name, e.g. `Lumogon` (https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters#an)
- The Lumogon application version, e.g. `1.0.0` (https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters#av)
- The Lumogon analytics hit type, e.g. `screenview` (https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters#t)
- The Lumogon analytics screen view, e.g. `list` (https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters#cd)


Lumogon's analytics records the following different events:

- a `screenview` hit type with the official Lumogon command you have run (with arguments stripped), e.g. `lumogon capabilities`
- an `event` type with any ancillary actions taken with a given command, e.g. user uploads to SaaS service

You can also view all the information that is sent by Lumogon's analytics by setting `LUMOGON_DISABLE_ANALYTICS=1` in your environment. Please note this will also stop any analytics from being sent.

It is impossible for the Lumogon developers to match any particular event to any particular user.


## When/Where?
Lumogon's analytics are sent throughout Lumogon's execution to Google Analytics over HTTPS.

## Who?

TODO: Document Who has access to this

## How?
The code is viewable in [analytics](https://github.com/puppetlabs/lumogon/blob/master/analytics/ga.go). They are done in a separate background process and fail fast to avoid delaying any execution. They will fail immediately and silently if you have no network connection.

## Opting out
Lumogon analytics helps us maintainers and leaving it on is appreciated. However, if you want to opt out of Lumogon's analytics, you can set this variable in your environment:

```sh
export LUMOGON_DISABLE_ANALYTICS=1
```

Alternatively, you can pass the following CLI argument at runtime

```sh
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock local/lumogon --disable-analytics
```