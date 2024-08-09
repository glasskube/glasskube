# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# FROM gcr.io/distroless/static:debug-nonroot
FROM ubuntu
RUN ["apt-get", "update"]
RUN ["apt-get", "install", "ca-certificates", "-y"]
RUN ["update-ca-certificates"]
WORKDIR /
COPY glasskube /glasskube
USER 65532:65532
ENTRYPOINT ["/glasskube"]
