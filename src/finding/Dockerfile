FROM golang:1.17 as builder
WORKDIR /go/src/github.com/ca-risken/core/src/finding/
ADD *.go go.* ./
ADD model/ ./model
ENV CGO_ENABLED=0 GOPRIVATE="github.com/ca-risken/*"
ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN echo "machine github.com" > ~/.netrc
RUN echo "login $GITHUB_USER" >> ~/.netrc
RUN echo "password $GITHUB_TOKEN" >> ~/.netrc
RUN go build && cp finding /go/bin/

FROM public.ecr.aws/risken/base/risken-base:v0.0.1
COPY --from=builder /go/bin/finding /usr/local/finding/bin/
ENV PORT= \
  DB_MASTER_HOST= \
  DB_MASTER_USER= \
  DB_MASTER_PASSWORD= \
  DB_SLAVE_HOST= \
  DB_SLAVE_USER= \
  DB_SLAVE_PASSWORD= \
  DB_SCHEMA=mimosa \
  DB_PORT=3306 \
  DB_LOG_MODE=false \
  DB_MAX_CONNECTION= \
  AWS_XRAY_CONTEXT_MISSING= \
  TZ=Asia/Tokyo
WORKDIR /usr/local/finding
ENTRYPOINT ["/usr/local/bin/env-injector"]
CMD ["bin/finding"]
