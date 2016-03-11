# Dockerfile for appsensor-reverse-proxy

FROM alpine:3.3

MAINTAINER John Melton <jtmelton@gmail.com>

# default values override with -e
ENV APPSENSOR_REST_ENGINE_URL=http://localhost:8085
ENV APPSENSOR_CLIENT_APPLICATION_ID_HEADER_NAME=X-Appsensor-Client-Application-Name
ENV APPSENSOR_CLIENT_APPLICATION_ID_HEADER_VALUE=reverse-proxy
ENV APPSENSOR_CLIENT_APPLICATION_IP_ADDRESS=127.0.0.1

# using ENV to allow dynamic loading of configuration files
ENV resource-verbs-mapping-file=testdata/sample-resource-verbs-mapping.yml
ENV resources-file=testdata/sample-resources.yml

# work around for permission error when ADD creates directories that binary is in
RUN mkdir /go && mkdir /go/bin

# add specially compiled binary directly
ADD appsensor-reverse-proxy /go/bin/proxy

# add config files
ADD $resource-verbs-mapping-file /tmp/resource-verbs-mapping.xml
ADD $resources-file /tmp/resources.yml

# cli args can be sent through as expected
ENTRYPOINT ["/go/bin/proxy", "-resource-verbs-mapping-file=/tmp/resource-verbs-mapping.xml", "-resources-file=/tmp/resources.yml"]

# if no cli args, default of -help is sent
CMD ["-help"]
# this is the default port to run appsensor-reverse-proxy on
EXPOSE 8080

