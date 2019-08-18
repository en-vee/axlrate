# ================================
# Build Stage
# ================================
#FROM golang:latest AS builder
FROM axlrate-base:latest AS builder


RUN go get -u -v github.com/en-vee/aconf/...
RUN go get -u -v github.com/en-vee/alog/...
RUN mkdir -p /go/src/github.com/en-vee/axlrate
# Get source code from local or github (?)
COPY . /go/src/github.com/en-vee/axlrate
RUN ls -ltR /go/src/github.com/en-vee/axlrate
#RUN go get -u -v github.com/en-vee/axlrate/...

# Create target directories which can then be COPY'ed to the final image
RUN mkdir -p /opt/app/axlrate/bin
RUN mkdir -p /etc/opt/app/axlrate/conf

# Get dependencies
#RUN go get -d -v ./...

WORKDIR $GOPATH/src/github.com/en-vee/axlrate/app
# Build it as a static executable
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/axlrate .

RUN ls -ltr /go/bin/axlrate

# ================================
# Final Stage
# ================================
FROM alpine:latest

RUN apk update && apk add bash

# Create directories
COPY --from=builder /opt/app/axlrate/bin .
COPY --from=builder /etc/opt/app/axlrate/conf .

# Copy application binary
COPY --from=builder /go/bin/axlrate /opt/app/axlrate/bin/axlrate

# Copy application config files
COPY --from=builder /go/src/github.com/en-vee/axlrate/app/axlrate.conf /etc/opt/app/axlrate/conf/axlrate.conf
COPY --from=builder /go/src/github.com/en-vee/axlrate/app/alog.conf /etc/opt/app/axlrate/conf/alog.conf

# Set Env Variables
ENV ALOG_CONF_DIR /etc/opt/app/axlrate/conf/

#RUN ls -ltr /etc/opt/app/axlrate/conf/

#CMD ["sleep","9999"]
CMD ["/opt/app/axlrate/bin/axlrate","-config-file-name","/etc/opt/app/axlrate/conf/axlrate.conf"]

