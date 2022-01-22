# sbcompress
compress relative or absolute directory to a tarball such that when extracted to a target directory, it creates the
original target directory withouth the parent tree

```shell
cd /home/myaccount
sbcompress .ssh
# echo /home/myaccount/.ssh | base64 -> L2hvbWUvbXlhY2NvdW50Ly5zc2gK
#  tar -czvf /home/myaccount/.stayback/.tmp/L2hvbWUvbXlhY2NvdW50Ly5zc2gK.tar.gz  /home/myaccount/.ssh
# if succeessful, mv /home/myaccount/.stayback/.tmp/L2hvbWUvbXlhY2NvdW50Ly5zc2gK.tar.gz to
# /home/myaccount/.stayback/L2hvbWUvbXlhY2NvdW50Ly5zc2gK.tar.gz overwriting  any previous backup of /home/myaccount/.ssh
```

```shell
/usr/bin/env bash
TMPDIR="$HOME/.stayback/.tmp"
mkdir -p "$TMPDIR"
TAR_FILE="$(realpath $1) | base64
```