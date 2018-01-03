# s3-du
Like the unix command `du` but for s3.

## Installation
There are two ways you can install s3-du:

  - From the [release page](https://github.com/tehmoon/s3-du/releases)
  - From the source -- requires Go:
```
git clone https://github.com/tehmoon/s3-du
cd s3-du
go get ./...
go build # A binary name s3-du will be generated in the directory
```

## Example

```
s3-du -b blih -d 0
```

## Usage
```
Usage of ./s3-du:
  -b string
    	Bucket to fetch keys from
  -d uint
    	Calculate directory sizes with specified depth
  -p string
    	Prefix for s3 object keys
```

## S3 Credentials
It uses the `S3` official SDK for `Go`, so you can use the same credential options as from `awscli` for example.

Supported environment variables:
  - AWS_ACCESS_KEY_ID
  - AWS_SECRET_ACCESS_KEY
  - AWS_PROFILE
  - AWS_CONFIG_FILE
  - AWS_SHARED_CREDENTIALS_FILE

You'll also need those access in order for the tool to work:
```
       {
            "Effect": "Allow",
            "Action": [
                "s3:ListBucketVersions"

            ],
            "Resource": [
                "arn:aws:s3:::*"
            ]
        }
```

## Caveats
  - If you have files and directories inside a directory, when the depth is greated than where the directory is, the size of the directory is the sum of all the regular files, not the regular files and its children.
  - Human readable option `-h` is to be implemented
