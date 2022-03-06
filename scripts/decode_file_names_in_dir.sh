#!/usr/bin/env bash

# example:
# bash  scripts/decode_file_names_in_dir.sh ~/.stayback/20220207-050726
# /Users/nmarks/.1password : /Users/nmarks/.stayback/20220207-050726/L1VzZXJzL25tYXJrcy8uMXBhc3N3b3Jk.tar.gz.asc
# /Users/nmarks/.aws : /Users/nmarks/.stayback/20220207-050726/L1VzZXJzL25tYXJrcy8uYXdz.tar.gz.asc
# /Users/nmarks/.ssh : /Users/nmarks/.stayback/20220207-050726/L1VzZXJzL25tYXJrcy8uc3No.tar.gz.asc
# /Users/nmarks/Pictures : /Users/nmarks/.stayback/20220207-050726/L1VzZXJzL25tYXJrcy9QaWN0dXJlcw==.tar.gz.asc
# /Users/nmarks/bin : /Users/nmarks/.stayback/20220207-050726/L1VzZXJzL25tYXJrcy9iaW4=.tar.gz

# given /path/to/file/some_base64_string.tar.gz.asc
# return the decoded base64 string
function decode() {
  basename "${1}" | cut -f 1 -d '.' | base64 -d
}
for i in ${1}/*; do
    echo "$(decode ${i}) : ${i}"
done