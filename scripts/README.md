## sbcompress
Tarballs a list of directories. Each tarball base name is the base64 encode of the realpath of the original directory. The default destination is $HOME/.stayback. A temporary directory is used to create the tarballs. The tarballs replace existing tarballs of the same name only if the tgz process is successful.

## sbencrpyt
NOTE: sbencrypt can use [makemine](https://github.com/natemarks/makemine) to set a default recipient if it's used.

Uses gpg to encrypt a list of files - usually the tarballs created by sbcompress

First create and export/backup your gpg keys and passphrase

```shell
# gen-key will prompt for a passphrase . be sure to save it right away
gpg --gen-key
gpg --output my-public.pgp --armor --export my@email.com
gpg --output my-private.pgp --armor --export-secret-keys my@email.com

```

sbencrypt uses the gpg key to encrypt the file, but it also needs a recipient to use the public key to encrypt. It will look for the /etc/makemine/makemine.sh script. If it exists, it sources that script and uses the environment variable EMAIL as the default recipient. 


The encrypyted files can be decrypted using the private key:
```shell
gpg --output file.txt --decrypt file.txt.gpg
```

## Example
This is how I use it to backup and encrypt sensitive data 
```shell
# compress my aws and sssh config directories to ~/.stayback
bash scripts/sbcompress.sh ~/.ssh ~/.aws
# encrypt the new tarballs and append '.asc' extension
# NOTE: these are ascii-armor files
bash scripts/sbencrypt.sh $(find ~/.stayback/ -type f -name "*.tar.gz")
```