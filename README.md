# respspec

地震動の応答スペクトル（Sa, Sv, Sd）を求めます。
また、スペクトル強度（SI）を求めることもできます。

## Install

    > go get github.com/takatoh/respspec

## Usage

入力ファイルはCSV形式と固定長形式に対応しています。

CSV形式の場合は、1波のみで、1行目はヘッダ（`TIME`と地震動名）にします。

    > respspec sample.csv

固定長形式の場合は、データに関する情報をコマンドラインオプションで
与えてやる必要があります。

    > respspec -format 10F8.2 -dt 0.01 -ndata 6000 -skip 2 sample.dat

上の例ではデータフォーマットが 10F8.2、時刻刻み 0.01 秒、データ数 6000で、
初めの2行を読み飛ばします。

スペクトル強度（SI）を求めるには、`-si` オプションを指定します。

    > respspec -si -h 0.2 sample.csv

スペクトル強度を求める場合、減衰定数 h を 0.2 とすることが多いようですので、`-h` オプションで指定します。

## License

MIT License
