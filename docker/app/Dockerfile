FROM golang:1.25

ENV GOPRIVATE=github.com/MousaZa/*

ARG GITHUB_TOKEN

RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

RUN go install github.com/cespare/reflex@v0.3.1
COPY reflex.conf /

COPY start.sh /
COPY entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
