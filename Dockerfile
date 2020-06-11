FROM golang:1.14 AS builder

# Build only for linux-amd64
ENV GO111MODULE=on \
    GO_BUILD_OS=linux \
    GO_BUILD_ARCH=amd64

WORKDIR /src
COPY . .
RUN make build


FROM ubuntu:18.04
WORKDIR /fringe-runner

# Update the image and install dependencies
ENV GO_BUILD_OS=linux \
    GO_BUILD_ARCH=amd64

RUN apt-get update -y && apt-get upgrade -y \
    ca-certificates \
    nmap

# Copy the runner and download the modules
COPY --from=builder /src/build/fringe-runner-${GO_BUILD_OS}-${GO_BUILD_ARCH} ./fringe-runner
COPY --from=builder /src/config.yml /fringe-runner/config.yml
RUN ./fringe-runner update

ENTRYPOINT [ "./fringe-runner" ]
