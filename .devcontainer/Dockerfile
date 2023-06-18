FROM mcr.microsoft.com/devcontainers/base:bullseye

# Update and install required packages
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
        git \
        make \
        curl \
        htop \
        tar \
        gcc \
        ffmpeg
RUN wget https://github.com/tinygo-org/tinygo/releases/download/v0.27.0/tinygo_0.27.0_amd64.deb
RUN sudo dpkg -i tinygo_0.27.0_amd64.deb

RUN curl -sfL https://github.com/upx/upx/releases/download/v4.0.2/upx-4.0.2-amd64_linux.tar.xz | tar -xJvf -; \
sudo mv upx-4.0.2-amd64_linux/upx /usr/local/bin; \
rm -rf upx-4.0.2-amd64_linux/

RUN sudo rm -rf /usr/local/go && curl -sfL https://go.dev/dl/go1.20.3.linux-amd64.tar.gz | sudo tar -C /usr/local -xzf -
RUN sudo ln -s /usr/local/go/bin/go /usr/bin/go
RUN sudo ln -s /usr/local/go/bin/gofmt /usr/bin/gofmt

RUN wget -q -O tmp.zip https://github.com/protocolbuffers/protobuf/releases/download/v22.0/protoc-22.0-linux-x86_64.zip && \
    unzip tmp.zip -d /usr/local/ && \
    rm tmp.zip

RUN sudo mkdir -p /usr/local/lib/nodejs
ENV VERSION=v20.0.0
ENV DISTRO=linux-x64
RUN curl -sfL https://nodejs.org/dist/${VERSION}/node-${VERSION}-${DISTRO}.tar.xz | sudo tar -xJvf - -C /usr/local/lib/nodejs 

RUN sudo ln -s /usr/local/lib/nodejs/node-$VERSION-$DISTRO/bin/node /usr/local/bin/node
RUN sudo ln -s /usr/local/lib/nodejs/node-$VERSION-$DISTRO/bin/npm /usr/local/bin/npm
RUN sudo ln -s /usr/local/lib/nodejs/node-$VERSION-$DISTRO/bin/npx /usr/local/bin/npx

USER vscode
RUN mkdir -p /home/vscode/go/src/github.com/tzapio/tzap
WORKDIR /home/vscode/go/src/github.com/tzapio/tzap
RUN echo set enable-bracketed-paste off > ~/.inputrc
RUN echo "export PATH=/home/vscode/go/bin:/usr/local/lib/nodejs/node-$VERSION-$DISTRO/bin:\$PATH" >> /home/vscode/.bashrc
RUN echo 'alias tt="go run ./cli/main.go"' >> /home/vscode/.bashrc
# Install Go tools
RUN go install -v golang.org/x/tools/gopls@latest
RUN go install -v golang.org/x/tools/cmd/goimports@latest
RUN go install -v github.com/ramya-rao-a/go-outline@latest
RUN go install -v github.com/go-delve/delve/cmd/dlv@latest
RUN go install -v github.com/rogpeppe/godef@latest
RUN go install -v honnef.co/go/tools/cmd/staticcheck@latest
RUN go install -v github.com/josharian/impl@latest
RUN go install -v golang.org/x/tools/cmd/gorename@latest

RUN go install -v github.com/mitranim/gow@latest
RUN npm install -g live-server

USER root
# Cleanup
RUN apt-get clean autoclean && \
    apt-get autoremove -y && \
    rm -rf /var/lib/{apt,dpkg,cache,log}/

# Set user
USER vscode

# Print directory contents and current working directory
RUN ls -la && \
    pwd