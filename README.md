# commerce-system

안녕하세요, 홍정완입니다.

본 Repository는 1주일 간의 개발 성과를 평가하기 위해 작성되었습니다.

평가가 원활히 진행될 수 있도록, 아래 README.md에 작성된 내용을 충분히 검토해 주시면 감사하겠습니다.

---

### 실행 방법 --- 작성 중 ⚙ ---
1. `docker-compose up` 명령어 실행
2. `docker build --no-cache -t my-app-image .` 명령어 실행

<br>

```
# 📌 GoLand IDE에서 프로그램 실행 시 설정하는 방법

- Program arguments: Go 프로그램을 실행할 때 전달할 명령줄 인수. 
- 여기서는 `-c deploy/config.toml`을 전달하여 `config.toml` 파일을 설정 파일로 사용합니다.

- 설정 방법
  1. GoLand에서 Run/Debug Configurations를 엽니다.
  2. Program arguments 필드에 `-c deploy/config.toml`을 입력합니다.
  3. 이 설정은 프로그램이 `config.toml` 파일을 읽어들이도록 하여, 지정된 환경 설정을 로드하게 합니다.
  4. 설정을 저장하고, Run 버튼을 클릭하여 프로그램을 실행합니다.
```

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