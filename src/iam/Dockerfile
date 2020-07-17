FROM golang AS builder
WORKDIR /go/src/github.com/CyberAgent/mimosa-core/src/iam/
ADD *.go go.* ./
ENV CGO_ENABLED=0 GO111MODULE=on GOPRIVATE="github.com/CyberAgent/*"
ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN echo "machine github.com" > ~/.netrc
RUN echo "login $GITHUB_USER" >> ~/.netrc
RUN echo "password $GITHUB_TOKEN" >> ~/.netrc
RUN go install -v

FROM alpine
RUN apk add --no-cache ca-certificates tzdata \
  && cd /tmp/ && wget -q https://github.com/okzk/env-injector/releases/download/v0.0.5/env-injector_0.0.5_linux_amd64.tar.gz \
  && tar xvfz env-injector_0.0.5_linux_amd64.tar.gz && mv env-injector /usr/local/bin/ && rm /tmp/*
COPY --from=builder /go/bin/iam /usr/local/iam/bin/
ENV PORT= \
    DB_MASTER_HOST= \
    DB_MASTER_USER= \
    DB_MASTER_PASSWORD= \
    DB_SLAVE_HOST= \
    DB_SLAVE_USER= \
    DB_SLAVE_PASSWORD= \
    DB_SCHEMA=mimosa \
    DB_PORT=3306 \
    DB_LOG_MODE= \
    TZ=Asia/Tokyo
WORKDIR /usr/local/iam
ENTRYPOINT ["/usr/local/bin/env-injector"]
CMD ["bin/iam"]
