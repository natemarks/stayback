#!/usr/bin/env bash
# sbcompress.sh is a simple mockup of the general way sbs3 will compress files
# It works but it's intent is to just help me think the process through

# https://gist.github.com/natemarks/aebb7e84010d4bc37270d554106cb38b
set -Eeuo pipefail
trap cleanup SIGINT SIGTERM ERR EXIT

JOBID="$(date +%Y%m%d-%H%M%S)"
# shellcheck disable=SC2034
script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)

usage() {
  cat <<EOF
Usage: sbcompress.sh [-h] [-v] -d destination arg1 [arg2...]

compress a list of directories into a set of tarballs in a given destination directory. Each tarball base name is the
base64 encode of the realpath of the original directory

The default destination directory is $HOME/.stayback

Available options:

-h, --help        Print this help and exit
-v, --verbose     Print script debug info
-d, --destination Destination for backup tarball files
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
compress_dir() {
  ABSPATH=$(realpath "${1}")
  TAR_FILE=$(echo "${ABSPATH}" | base64).tar.gz
  msg "${GREEN}backing up ${ABSPATH} to ${destination}/${TAR_FILE}${NOFORMAT}"
  tar -czvf "${SBTMP}/${TAR_FILE}" "${ABSPATH}"
  rm -f "${destination}/${TAR_FILE}"
  mv "${SBTMP}/${TAR_FILE}" "${destination}"
  msg "${GREEN}SUCCESS: ${ABSPATH} -> ${destination}/${TAR_FILE}${NOFORMAT}"
}

check_dirs() {
  for i in "${args[@]}"
  do
    ls "${i}"
  done
}

# fail if any of the target directories don't exist
check_dirs
# ensure working dir
SBTMP="${destination}/${JOBID}"
mkdir -p "${SBTMP}"

# compress the target directories
for i in "${args[@]}"
do
  compress_dir "${i}"
done