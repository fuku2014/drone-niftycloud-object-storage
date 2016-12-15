Use this plugin to upload files and build artifacts to an NIFTY Cloud  Object Storage bucket.

## Config

The following parameters are used to configure the plugin:

* **access_key** - niftycloud key
* **secret_key** - niftycloud secret key
* **bucket** - bucket name
* **region** - bucket region (`jp-east-2`, etc)
* **acl** - access to files that are uploaded (`private`, `public-read`, etc)
* **source** - source location of the files, using a glob matching pattern
* **target** - target location of files in the bucket
* **strip_prefix** - strip the prefix from source path
* **exclude** - glob exclusion patterns

## Example

Common example to upload to NIFTY Cloud Object Storage:

```yaml
pipeline:
  niftycloud:
    image: plugins/niftycloud-object-storage
    acl: public-read
    region: "jp-east-2"
    bucket: "my-bucket-name"
    access_key: "xxxxxxxxxxxxxxxxxxxx"
    secret_key: "xxxxxxxxxxxxxxxx"
    source: public/**/*
    strip_prefix: public/
    target: /target/location
    exclude:
      - **/*.xml
```
