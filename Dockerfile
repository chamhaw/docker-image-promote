# docker build --rm -f docker/Dockerfile -t plugins/docker .

FROM docker:18.03.0-ce-dind

ADD release/linux/amd64/image-promote /bin/image-promote
ENTRYPOINT ["/usr/local/bin/dockerd-entrypoint.sh", "/bin/image-promote"]
