# SRUUI

[『SHOWROOM イベント 獲得ポイント一覧』(1) 関連ソースの公開について](https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/4fccae)

[『SHOWROOM イベント 獲得ポイント一覧』(5) ルーム情報の更新](https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/d9f919)

----------------------

これは[『SHOWROOM イベント 獲得ポイント一覧』(1) 関連ソースの公開について](https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/4fccae)から始まる一連の記事の一つです。

[Github - SRCGI](https://github.com/Chouette2100/SRCGI)（[『SHOWROOM イベント 獲得ポイント一覧』(3) 獲得ポイントの推移の表示とデータ取得対象イベント・配信者の設定](https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/56ec9b)）でイベントを登録したときイベントに参加しているルームの情報が保存されます。

しかし、ルーム情報つまりルーム名、ランク、レベル、フォロワー数などは刻々と変化します。

このモジュールは現在イベント参加中のルームのルーム情報を最新の状況に更新するためのものです。

このモジュールを定期的に起動することにより保存されている情報を最新の状態に保つことができます。

---

実行方法

まず次のようなシェルを用意します。

```
$ cd ~/MyProject/Showroom/UpdateUserInf
$ cat UpdateUserInf
#! /bin/sh
cd ~/MyProject/Showroom/UpdateUserInf
export DBNAME=xxxxxxxxx
export DBUSER=xxxxxxxxx
export DBPW=xxxxxxxx
./SRUUI  1>>tmplog.txt 2>&1
```

データベースのログインパスワードはServerConfig.ymlの方に直接書いてもいいのですが、ここでは環境変数を介して与える方法を使っています。これはServerConfig.ymlが公開するパッケージに含まれているからというのが理由です。

それから、SRUUIについは現状ではログファイル出力と標準出力があるので　「  1>>tmplog.txt 2>&1」　の部分は余計です。ログの出力はご自身の運用形態に合わせて調整していただければと思います。いずれにしても「 2>>tmplog.log」はあったほうがいいかも。

調整というのは SRUUI.go の
```
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))
```
のところです。


次にこれをcronで起動します。

```
$ crontab -e
.
.
10 3-23/12 * * * ~chouette/MyProject/Showroom/UpdateUserInf/UpdateUserInf.sh
. 
.
```

これは毎日03時10分と15時10分に起動する例です。

---

ロードモジュールの作成と設置

現在 Linux Mint 21.1 Vera base: Ubuntu 22.04 jammy 、 go version go1.20.4 linux/amd64　で作成したロードモジュールを VPS（Ubuntu 20.04.4 LTS focal）に持っていって動かしているのですが、この場合次のような手順になります。

まず[Github - SRUUI](https://github.com/Chouette2100/SRUUI)から入手したソースを~/go/src/SRUUI 以下におきます。

以下

```
$ cd ~/go/src/SRUUI
$ go mod init
$ go mod tidy
$ CGO_ENABLED=0 go build SRUUI.go
$ sftp -oServerAliveInterval=60 -i ~/.ssh/id_ed25519 -P nnnn xxxxxxxxnnn.nnn.nnn.nnn
sftp> cd ~/MyProject/Showroom/UpdateUserInf
sftp> put SRUUI
sftp> put ServerConfig.yml
```

みたいな感じで進めます。

なお CGO_ENABLED=0 は最近VPSにもっていったときライブラリーのエラー（/lib/x86_64-linux-gnu/libc.so.6: version `GLIBC_2.32' not found）が出るようになったので入れています、

VPSの方が
```
ldd (Ubuntu GLIBC 2.31-0ubuntu9.9) 2.31
```
でローカルの方が

```
ldd (Ubuntu GLIBC 2.35-0ubuntu3.1) 2.35
```
なので GLIBC_2.32 以上が必要ということでしょうか。

Goの場合ロードモジュールは必要なライブラリーをぜんぶ持ってるはず、なぜ？、ということでざっとググったところ、こういうエラーが起きるのはnet/http などのパッケージを使った場合内部的にはCライブラリを使う方法と純粋なGoですませる方法があり、何もしないと（CGO_ENABLED=0 を指定しないと）前者の方法になりCライブラリーがダイナミックリンクされてしまう、というのが原因のようです。このあたりの事情は正直よくわかってません。すみません。

なおVSCodeを使っている場合、tasks.jsonは次のような書き方でいいようです。ふつうと違うのはoptionsの部分です。

```
        {
            "label": "go build",
            "type": "shell",
            "options": {
                "env": {
                  "CGO_ENABLED": "0"
                }
            },
            "command": "go",
            "args": [
                "build",
                "-v",
                "./..."
            ],
            "problemMatcher": [],
            "group": {
                "kind": "build",
                "isDefault": true
            }
        },
```
