# 베이스 이미지 지정 - build 환경
FROM docker.io/library/golang:1.20 AS builder
RUN apt-get update && apt-get install -y tzdata

# 작업 디렉토리 설정
WORKDIR /app

# 의존성 (파일 복사, 다운로드)
COPY go.mod go.sum ./
RUN go mod tidy

# 애플리케이션 소스 코드 복사
COPY cmd/ ./cmd/
COPY database/ ./database/
COPY deploy/ ./deploy/

# 애플리케이션 컴파일
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o commerce-system ./cmd/main.go

# 베이스 이미지 지정 - 실행 환경
FROM docker.io/library/alpine:3.12.3
RUN apk add --no-cache tzdata

# 환경 변수로 시간대 설정
ENV TZ=Asia/Seoul
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 설정 및 실행 파일 복사
COPY --from=builder /app/commerce-system /usr/local/bin/commerce-system
COPY deploy/config.toml /configs/config.toml

# 시작 명령어 설정
ENTRYPOINT ["/usr/local/bin/commerce-system"]
CMD ["-c", "/configs/config.toml"]
