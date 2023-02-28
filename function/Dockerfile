FROM golang:1.20.1
WORKDIR /entrypoint
COPY main.go .
COPY go.* .
RUN go build -ldflags "-s -w" -tags lambda.norpc -o entrypoint .

FROM logstash:8.6.2
COPY logstash.conf /usr/share/logstash/pipeline/logstash.conf
COPY logstash.yml /usr/share/logstash/config/logstash.yml
USER root
RUN find /usr/share/logstash -type f -exec chmod 777 {} \;
RUN find /usr/share/logstash -type d -exec chmod 777 {} \;
USER logstash
RUN mkdir -p /tmp/logstash/data
COPY --from=0 /entrypoint/entrypoint .
CMD ["./entrypoint"]
