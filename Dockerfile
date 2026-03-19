FROM alpine:3.23

COPY ai-server /

EXPOSE 8080

RUN mkdir -p /uploads

ENTRYPOINT ["/ai-server"]
