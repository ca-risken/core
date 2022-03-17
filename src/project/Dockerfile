FROM golang:1.17 as builder
WORKDIR /go/src/github.com/ca-risken/core/src/project/
ADD *.go go.* ./
ADD model/ ./model
ENV CGO_ENABLED=0
ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN echo "machine github.com" > ~/.netrc
RUN echo "login $GITHUB_USER" >> ~/.netrc
RUN echo "password $GITHUB_TOKEN" >> ~/.netrc
RUN go build && cp project /go/bin/

FROM public.ecr.aws/risken/base/risken-base:v0.0.1
COPY --from=builder /go/bin/project /usr/local/project/bin/
ENV PORT= \
    PROFILE_EXPORTER= \
    PROFILE_TYPES= \
    IAM_SVC_ADDR= \
    DB_MASTER_HOST= \
    DB_MASTER_USER= \
    DB_MASTER_PASSWORD= \
    DB_SLAVE_HOST= \
    DB_SLAVE_USER= \
    DB_SLAVE_PASSWORD= \
    DB_SCHEMA=mimosa \
    DB_PORT=3306 \
    DB_LOG_MODE=false \
    TZ=Asia/Tokyo
WORKDIR /usr/local/project
ENTRYPOINT ["/usr/local/bin/env-injector"]
CMD ["bin/project"]
