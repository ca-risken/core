FROM golang:1.18.2 as builder
WORKDIR /src/

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /bin/core

FROM public.ecr.aws/risken/base/risken-base:v0.0.1
COPY --from=builder /bin/core /usr/local/bin/
ENV PORT= \
  PROFILE_EXPORTER= \
  PROFILE_TYPES= \
  TRACE_DEBUG= \
  FINDING_SVC_ADDR= \
  PROJECT_SVC_ADDR= \
  DB_MASTER_HOST= \
  DB_MASTER_USER= \
  DB_MASTER_PASSWORD= \
  DB_SLAVE_HOST= \
  DB_SLAVE_USER= \
  DB_SLAVE_PASSWORD= \
  DB_SCHEMA=mimosa \
  DB_PORT=3306 \
  DB_LOG_MODE=false \
  NOTIFICATION_ALERT_URL= \
  TZ=Asia/Tokyo
WORKDIR /usr/local/
ENTRYPOINT ["/usr/local/bin/env-injector"]
CMD ["bin/core"]