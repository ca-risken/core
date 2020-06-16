FROM golang AS builder
WORKDIR /go/src/github.com/CyberAgent/mimosa-core/
ADD *.go go.* ./
ENV CGO_ENABLED=0 GO111MODULE=on GOPRIVATE="github.com/CyberAgent/*"
ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN echo "machine github.com" > ~/.netrc
RUN echo "login $GITHUB_USER" >> ~/.netrc
RUN echo "password $GITHUB_TOKEN" >> ~/.netrc
RUN go install -v

FROM alpine
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /go/bin/gateway /usr/local/gateway/bin/
ENV PORT= \
    FINDING_SVC_ADDR= \
    TZ=Asia/Tokyo
WORKDIR /usr/local/gateway
CMD ["bin/gateway"]
