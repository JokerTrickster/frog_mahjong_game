FROM golang:1.22.1-alpine AS builder
WORKDIR /build/
RUN rm -rf ./*
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -ldflags '-s' -installsuffix cgo -o main .

# FINAL STAGE
FROM alpine:3.15 AS runner
WORKDIR /app/
RUN mkdir -p /logs && chmod 755 /logs
RUN apk add --no-cache curl unzip \
    && curl -L https://github.com/grafana/loki/releases/download/v2.6.1/promtail-linux-arm64.zip -o /promtail.zip \
    && unzip /promtail.zip -d /usr/local/bin \
    && mv /usr/local/bin/promtail-linux-arm64 /usr/local/bin/promtail \
    && chmod +x /usr/local/bin/promtail \
    && rm /promtail.zip
COPY --from=builder /build/main .
# 현재 위치와 /app 디렉토리 확인
RUN mkdir -p /etc/promtail && chmod 755 /etc/promtail
ARG LOKI_URL
ENV LOKI_URL=${LOKI_URL}
COPY ./promtail-config.yaml /etc/promtail/promtail-config.yaml
RUN chmod 755 /etc/promtail/promtail-config.yaml 
RUN sed -i "s|http://loki:3100/loki/api/v1/push|http://${LOKI_URL}:3100/loki/api/v1/push|" /etc/promtail/promtail-config.yaml
RUN echo "cat /etc/promtail/promtail-config.yaml" | cat /etc/promtail/promtail-config.yaml
RUN echo 'export $(strings /proc/1/environ | grep AWS_CONTAINER_CREDENTIALS_RELATIVE_URI)' >> /root/.profile
ARG VERSION
ENV VER $VERSION
EXPOSE 80
EXPOSE 3100
CMD ["sh", "-c", "promtail -config.file=/etc/promtail/promtail-config.yaml & ./main"]