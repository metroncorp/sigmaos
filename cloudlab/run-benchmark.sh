#!/bin/bash

usage() {
  echo "Usage: $0 --command COMMAND [--vm VM]" 1>&2
}

VPC=""
VM=0
COMMAND=""
while [[ $# -gt 0 ]]; do
  case $1 in
  --vpc)
    shift
    shift
    ;;
  --command)
    shift
    COMMAND=$1
    shift
    ;;
  --vm)
    shift
    VM=$1
    shift
    ;;
  -help)
    usage
    exit 0
    ;;
  *)
    echo "Error: unexpected argument '$1'"
    usage
    exit 1
    ;;
  esac
done


if [ -z "$COMMAND" ] || [ $# -gt 0 ]; then
  usage
  exit 1
fi

DIR=$(dirname $0)
source $DIR/env.sh

vms=`cat servers.txt | cut -d " " -f2` 
vma=($vms)
MAIN="${vma[0]}"

SSHVM="${vma[$VM]}"

# Get the pubkey and privkey for the SigmaOS deployment
MASTER_PUB_KEY="$(ssh -i key-$VPC.pem ubuntu@$MAIN cat /tmp/sigmaos/master-key.pub)"
MASTER_PRIV_KEY="$(ssh -i key-$VPC.pem ubuntu@$MAIN cat /tmp/sigmaos/master-key.priv)"

echo "Run [$SSHVM]: $COMMAND"
ssh -i $DIR/keys/cloudlab-sigmaos $LOGIN@$SSHVM <<ENDSSH
  # Make sure swap is off on the benchmark machines.
  sudo swapoff -a
  cd sigmaos
  source ./env/env.sh
  echo "$MASTER_PRIV_KEY" > /tmp/sigmaos/master-key.priv 
  echo "$MASTER_PUB_KEY" > /tmp/sigmaos/master-key.pub
  export SIGMAPERF="KVCLERK_TPT;MRMAPPER_TPT;MRREDUCER_TPT;HOTEL_WWW_TPT;TEST_TPT;BENCH_TPT;THUMBNAIL_TPT;"
#  export SIGMAPERF="KVCLERK_TPT;MRMAPPER_TPT;MRREDUCER_TPT;HOTEL_WWW_TPT;TEST_TPT;BENCH_TPT;HOTEL_RESERVE_PPROF_MUTEX;CACHESRV_PPROF_MUTEX;"
  $COMMAND
ENDSSH
