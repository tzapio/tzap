FROM --platform=linux/amd64 ubuntu

WORKDIR /app
RUN apt-get update && apt-get install -y \
ca-certificates \
&& rm -rf /var/lib/apt/lists/*
COPY release/tzap-linux-amd64 /usr/local/bin/tzap
ENV TZAPEDITOR=api
RUN useradd -m -s /bin/bash user
USER user
ENTRYPOINT [ "tzap" ]