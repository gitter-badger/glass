# Payloads
(alpha 0 format version)


## Payload structure
* public metadata
* ... encryption with recipient's key
* ... signature with sender's key
* private metadata
  * creator app
  * timestamp
  * ...
* content

## Metadata format
probably XML, but why not MIME headers


## Payload processing status codes
* 0 to be processed
* 1 waiting for sender's key
* ...
* 100 processed
* error codes
  * -1 unsupported format version
  * -2 structural error
