FROM alpine:3.19
ARG TARGETOS
ARG TARGETARCH
COPY /bin/vt-manager-${TARGETOS}-${TARGETARCH}* /vt-manager
ENTRYPOINT [ "/vt-manager" ]