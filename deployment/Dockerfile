# 빌드 단계
FROM golang:1.24-alpine AS builder
WORKDIR /build/

# 의존성 캐시 단계 추가
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

# 빌드 실행
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags '-s -w' -o main ./src

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

# 로그 디렉토리 생성 및 권한 설정
RUN mkdir -p /logs && chmod 755 /logs

# 실행 파일 복사
COPY --from=builder /build/main .

# Promtail 설정 디렉토리 생성
RUN mkdir -p /etc/promtail && chmod 755 /etc/promtail

# Promtail 설정 파일 복사 (상대 경로 수정)
COPY promtail-config.yaml /etc/promtail/promtail-config.yaml

# 환경변수 설정
ARG LOKI_URL
ENV LOKI_URL=${LOKI_URL}

# Loki URL 업데이트 (LOKI_URL이 설정된 경우에만)
RUN if [ -n "$LOKI_URL" ]; then \
        sed -i "s|http://loki:3100/loki/api/v1/push|http://${LOKI_URL}:3100/loki/api/v1/push|" /etc/promtail/promtail-config.yaml; \
    fi

# AWS 자격 증명 설정
RUN echo 'export $(strings /proc/1/environ | grep AWS_CONTAINER_CREDENTIALS_RELATIVE_URI)' >> /root/.profile

# 버전 설정
ARG VERSION
ENV VER=$VERSION

# 포트 노출
EXPOSE 80
EXPOSE 3100

# 실행 명령
CMD ["sh", "-c", "promtail -config.file=/etc/promtail/promtail-config.yaml & ./main"]