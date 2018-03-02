# email
library to send email messages, written in go. This library is using the standard library to make up all the moving parts in sending email messages. No external dependencies required.

## Use
To use the library you need to construct three parts:
* *SMTPServer* The smtp server the message is being submitted to
* *EmailUser* The authorized email user sending the message
* *Message* The message that is to be constructed and sent

## Status
This library is a work in progress but has been tested against a MS Exchange server and a Postfix server configured with starttls.

## Issues
* bcc not functioning as expected
* testing on exchange did not send messages outside of organization, but postfix worked sending to internal/external recipients
