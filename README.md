# JWT Authentication App

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
  - `GET /health`
  - `POST /users`
  - `POST /login`
- JWT 認証が必要
  - `GET /`
  - `GET /users`
  - `GET /users/{id}`
  - `PUT /users/{id}`
  - `DELETE /users/{id}`

## API

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
