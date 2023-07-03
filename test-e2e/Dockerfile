# Basisimage definieren
FROM ubuntu:latest

# Umgebungsvariablen setzen
ARG OPENAI_APIKEY
ENV OPENAI_APIKEY=${OPENAI_APIKEY}
RUN if [ -z "$OPENAI_APIKEY" ]; then echo "OPENAI_APIKEY is not set or is an empty string. Exiting." && exit 1; fi

# System aktualisieren und notwendige Pakete installieren
RUN apt-get update && apt-get install -y \
    python3.10 \
    python3-pip \
    curl \
    build-essential \
    libssl-dev \
    libffi-dev \
    python3-dev \
    && curl -sL https://deb.nodesource.com/setup_16.x | bash - \
    && apt-get install -y nodejs


RUN npm install -g tzap@0.8.6
RUN pip install poetry==1.5.1

WORKDIR /app

# copy dependencies
RUN  touch ./README.md # required by poetry somehow
COPY pyproject.toml poetry.lock ./
COPY e2e_examples ./e2e_examples
COPY tests ./tests

# set up poetry
RUN poetry check
RUN poetry install

COPY Makefile ./Makefile
