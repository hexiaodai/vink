FROM alpine:3.18.4

ARG TARGETARCH

# Kubectl environment variable
ENV OS_ARCH="${TARGETARCH:-amd64}" \
    OS_NAME="linux" \
    KUBECTL_VERSION="v1.24.3" \
    HELM_VERSION="v3.17.1"

SHELL ["/bin/sh", "-o", "pipefail", "-c"]
# Install required system packages and dependencies
RUN apk update && \
    apk add bash && \
    apk add ca-certificates && \
    apk add curl && \
    apk add openssl

# Download kubectl binaries
RUN mkdir -p /tmp/pkg/cache/ && cd /tmp/pkg/cache/ && \
    if [ ! -e kubectl ]; then \
      curl -LO https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/${OS_NAME}/${OS_ARCH}/kubectl ; \
      curl -LO https://dl.k8s.io/${KUBECTL_VERSION}/bin/${OS_NAME}/${OS_ARCH}/kubectl.sha256 ; \
    fi && \
    # Note that there are 2 spaces between sha256 and kubectl files
    echo "$(cat kubectl.sha256)  kubectl" | sha256sum -c && \
    cp kubectl /usr/bin && \
    rm -rf kubectl kubectl.sha256

# Download helm binaries
RUN mkdir -p /tmp/pkg/cache/helm && cd /tmp/pkg/cache/helm && \
    curl -LO https://get.helm.sh/helm-${HELM_VERSION}-${OS_NAME}-${OS_ARCH}.tar.gz && \
    tar -zxf helm-${HELM_VERSION}-${OS_NAME}-${OS_ARCH}.tar.gz && \
    ls && \
    cp ${OS_NAME}-${OS_ARCH}/helm /usr/bin && \
    cd /tmp/pkg/cache && rm -rf helm

# Create files and authorizations
RUN mkdir /.kube && \
    chmod g+rwX /.kube && \
    chmod g+x /usr/bin/kubectl && \
    chmod g+x /usr/bin/helm

# USER 1001
