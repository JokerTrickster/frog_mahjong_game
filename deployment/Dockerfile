# 빌드 단계
FROM golang:1.22.1-alpine AS builder
WORKDIR /build/

# 의존성 캐시 단계 추가
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

# 빌드 실행
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags '-s -w' -o main .

# 실행 단계
FROM alpine:3.15 AS runner
WORKDIR /app/

# 필수 패키지 설치 및 Promtail 다운로드
RUN apk add --no-cache curl unzip && \
    curl -L https://github.com/grafana/loki/releases/download/v2.6.1/promtail-linux-arm64.zip -o /promtail.zip && \
    unzip /promtail.zip -d /usr/local/bin && \
    mv /usr/local/bin/promtail-linux-arm64 /usr/local/bin/promtail && \
    chmod +x /usr/local/bin/promtail && \
    rm /promtail.zip && \
    apk del curl unzip

# 실행 파일 복사
COPY --from=builder /build/main .

# Promtail 설정 파일 복사
RUN mkdir -p /etc/promtail && chmod 755 /etc/promtail
COPY ./promtail-config.yaml /etc/promtail/promtail-config.yaml
RUN sed -i "s|http://loki:3100/loki/api/v1/push|http://${LOKI_URL}:3100/loki/api/v1/push|" /etc/promtail/promtail-config.yaml

# 환경 변수 및 실행 설정
ARG LOKI_URL
ENV LOKI_URL=${LOKI_URL}
EXPOSE 80
EXPOSE 3100
CMD ["sh", "-c", "promtail -config.file=/etc/promtail/promtail-config.yaml & ./main"]
