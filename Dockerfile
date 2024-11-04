FROM alpine:3.20

COPY ai-server /

EXPOSE 8080

RUN mkdir -p /uploads

ENTRYPOINT ["/ai-server"]
