# 目標
windowsでGO言語、node.jsの環境構築をして、smarket_draftを動かすこと。

## 必要な環境
golang
node.js
npm

### 筆者のPC環境
windows 10 intel core i7 (64bit)

# 1.Go言語のための環境構築
下記の３つの手順に従って環境構築をします。
なお、[はじめてのGo言語(on windows)](https://qiita.com/spiegel-im-spiegel/items/dca0df389df1470bdbfa)を参考しました。

## 手順①Go言語のコンパイラをインストールする。
windowsのコマンドプロンプトでGo言語を使用するために、[ダウンロードページ](https://golang.org/dl/)でwindows版のコンパイラをダウンロードします。

現在(2018/09/01)の最新versionは「go1.11.windows-amd64.msi」。

![スクリーンショット (236)_LI.jpg](https://qiita-image-store.s3.amazonaws.com/0/251551/98547195-c938-1597-7110-2b0dbed468b7.jpeg)

うまくインストールされれば、C:直下に「GO」フォルダが作成されます。

## 手順②環境変数を設定する。

「環境変数」のGOROOTを「C:\GO\」に設定し、pathを「C:\GO\bin\」に設定します。

（手順①のインストールで.msiパッケージをインストールした場合は、自動的に環境変数を設定してくれているので、正しく設定されているか確認するだけでよい。）

## 手順③動作確認

windowsのコマンドプロンプトでversionを確認する。

~~~
C:go version
go version go1.11 windows/amd64
~~~

インストールされていることが確認できる。
以上でコンパイラのインストールは完了です。

#smarket_draftを動かすための前準備(GO言語)

まずは、sorcetreeからsmarket_draftフォルダをダウンロードしておきます。

## 必要なパッケージをインポートする

このプロジェクトでimportするパッケージは以下のパッケージです。

「github.com/golang/protobuf/proto」
「google.golang.org/grpc」
「golang.org/x/net/context」
「github.com/mtfelian/golang-socketio」
「github.com/sirupsen/logrus」
「github.com/bwmarrin/snowflake」

gitを開き「C:GO/src」ディレクトリに移動し、以下のコマンドでそれぞれインポートしてください。

~~~
$ go get github.com/golang/protobuf/proto
$ go get google.golang.org/grpc
$ go get golang.org/x/net/context
$ go get github.com/mtfelian/golang-socketio
$ go get github.com/bwmarrin/snowflake
$ go get github.com/sirupsen/logrus
~~~

これらのパッケージをインストールすることでプロジェクト内で使用することができるようになります。

## エラーと対処法
Goファイルを実行した際、以下のようなエラーが発生することがあります。

~~~
$ go run fleet-provider.go             
fleet-provider.go:11:2: cannot find package "github.com/gorilla/websocket" in any of:                                
   C:\Go\src\github.com\gorilla\websocket (from $GOROOT)                 
   C:\Users\Rui Hirano\go\src\github.com\gorilla\websocket (from $GOPATH)
~~~

github.com/gorilla/websocketがインストールされていないというエラーですので
その都度、C:Go/srcディレクトリ内で下記コマンドでパッケージをインストールしてください。

~~~
go get github.com/gorilla/websocket
~~~

# 2.node.jsのための環境構築
続いてjavascript言語のコンパイラ(node.js,npm)をインストールします。

windowsの場合、nodistというnode.jsのバージョン管理を行えるツールをインストールすることになります。

## 手順①nodistのインストール

[ダウンロードページ](https://github.com/marcelklehr/nodist/releases)にてnodistをインストールします。
現在(2018/09/01)の最新版はv0.88です。

環境変数はインストーラが自動で行ってくれるそうです。

## 手順②動作確認
コマンドプロンプトで以下のコマンドを実行。

~~~
C: nodist -v
0.8.8

C: npm -v
4.0.5
~~~

続いて最新版をインストールします。

執筆時点では10.9.0が最新でした。

~~~
C: nodist dist  //インストール可能なバージョンを表示

C: nodist + 10.9.0  //表示された中から最新バージョンをインストール

C: nodist 10.9.0   //使用するバージョンを指定
~~~

npmを最新版にします。
最低でも5.6.0以上にしましょう。
執筆時点では6.1.0が最新版でした。

~~~
nodist npm 6.1.0   //npmの最新バージョンをインストール
~~~

再度、バージョンを表示して最新版になっているか確認してみてください。
インストールの動作確認は以上です。

# smarket-draftを動かすための前準備(node.js,npm)


## node_modulesにパッケージをインストールする。
clientディレクトリにて下記コマンドを実行

~~~
npm outdated
~~~

必要なパッケージの最新版が表示されます。
![スクリーンショット (241).png](https://qiita-image-store.s3.amazonaws.com/0/251551/3de65e55-98a0-f7cf-9cee-e75f9ccd1883.png)

react-scriptsのみ最新版があるので指定してインストールします。

~~~
npm install react-scripts@1.1.5
~~~
node_modulesフォルダにreacts-scriptsフォルダができインストールが成功したのが確認できると思います。

他のパッケージは下記コマンドでまとめてインストールします。

~~~
npm install
~~~
node_modulesにそれぞれフォルダができているのを確認してみてください。

その後、以下のコマンドを実行してパッケージをビルドすれば終了です。

~~~
npm run build
~~~

## エラーとその対処法

npm run buildにて以下のエラーが表示された場合、パッケージ(下記の場合、react-dom)がインストールされていない可能性があります。

~~~
> npm run build

Creating an optimized production build...
Failed to compile.

Module not found: Error: Can't resolve 'react-dom' in 'C:\Users\Rui Hirano\uclab_nu_smarket_draft\src\monitor\client\src'
~~~

clientディレクトリにて下記コマンドを実行してください。

~~~
npm install
~~~

# smarket-draftの実行、起動手順

動作確認には複数のコマンドプロンプトが必要になります。
windowsの場合、conEmuという複数のコマンドプロンプトを同時に表示できるソフトがあるのでそれを用いることをお勧めします。

![スクリーンショット (242).png](https://qiita-image-store.s3.amazonaws.com/0/251551/d53b4bc2-8d01-acbc-da45-5f9fd272e707.png)

まずserverを三つ起動します。

起動する順番は
①nodeid-sever
②moniter-server
③smarket-server
という順番です。

起動すると [http://127.0.0.1:9999](http://127.0.0.1:9999) にてモニターを見ることができます。

![スクリーンショット (243).png](https://qiita-image-store.s3.amazonaws.com/0/251551/922f81d9-be17-325b-0835-7e17f803408e.png)


そのあとはシナリオ次第でpub/subできます。

[ad]

~~~
go run ad-provider.go
~~~

[taxi]

~~~
go run taxi-provider.go -price 100
~~~

[user]

~~~
go run user-provider.go
~~~

[fleet]
自動車情報をリアルタイムで受信します。

Read.meに記載されたリンクから自動車情報の可視化状況を見ることができます。

~~~
go run fleet-provider.go
~~~

![スクリーンショット (244).png](https://qiita-image-store.s3.amazonaws.com/0/251551/1915b8ea-3760-1194-3a44-ee1c9738c763.png)


以上です。
そのほかの詳細はbacklogに記載されているのでそちらを参考してください。
ありがとうございました。

## 参考文献
[はじめてのGo言語(on windows)](https://qiita.com/spiegel-im-spiegel/items/dca0df389df1470bdbfa)
[Github上のパッケージを参照する](https://maku77.github.io/hugo/go/github.html)
[nodistでnode.jsをバージョン管理](https://qiita.com/satoyan419/items/56e0b5f35912b9374305)
