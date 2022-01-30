#!/usr/bin/env bash
# sbencrypt.sh encrypts a list of files, generally tarballls created by sbscompress.sh

# https://gist.github.com/natemarks/aebb7e84010d4bc37270d554106cb38b
set -Eeuo pipefail
trap cleanup SIGINT SIGTERM ERR EXIT

# shellcheck disable=SC2034
script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)


# if /etc/makemine/makemine.sh exists, source it
if [ -f "/etc/makemine/makemine.sh" ]; then
  # shellcheck disable=SC1091
  . /etc/makemine/makemine.sh
fi

usage() {
  cat <<EOF
Usage: sbencrypt.sh [-h] [-v] -r recipient arg1 [arg2...]

encrypt a list of files to the same directory as each target file, overwriting an existing gpg files

Available options:

-h, --help        Print this help and exit
-v, --verbose     Print script debug info
-r, --recipient   Recipient email used to choose the public key to encrypt the file
EOF
  exit
}


cleanup() {
  trap - SIGINT SIGTERM ERR EXIT
}

setup_colors() {
  if [[ -t 2 ]] && [[ -z "${NO_COLOR-}" ]] && [[ "${TERM-}" != "dumb" ]]; then
    NOFORMAT='\033[0m' RED='\033[0;31m' GREEN='\033[0;32m' ORANGE='\033[0;33m' BLUE='\033[0;34m' PURPLE='\033[0;35m' CYAN='\033[0;36m' YELLOW='\033[1;33m'
  else
    # shellcheck disable=SC2034
    NOFORMAT='' RED='' GREEN='' ORANGE='' BLUE='' PURPLE='' CYAN='' YELLOW=''
  fi
}

msg() {
  echo >&2 -e "${1-}"
}

die() {
  local msg=$1
  local code=${2-1} # default exit status 1
  msg "$msg"
  exit "$code"
}

parse_params() {
  # makemine.sh should export EMAIL. check and use it as the default if it is set
  recipient="${EMAIL}"

  while :; do
    case "${1-}" in
    -h | --help) usage ;;
    -v | --verbose) set -xv ;;
    --no-color) NO_COLOR=1 ;;
    -r | --recipient) # example named parameter
      recipient="${2-}"
      shift
      ;;
    -?*) die "Unknown option: $1" ;;
    *) break ;;
    esac
    shift
  done

  args=("$@")

  # check required params and arguments
  [[ -z "${recipient}" ]] && die "No recipient provided"
  [[ ${#args[@]} -eq 0 ]] && die "Missing script arguments"

  return 0
}

parse_params "$@"
setup_colors



# --openppg is required. If I don't use it, files encrypted on mac won't be decrypted on linux
encrypt_file() {
  msg "${GREEN}encrypting:  ${1}${NOFORMAT}"
  gpg --openpgp --batch --yes --output \
  "${1}.gpg" --encrypt --recipient "${recipient}" "${1}"
  msg "${GREEN}SUCCESS: encrypted ${1}${NOFORMAT}"
}


# check all of the target files for missing targets so they can all be addressed at once
check_files() {
  status=0
  for i in "${args[@]}"
  do
    if test -f "${i}"; then
      continue
    else
      msg "${RED}File doesn't exist: ${i}${NOFORMAT}"
      status=1
    fi
  done
  if [ $status -eq 1 ];then
    die "Missing target files"
  fi
}

msg "${GREEN}encrypting for recipient:  ${recipient}${NOFORMAT}"

# fail if any of the target directories don't exist
check_files

# encrypt the target directories
for i in "${args[@]}"
do
  encrypt_file "${i}"
done