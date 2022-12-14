FROM chromedp/headless-shell:latest

# To reap zombie processeses
# Add Tini
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

# add ca-certificates
RUN \
	apt-get update -y \
    && apt-get install -y ca-certificates \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ADD bin/service /opt/fireglass/service
ADD configuration /tmp/configuration

ARG BUILD_COMMIT
ENV BUILD_COMMIT ${BUILD_COMMIT:-no-commit}
ARG BUILD_BRANCH
ENV BUILD_BRANCH ${BUILD_BRANCH:-no-branch}

CMD [ "./opt/fireglass/service" ]