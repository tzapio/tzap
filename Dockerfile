FROM --platform=linux/amd64 ubuntu

RUN apt-get update && apt-get install -y \
ca-certificates \
&& rm -rf /var/lib/apt/lists/*
COPY release/tzap-linux-amd64 /usr/local/bin/tzap
ENV TZAPEDITOR=api
RUN useradd -m -s /bin/bash user
RUN mkdir /app
RUN chown -R user:user /app
RUN chmod -R 777 /app
USER user
WORKDIR /app
ENTRYPOINT [ "tzap" ]