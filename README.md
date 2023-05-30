# SRUUI

[SRCGI](https://github.com/Chouette2100/SRCGI)でイベントを登録したときにルーム情報が保存されます。

ルーム情報つまりルーム名、ランク、レベル、フォロワー数などは刻々と変化します。

このモジュールは現在イベント参加中のルームのルーム情報を最新の状況に更新するためのものです。

このモジュールを定期的に起動することにより保存された情報を最新の状態に保つことができます。

具体的には

まず次のようなシェルを用意します。

```
$ cd ~/MyProject/Showroom/UpdateUserInf
$ cat UpdateUserInf
#! /bin/sh
cd ~/MyProject/Showroom/UpdateUserInf
export DBUSER=xxxxxxxxx
export DBPW=xxxxxxxx
./SRUUI  1>>tmplog.txt 2>&1
```

※ SRUUIについてはログ出力がるので　「  1>>tmplog.txt 2>&1」　の部分は余計ですが、参考として...

※ 「 2>>tmplog.log」くらいはあったほうがいいかも。

次にこれをcronで起動します。

```
$ crontab -e
.
.

15-59/30 * * * * ~chouette/MyProject/Showroom/UpdateUserInf/UpdateUserInf.sh
.
.
```

これは毎時15分と45分に起動する例です。まあじっさいは1日に一回も更新すればじゅうぶんな気がしますが。


ロードモジュールの作成と設置

現在 Linux Mint 21.1 Vera base: Ubuntu 22.04 jammy 、 go version go1.20.4 linux/amd64　で作成したロードモジュールを Ubuntu 20.04.4 LTS ocal のVPSに持っていって動かしているのですが、この場合次のような手順になります。

まずGithubで入手したソースを~/go/src/SRUUI 以下におきます。
以下
```
$ cd ~/go/src/SRUUI
$ go mod init
$ go mod tidy
$ CGO_ENABLED=0 go build SRUUI.go
$ sftp -oServerAliveInterval=60 -i ~/.ssh/id_ed25519 -P nnnn xxxxxxxxnnn.nnn.nnn.nnn
$ cd ~/MyProject/Showroom/UpdateUserInf
$ put SRUUI
```

みたいな感じで進めます。

なお CGO_ENABLED=0 は最近VPSにもっていったときライブラリーのエラー（/lib/x86_64-linux-gnu/libc.so.6: version `GLIBC_2.32' not found）が出るようになったので入れています、このあたりの事情は正直よくわかってません。すみません。



