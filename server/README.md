# JWT Authentication App Server

Go の `net/http` と SQLite を使ったシンプルな Web API サンプルです。

## 機能

- ユーザ登録
- ログイン
- JWT 認証
- ユーザ一覧取得
- ユーザ参照
- ユーザ更新
- ユーザ論理削除
- `health` エンドポイント

## 技術スタック

- Go 1.25.5
- `net/http`
- SQLite
- JWT
- bcrypt

## セットアップ

```bash
go mod download
```

## 起動方法

```bash
go run .
```

`JWT_SECRET` を指定したい場合は環境変数で渡せます。

```bash
JWT_SECRET=my-secret go run .
```

サーバーは `http://localhost:8080` で起動します。

## 認証ルール

- 未認証でアクセス可能
  - `GET /docs`
  - `GET /openapi.yaml`
  - `GET /health`
  - `POST /users`
  - `POST /login`
- JWT 認証が必要
  - `GET /`
  - `GET /users`
  - `GET /users/name/{name}`
  - `GET /users/{id}`
  - `PUT /users/{id}`
  - `DELETE /users/{id}`

## API

### Swagger UI

ブラウザで次にアクセスすると Swagger UI を確認できます。

```text
http://localhost:8080/docs
```

### Health Check

```bash
curl http://localhost:8080/health
```

### ユーザ登録

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"taro","email":"taro@example.com","password":"secret123"}'
```

レスポンス例:

```json
{"id":"1","name":"taro","email":"taro@example.com"}
```

### ログイン

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"name":"taro","password":"secret123"}'
```

レスポンス例:

```json
{"token":"<JWT_TOKEN>"}
```

以降の保護された API では、`Authorization` ヘッダに JWT を付与します。

```bash
-H "Authorization: Bearer <JWT_TOKEN>"
```

### ユーザ一覧

```bash
curl http://localhost:8080/users \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### ユーザ参照

```bash
curl http://localhost:8080/users/1 \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### ユーザ名で参照

```bash
curl http://localhost:8080/users/name/taro \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### ユーザ更新

```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{"name":"taro","email":"taro-updated@example.com","password":"new-secret123"}'
```

### ユーザ削除

論理削除です。物理削除は行いません。

```bash
curl -X DELETE http://localhost:8080/users/1 \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

## データ保存

- SQLite ファイルは `app.db` に保存されます
- パスワードは bcrypt でハッシュ化して保存されます
- 削除済みユーザは `deleted_at` に時刻を保存して論理削除されます

## 備考

- 初回起動時に `users` テーブルを自動作成します
- 既存 DB に不足カラムがある場合は、起動時に追加されます

## TODO

- JWT 関連処理を `github.com/golang-jwt/jwt/v5` に依存せず、標準ライブラリを使った独自実装へ置き換える
- 置き換え対象には、トークン生成、署名検証、期限チェック、Bearer トークンの取り扱いを含める
