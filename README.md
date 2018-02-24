# procman

Process manager for docker container.
Support to start/stop the process via API.

# Example

Manage nginx process by procman.

```Dockerfile
ARG PROCMAN_VERSION=0.0.1
RUN curl -sL https://github.com/uphy/procman/releases/download/${PROCMAN_VERSION}/procman_linux_amd64.tar.gz | tar zx --strip 1 -C /usr/local/bin

CMD procman -p 81 -- nginx -g "daemon off;"
```