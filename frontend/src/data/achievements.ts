// Figmaの表示確認用データ。API接続時は、この層をバックエンドからの取得に置き換えます。
export const achievements = Array.from({ length: 5 }, (_, index) => ({
  id: index + 1,
  href:
    index === 0 ? "/blog/recommended-clubs/" : `/blog/article-${index + 1}/`,
  title: "サークルのホームページを開発",
}));
