FROM --platform=linux/amd64 ubuntu

WORKDIR /app

RUN apt-get update && apt-get install -y curl sudo
RUN curl https://tzap.io/sh | bash

ENTRYPOINT [ "tzap" ]