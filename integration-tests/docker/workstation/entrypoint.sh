#!/bin/sh

set -e

if [ $(qemu-img snapshot -l rootfs.qcow | wc -l) -eq 0 ] ; then
  echo "Create snapshot of initial rootfs state"
  qemu-img snapshot -c initial rootfs.qcow
else
  echo "Restore initial rootfs state"
  qemu-img snapshot -a initial rootfs.qcow
fi

exec "$@"
