# MakeIT フロントエンド

Figma `makeit-lp` を既存のAstro 5に実装。React・Tailwindは使用せず、AstroコンポーネントとCSSで構成しています。contactページとフォーム送信は対象外です。公開・デプロイは行いません。

## 起動・確認

リポジトリのルートで `docker compose up -d`。フロントだけを動かす場合はこのディレクトリで以下を実行します。

```sh
npm ci
npm run dev
# http://localhost:4321
npm run check
npm run build
```

バックエンドやMySQLを起動しなくても画面は確認できます。

## ページ

| URL                        | 内容                                          | Figma参照ノード                 |
| -------------------------- | --------------------------------------------- | ------------------------------- |
| `/`                        | トップ                                        | `443:270` / モバイル `635:1178` |
| `/member/`                 | メンバー一覧                                  | `634:749` / モバイル `1414:764` |
| `/blog/`                   | ブログ・ニュース一覧                          | `668:967`, `668:968`, `668:970` |
| `/blog/recommended-clubs/` | 記事詳細のサンプル                            | `669:1256`                      |
| `/invitation/`             | 活動・入部案内（Figmaページ名 `invitatioin`） | `3062:1181`                     |
| `/recruit/`                | 募集案内・実績                                | `3043:1300`                     |

[元のFigmaデザイン](https://www.figma.com/design/NWjPGQAnAivTqpD2ezuaxs/makeit-lp?node-id=443-270)

FigmaのMember画面には `contact-page` というフレーム名が付いていますが、内容はメンバー一覧なので実装対象です。`contact` ページ（`69:17`）は実装していません。

`Join/Job` は現段階では `/invitation/#activity` に接続しています。入部・依頼フォームのURLが確定したら `Header`、`Hero`、`ContactBanner` のリンクを差し替えてください。Figma内の `Mender` は `Member` に統一しています。

## ディレクトリ構成

```text
frontend/
├── public/
│   └── assets/figma/  # Figmaから取得した画像・SVG。期限付きURLへの依存なし
├── src/
│   ├── components/    # ヘッダー、フッター、セクション、カルーセルなど
│   ├── data/          # メンバー、記事、実績の仮データ
│   ├── layouts/       # HTML、メタ情報、共通レイアウト
│   ├── pages/         # Astroのファイルベースルーティング
│   └── styles/        # 色・寸法・レスポンシブスタイル
└── astro.config.mjs
```

## フォント

Outfit・Inter・Noto Sans JPは、`src/layouts/BaseLayout.astro` の共通`head`からGoogle FontsのCSS APIで読み込みます。`fonts.googleapis.com` と `fonts.gstatic.com` への `preconnect`、読み込み中に代替フォントを表示する `display=swap` を指定しています。

書体の使い分けは `src/styles/global.css` にあります。ローカルのフォントファイルと `@font-face` 定義は不要です。フォント取得にはインターネット接続が必要で、取得できない場合はCSSで指定した代替フォントが使われます。

## 後からバックエンドを接続するとき

- `../backend/` はGo、DBはMySQL 8.4。既存の `compose.yaml`、ボリューム、サービス間依存は維持しています。
- 現在のAPIは `/api/health` のみです。未実装APIへのリクエストは画面表示時に発行しません。
- `src/data/content.ts` にFigmaの仮データを集約しています。氏名・本文・日付・画像の灰色枠は実データではありません。個人プロフィールの仮URLは外部へ遷移しない表示のみです。
- 記事は `src/pages/blog/[slug].astro` の `getStaticPaths()` で静的生成します。ビルド時取得ならこのデータ層をAPI取得に置き換えます。随時更新が必要ならAstroのSSRアダプター、またはクライアント側での取得を追加します。
- ブラウザから呼ぶAPIのベースURLには、既存の `PUBLIC_API_URL` を使えます。Docker内のビルド時・SSRから取得する場合は `backend:8080` を使うサーバー専用環境変数を別途追加してください。

カルーセルはドット・横スクロール・左右キーで操作できます。記事一覧は新旧順と年の絞り込み、モバイルナビは開閉とEscapeキーに対応しています。

## ブランドアイコン

Discord・GitHubのアイコンは[Bootstrap Icons](https://icons.getbootstrap.com/)（MITライセンス）を使用しています。
`BrandIcon.astro`で必要なSVGだけをnpmパッケージから読み込み、HTMLへ埋め込んでいます。
ライセンスは`public/licenses/bootstrap-icons.txt`に同梱しています。
