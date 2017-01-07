FROM alpine
RUN apk add --no-cache --quiet qemu-system-x86_64 qemu-img
RUN wget -q https://github.com/transprouter/sandbox-system/releases/download/v0.1/bzImage \
 && wget -q https://github.com/transprouter/sandbox-system/releases/download/v0.1/rootfs.qcow
EXPOSE 22
ADD entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
CMD qemu-system-x86_64 \
        -kernel bzImage -hda rootfs.qcow \
        -boot c -m 64 -append "root=/dev/sda console=ttyS0 rw" -localtime -no-reboot \
        -name sandbox \
        -device virtio-net,netdev=net0 -netdev user,id=net0 -redir tcp:22::22 \
        -nographic
