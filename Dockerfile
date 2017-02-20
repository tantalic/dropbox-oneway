FROM scratch
MAINTAINER Kevin Stock <kevinstock@tantalic.com>

ADD certs/ca-certificates.crt /etc/ssl/certs/

# Labels: http://label-schema.org
ARG BUILD_DATE
ARG VCS_REF
ARG VERSION
LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.name="dropbox-oneway" \
      org.label-schema.description="One-way syncing files/directories from Dropbox" \
      org.label-schema.url="https://tantalic.com/dropbox-oneway" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/tantalic/dropbox-oneway" \
      org.label-schema.version=$VERSION \
      org.label-schema.schema-version="1.0"

ADD build/dropbox-oneway-linux_amd64 /dropbox-oneway

VOLUME /dropbox

ENV LOCAL_DIRECTORY /dropbox
# Other environment variables to set in child containers or via the CLI:
# ENV DROPBOX_TOKEN xxxxxxxxxxxxxxxxxxxxxxxxxx
# ENV DROPBOX_DIRECTORY /Path/in/Dropbox

CMD ["/dropbox-oneway"]
