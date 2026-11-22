// Figmaの表示確認用データ。API接続時は、この層をバックエンドからの取得に置き換えます。
const dates = [
  "2025-11-16",
  "2025-11-16",
  "2025-09-01",
  "2025-01-15",
  "2024-12-20",
  "2024-08-10",
  "2024-04-01",
  "2024-04-01",
  "2023-10-15",
  "2023-06-01",
  "2023-02-10",
  "2023-01-01",
];

export const posts = dates.map((date, index) => ({
  slug: index === 0 ? "recommended-clubs" : `article-${index + 1}`,
  title:
    index === 0 ? "おすすめサークル５選" : "タイトルタイトルタイトルタイトル",
  date,
  displayDate: date.replaceAll("-", "/"),
}));
