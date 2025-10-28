あなたは優秀なエンジニアです。

## ゴール
ゲームボーイアドバンスのゲームを作る事


## 制約
- docs
- GameBoyアドバンスの仕様については、Webで検索し、結果を実装docs配下におくこと
- 便利関数は共通化したところに配置する事
- ゲームは複数作成するため、各ディレクトリ(=Game)でbuild可能にする事。go.workを使う事で、便利関数を利用できるようにする
- Go言語だけで実装する

- testも実装しながら、キリの良いところで、git add/commit/pushすること


## ディレクトリ構成

- .
 |- docs/ ... gameboyの実装の方針や共通処理の設計など
 | bin/ xx.gb のファイルが存在する
 |- {game名}
    |- docs ... gameの仕様
    |- main.go ..
    |- go.mod
 |- {game}
    |- docs ... gameの仕様
    |- main.go ..build対象
    |- go.mod

