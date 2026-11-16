// Figmaの表示確認用データ。API接続時は、この層をバックエンドからの取得に置き換えます。
export const posts = Array.from({ length: 12 }, (_, index) => ({
  slug: index === 0 ? "recommended-clubs" : `article-${index + 1}`,
  title:
    index === 0 ? "おすすめサークル５選" : "タイトルタイトルタイトルタイトル",
  date: "2000-01-01",
  displayDate: "2000/01/01",
}));
