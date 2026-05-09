# Client Web App

Nuxt で作ったブラウザ向けクライアントです。`Login` で取得した JWT を保持しながら、各APIを実行してレスポンスを確認できます。

## 起動

```bash
cd server
go run .
```

別ターミナルで:

```bash
cd client
npm install
npm run dev
```

Nuxt アプリは通常 [http://localhost:3000](http://localhost:3000) で起動します。

## バックエンド接続先

デフォルトでは `http://127.0.0.1:8080` を参照します。必要に応じて Nuxt 起動時に変更できます。

```bash
BACKEND_BASE_URL=http://127.0.0.1:8080 npm run dev
```

画面上部の `Backend Base URL` を変更すると、その値を優先して Nuxt のプロキシ経由で接続します。ブラウザから直接バックエンドへはアクセスせず、Nuxt サーバーが中継するので CORS 設定は不要です。

## 操作

- 左の `API Menu` から実行したいAPIを選択
- `Backend Base URL` を必要に応じて変更
- `Create User` または `Login` で初期セットアップ
- 認証が必要なAPIでは `Login` 後のJWTが自動送信されます
- `Clear Token` で保存中のJWTを破棄できます

## 想定フロー

1. `Create User` でユーザーを作成
2. `Login` でJWTを取得
3. `Protected Root` または `List Users` で認証付き通信を確認
4. `Get User` / `Get User By Name` / `Update User` / `Delete User` を必要に応じて実行
