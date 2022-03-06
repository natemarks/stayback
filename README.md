# stayback
I need sensitive, important and convenient files like ssh keys and aws credential configs on all of my mac and linux machines.  
stayback is a convenient way to encrypt and store them in an S3 bucket and dowload them to whatever machine I'm on.

stayback is not a good backup tool. I'd look elsewhere for that.


## Usage
I create ~/.stayback.json:
```json
{
  "homeDirectory": "/Users/nmarks",
  "backupDirectory": "/Users/nmarks/.stayback",
  "s3Bucket": "com.imprivata.371143864265.us-east-1.personal",
  "recipient": "npmarks@gmail.com",
  "encryptedDirs": [
    ".aws",
    ".ssh",
    ".1password",
    "/Users/nmarks/Pictures"
  ],
  "unEncryptedDirs": [
    "bin"
  ]
}
```
When I run:
```shell
s3backup
```
It resolves my target directories to absolute paths from my $HOME directory.  It tarballs each one using the base64 value of the original absolute path as the file name, so the tarball for 

/home/natemarks/bin

would be:
L2hvbWUvbmF0ZW1hcmtzL2Jpbgo=.tar.gz

Why unreadable names? It was the easiest way I could come up with to keep the origin completely in the S3 data.

If the path is one of the 'encryptedDirs', s3backup also gpg encrypts the tarball using the recipient to identify the locally imported gpg public key.

Finally, it uploads the files to s3://s3Bucket/stayback/jobId


If I want to download the latest job data, I run:
```shell
s3restore
```

This finds the latest job ID in S3, creates a jobid directory (~/.stayback/jobId) sand copies the tarballs down to it and prints the origins for each file

## Requirements

### gpg is installed with a gpg key imported
make sure a recent version of  gpg is installed and generate a key 
```shell
gpg --full-generate-key
```
This doc may also help:
https://docs.github.com/en/authentication/managing-commit-signature-verification/generating-a-new-gpg-key

This new key lets gpg encrypt files for the recipient. You' ll also want to import the key on all of your 
machines so you cna decrypt files anywhere with just the key passphase.  The Ubuntu seahorse/nautilus  integration is handy for decrypting, too.

### Install these tools
awscli v2
base64
realpath
tar
gzip/gunzip


## Installation
download


## Configure
Create ~/.stayback.json using the example above