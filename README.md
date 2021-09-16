# About

http-file-upload is help utility/container to allow uploading files to a directory.

POST to /api/1.0/upload a JSON body with the following format:

```json
[
  {"name": "<filename>", "data": "<urlbase64 encoded file data>"},
  ...
]
```

# Configuring

http-file-upload is configured with environment variables:

Variable | Description
-------- | -----
TOKEN | **Required**. Should be passed in the Authentication header as `Bearer TOKEN`
ROOT | **Defaults** to current working directory. The file root for uploaded files. Uploads are "chrooted" to the directory and can't be written outside of this folder
LISTENADDR | **Defaults** to `:80`. Can be used to set the IP and port to listen on. Example: `10.0.0.1:8080`

# Container

korylprince/http-file-upload:<version>
