# stayback
MOTHBALLED - use duplicity
I liked this project for  some testability learning and more time spent with cobra.  check out the last commits for the best parts


stayback is a command line tool used for storing/retrieving sensitive linux desktop data in various AWS services.

sbs3  is for storing big bundles of directories in s3, optionally encrypted
sbsecrets is for packing a small local file into secrets manager
sbkp is for packing a local ssh key into EC2 keypair. this is a 1-way operation


sbs3 backup requirements: 
awscli v2
base64
realpath
tar
gzip/gunzip
gpg





```markdown
sbs3
Using job config /home/nmarks/.stayback/default.json
Local Destination: /my/usb/drive/backups (via /my/usb/drive/backups/.tmp/)
Uploading backup to S3 bucket: my-backup-bucket/stayback/20220105-080910
encrypted: /home/nmarks/.ssh -> LnNzaA==.tar.gz.gpg
encrypted: /home/nmarks/.aws -> LmF3cw==.tar.gz.gpg
unencrypted: /home/nmarks/Documents -> RG9jdW1lbnRz.tar.gz
Press 'c' to continue...

```


sbs3:
WARNING:if you back up the parent of a secure directory, it will defeat the purpose of backing up the secure directory with encryption.


 - manage a single local backup
 - upload the full contents of each backup run to s3

if my backup job  targets dir1, dir2 and dir3,a backup file is created for each in a local directory
if I run it again and target dir1 and dir4, the dir4 backup file will be added to the directory. the dir1 file will be replaced IFF the dir1 backup is successful



