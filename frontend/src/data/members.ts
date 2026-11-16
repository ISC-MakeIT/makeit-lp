// Figmaの表示確認用データ。API接続時は、この層をバックエンドからの取得に置き換えます。
export const memberDescription =
  "親譲りの無鉄砲で小供の時から損ばかりしている。小学校に居る時分学校の二階から飛び降りて一週間ほど腰を抜かした事がある。なぜそんな無闇をしたと聞く人があるかも知る。";
export const members = Array.from({ length: 12 }, (_, index) => ({
  id: index + 1,
  name: "岩崎太郎",
  department: "情報処理科",
  description:
    index === 1
      ? "親譲りの無鉄砲で小供の時から損ばかりしている。"
      : memberDescription,
  github: "IWSK1225",
  website: "iwasaki-taro.com",
}));
