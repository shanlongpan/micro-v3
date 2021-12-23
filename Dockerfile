FROM centos:centos7

LABEL MAINTAINER=shanlongpan

COPY micro-v3-learn /data/

ENTRYPOINT ["/data/micro-v3-learn"]
