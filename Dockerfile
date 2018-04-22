FROM scratch
COPY ari-proxy /
ENTRYPOINT ["/ari-proxy"]
