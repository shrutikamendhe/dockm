FROM portainer/base

COPY dist app/

VOLUME /data

WORKDIR /

EXPOSE 9000

ENTRYPOINT ["/app/portainer-linux-amd64"]