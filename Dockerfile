# Use the official Golang image as the base image
FROM golang:1.21

# 환경 변수 설정
ENV APP_PORT=8080

# 작업 디렉토리 설정
WORKDIR /app

# 필요한 파일을 컨테이너로 복사합니다.
COPY ./src .

# 애플리케이션을 빌드합니다.
RUN go build -o app

# 컨테이너에서 사용할 포트를 노출합니다.
EXPOSE 8080

# 애플리케이션을 실행합니다.
CMD ["./app"]
