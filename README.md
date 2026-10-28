# makeit-lp

| 役割 | 技術 | URL / ポート |
|------|------|--------------|
| フロントエンド | Astro (Node 22) | http://localhost:4321 |
| バックエンド | Go 1.24 (air でホットリロード) | http://localhost:8080 |
| データベース | MySQL 8.4 | localhost:3306 |

**Docker上で全て動くので、Go・Node・MySQLを自分のPCにインストールする必要はありません。**
手順は **Windows(WSL2) / Mac / Linux** で共通です。違うのは「事前準備」だけです。

---

## 1. 事前準備

### Windows の方(WSL2 を使います)

> ⚠️ **PowerShell や コマンドプロンプトでは作業しないでください。** 以降の操作はすべて WSL(Ubuntu)の中で行います。

1. **WSL2 + Ubuntu を入れる**
   PowerShell を「管理者として実行」して次を実行し、PCを再起動します。
   ```powershell
   wsl --install -d Ubuntu
   ```
   再起動後に Ubuntu が開くので、ユーザー名とパスワードを設定します。
2. **Docker Desktop をインストール**
   https://www.docker.com/products/docker-desktop/ から入れて起動します。
3. **Docker Desktop を WSL と連携させる**
   Docker Desktop → ⚙ Settings →
   - General: 「Use the WSL 2 based engine」にチェック
   - Resources → WSL Integration: 「Ubuntu」をオンにして **Apply & restart**
4. **Ubuntu(WSL)のターミナルを開いて確認**
   ```bash
   docker --version
   docker compose version
   ```
   どちらもバージョンが表示されればOKです。
5. **Git を確認**(Ubuntuには通常入っています)
   ```bash
   git --version
   ```

### Mac / Linux の方

- Docker Desktop(Mac)または Docker Engine + Compose プラグイン(Linux)と Git を入れてください。

---

## 2. プロジェクトを取得する

> 🚨 **Windows の方へ(最重要)**
> プロジェクトは必ず **WSL の中のフォルダ(`~/` 以下)** に置いてください。
> `/mnt/c/...`(Windows側のCドライブ)に置くと、動作が極端に遅くなったり、ファイル変更が検知されなくなったりします。

```bash
cd ~
git clone <このリポジトリのURL> makeit-lp
cd makeit-lp
```

---

## 3. 起動する

```bash
docker compose up -d
```

これだけです。初回は イメージのダウンロードとビルドで **数分** かかります。
(`.env` の作成や `npm install` などの手作業は不要です。)

### 起動確認

```bash
docker compose ps
```

`db` が `healthy`、`backend` と `frontend` が `running`(`Up`)なら成功です。

ブラウザで次を開いてください。

| 確認先 | 期待する表示 |
|--------|--------------|
| http://localhost:4321 | 「Astro is running」と、`API: {"api":"ok","db":"ok"}` |
| http://localhost:8080/api/health | `{"api":"ok","db":"ok"}` |

> フロントエンドは初回、`node_modules` のインストールが終わるまで 1〜2分 つながらないことがあります。
> `docker compose logs -f frontend` に `Local http://localhost:4321/` と出れば準備完了です。

---

## 4. 日常の使い方

同じ操作は `make` でも実行できます。一覧は `make` または `make help`。

| やりたいこと | コマンド |
|--------------|----------|
| 起動 | `docker compose up -d` |
| 停止 | `docker compose stop` |
| 停止して削除(DBデータは残る) | `docker compose down` |
| ログを見る | `docker compose logs -f` (サービス指定: `... logs -f backend`) |
| DBに入る | `docker compose exec db sh -c 'mysql -u$MYSQL_USER -p$MYSQL_PASSWORD $MYSQL_DATABASE'` |
| **DBデータも含めて完全初期化** | `docker compose down -v` |
| Dockerfileや依存を変えた後 | `docker compose up -d --build` |

コードを編集して保存すると、**自動で反映**されます(Go は再ビルド、Astro はブラウザ更新)。
編集はWindowsの方も WSL 内のフォルダで行ってください。VS Code なら Ubuntu のターミナルで `code .` を実行すると WSL 接続で開けます(拡張機能「WSL」が必要)。

### 依存パッケージを追加するとき

```bash
# フロント(例: tailwind)
docker compose exec frontend npm install <パッケージ名>
# バックエンド
docker compose exec backend go get <モジュール名>
```
追加後に変更された `package.json` / `package-lock.json` / `go.mod` / `go.sum` もコミットしてください。

---

## 5. 環境差分が出ないための設計

- Go / Node / MySQL のバージョンは `Dockerfile` と `compose.yaml` で**固定**(各自のPCの環境に依存しません)
- `package-lock.json` / `go.sum` をコミットし、依存バージョンを固定(フロントは `npm ci`)
- `.gitattributes` で改行を **LF に固定**(Windowsでの改行コード事故を防止)
- ファイル変更検知はポーリング方式(Windows/WSLの bind mount でも確実に動く)
- DB接続情報は環境変数で変更できます(下記「環境変数」参照)。何も設定しなくてもデフォルト値で動きます

## 環境変数(DBのユーザー名・パスワードなど)

| 変数 | デフォルト | 説明 |
|------|-----------|------|
| `DB_NAME` | `app` | データベース名 |
| `DB_USER` | `app` | アプリ用ユーザー名 |
| `DB_PASSWORD` | `app` | アプリ用パスワード |
| `DB_ROOT_PASSWORD` | `root` | MySQL root パスワード |

変更したい場合:
```bash
cp .env.example .env   # .env を編集して値を変更
docker compose down -v && docker compose up -d
```
> ⚠️ DBユーザー・パスワードは **DBの初回作成時にだけ** 反映されます。変更後は上のように `down -v` でDBを作り直してください(**DBの中身は消えます**)。
> `.env` は Git に含まれません(秘密情報をコミットしないでください)。デフォルト値は開発用です。

## 6. ディレクトリ構成

```
.
├── compose.yaml        # 全サービスの定義
├── Makefile            # 起動・停止などのショートカット (make help)
├── .env.example        # 環境変数のサンプル(任意)
├── backend/            # Go (main.go, Dockerfile, .air.toml)
└── frontend/           # Astro (src/pages/index.astro, Dockerfile)
```

---

## 7. トラブルシューティング

| 症状 | 原因と対処 |
|------|-----------|
| `docker: command not found` (WSL) | Docker Desktop が起動していない、または WSL Integration がオフ。手順1-3を確認 |
| `Cannot connect to the Docker daemon` | Docker Desktop を起動して、1分ほど待ってから再実行 |
| `port is already allocated` / `address already in use` | 3306・4321・8080 を別のソフトが使用中。そのソフトを止める(ローカルMySQLが多い)か、`compose.yaml` の `ports` の左側の数字を変更 |
| http://localhost:4321 がつながらない | 初回は `npm ci` 中。`docker compose logs -f frontend` で完了を待つ |
| `api` は ok だが `db` がエラー | DBの起動待ち。数秒後に再読み込み。直らなければ `docker compose logs db` |
| 保存しても反映されない | プロジェクトが `/mnt/c/...` にある可能性。`~/` 以下に置き直す |
| `exec ... no such file or directory` / `\r` が出る | 改行コードがCRLFになっている。`git config --global core.autocrlf input` を実行し、クローンし直す |
| 何かおかしい・最初からやり直したい | `docker compose down -v && docker compose up -d --build` (**DBの中身も消えます**) |
