# image-server

AWS をインフラとした画像サーバー

## 説明

### GET メソッド

画像を取得します

#### クエリパラメータ

| key      | type   | 概要                       |
| :------- | :----- | :------------------------- |
| type     | string | 変換設定名                 |
| w        | \*int  | 幅                         |
| h        | \*int  | 高さ                       |
| crop     | string | 変換方式 cover / contain   |
| quality  | \*int  | クオリティ                 |
| webpq    | \*int  | webp のクオリティ          |
| avifq    | \*int  | avif のクオリティ          |
| jpegq    | \*int  | jpeg のクオリティ          |
| lossless | \*bool | ロスレスにする             |
| upscale  | \*bool | サイズアップを行うかどうか |

TODO: Etag や last-modified

### Put メソッド

画像をアップロードします

#### リクエストボディ

multipartform 想定

| key   | type     | 概要         |
| :---- | :------- | :----------- |
| image | ファイル | 画像ファイル |

## 環境変数

| key                                  | type   | 概要                                                         |
| :----------------------------------- | :----- | :----------------------------------------------------------- |
| IMAGE_SERVER_DISABLE_EXTERNAL_PARAMS | bool   | URL パラメータを外部パラメータを許可するかどうか             |
| IMAGE_SERVER_CONVERT_PARAMS          | JSON   | 変換設定                                                     |
| IMAGE_SERVER_BUCKET                  | string | S3 バケット名                                                |
| IMAGE_SERVER_S3_URL                  | string | S3 の URL AWS ででは指定しない                               |
| AWS_ACCESS_KEY_ID                    | string | 基本的には使わず、Fargate の TaskRole が使えない場合のみ指定 |
| AWS_SECRET_ACCESS_KEY                | string | 基本的には使わず、Fargate の TaskRole が使えない場合のみ指定 |
| AWS_REGION                           | string | 基本的には使わず、Fargate の TaskRole が使えない場合のみ指定 |

### IMAGE_SERVER_CONVERT_PARAMS

画像の変換のパラメータを事前に指定したいときにするときに使用する。
「GET リクエストで使えるパラメータ+画像」を JSON の配列形式で複数指定することができる。
GET リクエストの type によって指定した変換設定を呼び出すことができる。

```
[{"type":"mini","w":50,"h":50,"lossless":true}]
```

## AWS へのデプロイ

TODO: まだバグがある
