# Hermes

Hermes is a message queue but unlike AMQP queue this one is machine to human. It takes a message with arbitrary
connection string and delivers it to a human using specialized backend.

The rationale is that time of notifications solely over email are gone. Today, each user prefers to receive
messages/notifications over different services such as Facebook Messenger, Whatsapp, Telegram, SMS and others.

The application is a proxy to such services and therefore it is intended for small and frequent notifications.

## Setup

```bash
go get github.com/katomaso/hermes
./hermes
```

Hermes runs on localhost and port 5587 by default. It accepts JSON messages over HTTP(S) protocol. It is highly
recommended to use HTTPS that can be turned on by `--https /path/to/certificate --key /path/to/key` where
the key's password (if there is any) will be taken from `$HERMES_KEY_PASS`.

Hermes searches for its configuration in multiple places and if it finds one it stops looking further:
`./hermes.conf`, `/etc/hermes.conf`, custom location given by `--config` parameter.

Every backend will try to activate and configure itself from keys in the configuration. For example the
most simple backend is emailing backend. It will look for configuration keys prefixed with `email`.

## Usage

Once hermes server is running, you can send messages that are supposed to be passed further.

```bash
curl -XPOST -d '{recipient:"Name <user@example.com>", subject:"Confirm your account", message:"<a href="click.me/1234">Click me</a>"}' localhost:5587
{"uuid": acb123}  # email is the default backend
curl -XPOST -d '{recipient:"whatsapp:+41123456789", subject:"Your order has arrived", message:"Pick it up in our shop any time"}' localhost:5587
{"uuid": acb124}  # uuid is returned in HTTP 201 response
curl -XPOST -d '{recipient:"facebook:nikita36", uuid:"xyz987", message:"Reservation confirmed. Have a nice stay."}' localhost:5587
{"uuid": xyz987}  # you can use your own UUID that must be unique within the queue
```

You can ask for results of message sending by the UUID that you have received when POSTing the message.
```bash
curl -XGET localhost:5587/status/abc123
{delivered: "", "retries": 3}
curl -XGET localhost:5587/status/abc124
{delivered: "2020-02-29 12:13:14.156+02:00", "retries": 0}
```
