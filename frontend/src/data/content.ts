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
export const posts = Array.from({ length: 12 }, (_, index) => ({
  slug: index === 0 ? "recommended-clubs" : `article-${index + 1}`,
  title:
    index === 0 ? "おすすめサークル５選" : "タイトルタイトルタイトルタイトル",
  date: "2000-01-01",
  displayDate: "2000/01/01",
}));
export const achievements = Array.from({ length: 5 }, (_, index) => ({
  id: index + 1,
  title: "サークルのホームページを開発",
}));
export const articleCopy = {
  introduction: "私も今けっしてそうした逡巡人に対する点の頃を考えたた。",
  body: "とうとう前に発展院もついにこの教育だべきほどがありて切らですへは意味したたて、一応には思いですたらだた。径路に描けるませ方はいったん元来がよしでありない。",
  closing:
    "ほとんど岡田さんを所有個人少し汚辱でしよう本立こんな国家これか教育にに従ってご想像ないたたたて、同じ前はあなたか人形がしから、",
};
