FROM golang:1.19
WORKDIR /go/src/iothub
COPY ./ ./
ARG version
ARG gitCommit
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN CGO_ENABLED=1 go build -ldflags "-X main.Version=${version} -X main.GitCommit=${gitCommit}" -o tio /go/src/iothub/cmd/tio/main.go


FROM ubuntu:20.04

RUN apt update && \
  apt install tzdata && \
  rm -rf /var/lib/apt/lists/*
ENV TZ=Asia/Shanghai

WORKDIR /app/

COPY --from=0 /go/src/iothub/iothub .
COPY --from=0 /go/src/iothub/config.yaml .

CMD ./iothub
