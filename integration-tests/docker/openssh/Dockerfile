FROM alpine
RUN apk add --no-cache --quiet openssh
RUN echo "PermitEmptyPasswords yes" >> /etc/ssh/sshd_config \
 && echo "PermitRootLogin yes" >> /etc/ssh/sshd_config \
 && echo "root:" | chpasswd
EXPOSE 22
ADD entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
CMD ["/usr/sbin/sshd", "-D"]
