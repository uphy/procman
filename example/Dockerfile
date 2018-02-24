FROM nginx:latest

# Install curl to download procman binary
RUN apt-get update && \
    apt-get install -y curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install procman
ARG PROCMAN_VERSION=0.0.1
RUN curl -sL https://github.com/uphy/procman/releases/download/${PROCMAN_VERSION}/procman_linux_amd64.tar.gz | tar zx --strip 1 -C /usr/local/bin

CMD procman -p 81 -- nginx -g "daemon off;"
