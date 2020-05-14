# OpenDistro for Elasticsearch Alerts Handler

An HTTP server which used to handle webhooks triggered by [OpenDistro for Elasticsearch Alerting](https://opendistro.github.io/for-elasticsearch-docs/docs/alerting)

## Why?

As for time  of writing `destination` options that `ODFE` provides are limited.

1. It is not possible to send emails

## Features

- Ability to handle emails, and even send emails to multiple addresses within same webhook

## Usage


#### With Docker

    docker run --rm -p 8080:8080 youtous/odfe-alerts-handler --help

## Usage

    usage: odfe-alerts-handler [<flags>]

    Flags:
      -h, --help                   Show context-sensitive help (also try --help-long and --help-man).
          --web.listen-address=":8080"  
                                   Address to listen on for incoming HTTP requests.
          --smtp.host="localhost"  SMTP server hostname.
          --smtp.port=25           SMTP server port.
          --smtp.username=""       SMTP server login username.
          --smtp.password=""       SMTP server login password.
          --smtp.from="opendistro@localhost"    SMTP from address.
          --smtp.default-subject="Opendistro Alert fired"  
                                   SMTP default subject.
          --debug                  Enable debug mode.

    Env:
          LISTEN_ADDRESS           Address to listen on for incoming HTTP requests.
          SMTP_HOSTNAME            SMTP server hostname.
          SMTP_PORT                SMTP server port.
          SMTP_USERNAME            SMTP server login username.
          SMTP_PASSWORD            SMTP server login password.
          SMTP_FROM                SMTP from address.
          SMTP_DEFAULT_SUBJECT     SMTP default subject.
          DEBUG                    Enable debug mode.
    


## Configure ODFE Alerting destinations

### First

1. Go to `Alerting` > `Destinations`
2. Create the destination with type `Custom webhook`
3. Choose `Define endpoint by URL`
    - For `email` set the url to have path with `/email`, like `http://odfe-server:8080/email`
    
_(alternatively, useful if the frontend fail to validate hostname)_ query:
```
POST _opendistro/_alerting/destinations
{
  "name": "smtp",
  "type": "custom_webhook",
  "custom_webhook": {
    "url": "http://alerts-smtp-forwarder:8080/email"
  }
}

```

### Sending Email from triggers

1. Select destination which was created with the `/email` path
2. The `Message` body look like below:

```yaml
to: ['example@test.com', 'another@test.com']
subject: 'Optional subject param'
---
This is the body of the message
Here you can use the templeting as usual...
```

`subject` is optional, if not provided the default one used, see [usage](#usage).


## License

Licensed under the MIT License. See the LICENSE file for details.
