#!/usr/bin/env bash
# sbbackup.sh combines sbcompress and sb encrypt to place the encrypted backup files in an s3 bucket. It uses the
# recipient to choose the gpg public encryption key to use to encrypt files. it will default to the EMAIL value from
# makemine if it's in /etc/makemine/makemine.sh
# nNOTE: when I migrate this, I might want to split the task by destination so I cna run them in parallel AND clean up
# each destioation temp files as they complete

# https://gist.github.com/natemarks/aebb7e84010d4bc37270d554106cb38b
set -Eeuo pipefail
trap cleanup SIGINT SIGTERM ERR EXIT

JOBID="$(date +%Y%m%d-%H%M%S)"
# shellcheck disable=SC2034
script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)

usage() {
  cat <<EOF
Usage: sbbackup.sh [-h] [-v] -b bucket -d destination -r recipient arg1 [arg2...]

 - backup a list of directories into a set of tarballs in ~/.stayback/${JOBID}
 - encrypt the files in ~/.stayback/${JOBID} to *.tar.gz.asc
 - overwrite-move the unencrypted tarball files from ~/.stayback/${JOBID} to ~/.stayback
 - aws sync ${destination}/{$JOBID} to ${bucket}/stayback/${JOBID}
 - # cleanup ${destination}/{$JOBID} to ${bucket}/stayback/${JOBID}
base64 encode of the realpath of the original directory

The default destination directory is $HOME/.stayback

Available options:

-h, --help        Print this help and exit
-v, --verbose     Print script debug info
-b, --bucket      Backup bucket
-d, --destination Destination for backup tarball files
-r, --recipient   GPG public key recipient
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
  # default values of variables set from params
  destination="${HOME}/.stayback"

  while :; do
    case "${1-}" in
    -h | --help) usage ;;
    -v | --verbose) set -xv ;;
    --no-color) NO_COLOR=1 ;;
    -d | --destination) # example named parameter
      destination="${2-}"
      shift
      ;;
    -?*) die "Unknown option: $1" ;;
    *) break ;;
    esac
    shift
  done

  args=("$@")

  # check required params and arguments
  [[ ${#args[@]} -eq 0 ]] && die "Missing script arguments"

  return 0
}

parse_params "$@"
setup_colors


# script logic here
backup_dir() {
  ABSPATH=$(realpath "${1}")
  TAR_FILE=$(echo "${ABSPATH}" | base64).tar.gz
  msg "${GREEN}backing up ${ABSPATH} to ${destination}/${TAR_FILE}${NOFORMAT}"
  tar -cpzvf "${SBTMP}/${TAR_FILE}" "${ABSPATH}"
  rm -f "${destination}/${TAR_FILE}"
  mv "${SBTMP}/${TAR_FILE}" "${destination}"
  msg "${GREEN}SUCCESS: ${ABSPATH} -> ${destination}/${TAR_FILE}${NOFORMAT}"
}

check_dirs() {
  for i in "${args[@]}"
  do
    ls "${i}" > /dev/null
  done
}

# try to create the path to the backup bucket location. if it fails. this is a pretty good access test
# aws s3 sync  hhh/  s3://com.imprivata.371143864265.us-east-1.personal/stayback/123/
upload_job() {
  aws s3api put-object --bucket "${bucket}" --key "stayback/${JOBID}/"
  aws s3 sync ""  s3://com.imprivata.371143864265.us-east-1.personal/stayback/123/
}

# fail if any of the target directories don't exist
check_dirs
# ensure working dir
SBTMP="${destination}/${JOBID}"
mkdir -p "${SBTMP}"

# backup the target directories
for i in "${args[@]}"
do
  backup_dir "${i}"
done
