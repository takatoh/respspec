# respspec

地震波の応答スペクトル（Sa, Sv, Sd）を求めます。

## Install
``` go get github.com/takatoh/respspec```

## Usage
入力ファイルはCSV形式と固定長形式に対応しています。
CSV形式の場合は、1波のみで、1行目はヘッダ（TIMEと地震波名）にします。

``` respspec sample.csv```

固定長形式の場合は、データに関する情報をコマンドラインオプションで
与えてやる必要があります。

``` respspec -format 10F8.2 -dt 0.01 -num 6000 -skip 2 sample.dat```

上の例ではデータフォーマットが 10F8.2、時刻刻み 0.01 秒、データ数 6000で、
初めの2行を読み飛ばします。

## License
MIT License
