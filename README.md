# itsp-kadai
ITSPの課題(1Q)

[![CircleCI](https://circleci.com/gh/akky2501/itsp-kadai.svg?style=svg)](https://circleci.com/gh/akky2501/itsp-kadai)


## 概要
簡単なスケジュール管理APIサーバー

### API一覧
```
# イベント登録 API request
POST /api/v1/event
{"deadline": "2019-06-11T14:00:00+09:00", "title": "レポート提出", "memo": ""}

# イベント登録 API response
200 OK
{"status": "success", "message": "registered", "id": 1}

400 Bad Request
{"status": "failure", "message": "invalid date format"}
```

```
# イベント全取得 API request
GET /api/v1/event

# イベント全取得 API response
200 OK
{"events": [
    {"id": 1, "deadline": "2019-06-11T14:00:00+09:00", "title": "レポート提出", "memo": ""},
    ...
]}
```

```
# イベント1件取得 API request
GET /api/v1/event/${id}

# イベント1件取得 API response
200 OK
{"id": 1, "deadline": "2019-06-11T14:00:00+09:00", "title": "レポート提出", "memo": ""}

404 Not Found
```



## 構成
- Golang
  - main.go (サーバー本体)
  - server_test.go (レスポンスのテスト)
- Gin
- GORM
- mysql
- Docker / Docker Compose
- CircleCI