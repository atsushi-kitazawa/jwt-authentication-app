# JWT Authentication App

Go 製の JWT 認証 API と、その API を操作する TUI クライアントを同じリポジトリで管理しているサンプルです。

## 構成

- [server](server/README.md)
  `net/http` と SQLite を使った Web API
- [client](client/README.md)
  Charmbracelet を使った Go 製 TUI クライアント

API とクライアントはそれぞれ独立した Go モジュールです。

## 起動

### Server

```bash
cd server
go run .
```

### Client

```bash
cd client
go run .
```
