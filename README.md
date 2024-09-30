# commerce-system

<br>

### 실행 방법
1. `compose.yml` 파일이 존재하는 경로로 이동합니다.


2. `docker-compose up --build -d` 명령어를 실행합니다.


3. `Docker Desktop`에서 아래 이미지와 같이 컨테이너가 잘 실행 중인지 확인합니다.

<img src="readme/image/docker-desktop.png" width="800"/>

<br><br>

### 테스트 코드 실행 시키기 (Windows Powershell 기준)

📌 **모든 테스트 코드 확인 명령어**
```
go test ./internal/...
```

<br>

📌 **성공한 테스트 코드 확인 명령어**
```
go test ./internal/... -json | Select-String -Pattern '"Action":"pass"' | Measure-Object
```

<br>

📌 **실패한 테스트 코드 확인 명령어**
```
go test ./internal/... -json | Select-String -Pattern '"Action":"fail"' | Measure-Object
```

<br><br>

### Swagger 테스트 순서

... 작성 중 ⚙

<br><br>

### Swagger 테스트

... 작성 중 ⚙

<br><br>

### Application Server Architecture

<img src="readme/image/server-architecture.png" alt="Application Server Architecture" width="800"/>

📌 [참고 Link](https://github.com/bxcodec/go-clean-arch)

<br><br>

### 폴더 구조

... 작성 중 ⚙

<br><br>

### ERD(Entity Relationship Diagram)

... 작성 중 ⚙

<br><br>

### API Endpoint

| HTTP Method | URI                       | Description              |
|-------------|---------------------------|--------------------------|
| `작성 중...`   | `작성 중...`               | 작성 중...                  |

<br><br>

### Git 커밋 메시지 규칙

| Tag        | Description                                         |
|------------|-----------------------------------------------------|
| `feat`     | 새로운 기능을 추가한 경우 사용합니다.                               |
| `fix`      | 버그를 수정한 경우 사용합니다.                                   |
| `refactor` | 코드 리팩토링한 경우 사용합니다.                                  |
| `style`    | 코드 형식, 정렬, 주석 등의 변경(동작에 영향을 주는 코드 변경 없음)한 경우 사용합니다. |
| `test`     | 테스트 추가, 테스트 리팩토링(제품 코드 수정 없음, 테스트 코드에 관련된 모든 변경에 해당)한 경우 사용합니다.                                             |
| `docs`     | 문서를 수정(제품 코드 수정 없음)한 경우 사용합니다.                                             |
| `chore`    | 빌드 업무 수정, 패키지 매니저 설정 등 위에 해당되지 않는 모든 변경(제품 코드 수정 없음)일 경우 사용합니다.                                             |

<br><br>