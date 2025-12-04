# Kinu

Kinu is the server that resizing images. Inspired by [cookpad's](https://cookpad.com/) image server [Tofu](http://www.slideshare.net/mirakui/ss-8150494).

## Status

It's development now for first alpha release.

## Dependency

- ImageMagick >= 6.8.9-9 && < 7.0.0

## Installation

```
# Please to Install the ImageMagick before this
go get github.com/tokubai/kinu
```

or 

```
docker pull tokubai/kinu:1.0.0.alpha13
```

## Startup

#### go get

```
$ export KINU_LOG_LEVEL=info
$ export KINU_RESIZE_ENGINE=ImageMagick
$ export KINU_STORAGE_TYPE=File 
$ export KINU_FILE_DIRECTORY=/tmp/kinu
$ kinu
INFO[0000] setup resize engine                           resize_engine_type=ImageMagick
INFO[0000] setup storage                                 storage_type=File
INFO[0000] [starting kinu on 127.0.0.1:8080]
```

```
$ curl http://localhost:8080/version
```

#### docker run

```
$ docker run tokubai/kinu:1.0.0.alpha13
INFO[0000] setup resize engine                           resize_engine_type=ImageMagick
INFO[0000] setup storage                                 storage_type=File
INFO[0000] [starting kinu on 0.0.0.0:80]
```

```
$ curl http://localhost/version
```

## Usage

### ImageUpload

#### kinu.gem(Recommended)

```ruby
require 'kinu'
Kinu::Resource.new(:foods, 1).upload(open('/path/to/image.jpg'))
```

#### curl

```shell
$ curl -X POST -F id=1 -F name=foods -F image=@/path/to/image http://localhost/upload
```

### GetResizedImage

#### kinu.gem(Recommended)

```ruby
require 'kinu'
require 'open-uri'
uri = Kinu::Resource.new(:foods, 1).uri(width: 280, height: 300) #<URI::HTTP http://localhost/images/foods/w=280,h=300/1.jpg>
open(uri).read
```

#### curl

```shell
$ curl http://localhost/images/foods/w=280,h=300/1.jpg
```

### UploadToSandbox(Temporary)

#### kinu.gem(Recommended)

```ruby
require 'kinu'
Kinu::Sandbox.upload(open('/path/to/image.jpg')) #<Kinu::Sandbox:0x007fdf92bdc490 @id="db4f1509-e2f5-40a7-9944-a6b0024f2a24", @name="__sandbox__">
```

#### curl

```shell
$ curl -X POST -F image=@/path/to/image http://localhost/sandbox
```

### Attach from Sandbox

#### kinu.gem(Recommended)

```ruby
require 'kinu'
Kinu::Sandbox.upload(open('/path/to/image.jpg')).attach_to(:foods, 1) #<Kinu::Resource:0x007fdf92b44a00 @id="1", @name="foods">
```

#### curl

```shell
$ curl -X POST -F sandbox_id=db4f1509-e2f5-40a7-9944-a6b0024f2a24 -F name=foods -F id=1 http://localhost/sandbox
```

## Specification

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/images/:type/:geometry/:filename` | 画像の取得（リサイズ・トリミング対応） |
| POST | `/upload` | 画像のアップロード |
| POST | `/sandbox` | サンドボックスへの一時アップロード |
| POST | `/sandbox/attach` | サンドボックスから実画像への確定 |
| GET | `/version` | バージョン情報の取得 |
| GET | `/worker/stats` | ワーカー統計情報の取得 |

### Image resize query

画像取得エンドポイント `GET /images/:type/:geometry/:filename` のパラメータ詳細です。

#### Path Parameters

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `type` | string | 画像カテゴリー（リソース分類名） | `foods`, `products` |
| `geometry` | string | リサイズ・トリミング条件（後述） | `w=280,h=300` |
| `filename` | string | ファイル名＋拡張子 | `123.jpg`, `456.png` |

**Supported extensions**: `jpg`, `jpeg`, `png`, `webp`, `gif`, `data` (Data URI format)

#### Geometry Parameters

`geometry` は `key=value` をカンマで連結した形式です。**指定順序が固定**されています。

```
w → h → q → c → mc → wo → ho → cw → ch → aw → o → m
```

| Key | Type | Description | Default | Example |
|-----|------|-------------|---------|---------|
| `w` | int | リサイズ後の幅（px） | - | `w=280` |
| `h` | int | リサイズ後の高さ（px） | - | `h=300` |
| `q` | int | JPEG品質（0-100） | 80 | `q=85` |
| `c` | bool | 自動トリミング（Auto Crop） | false | `c=true` |
| `mc` | bool | 手動トリミング（Manual Crop） | false | `mc=true` |
| `wo` | int | 幅方向オフセット（mc時） | 0 | `wo=10` |
| `ho` | int | 高さ方向オフセット（mc時） | 0 | `ho=20` |
| `cw` | int | トリミング領域の幅（mc時必須） | - | `cw=200` |
| `ch` | int | トリミング領域の高さ（mc時必須） | - | `ch=150` |
| `aw` | int | 想定画像幅（mc時必須） | - | `aw=300` |
| `o` | bool | オリジナル画像を返す | false | `o=true` |
| `m` | string | 中間画像サイズ | - | `m=1000` |

**`m` valid values**: `"true"`, `"original"`, `"1000"`, `"2000"`, `"3000"`

#### Required Conditions

- `w`, `h`, `o=true`, `m` のいずれかが必要
- `mc=true` の場合は `cw`, `ch`, `aw` が全て必要

#### Examples

```shell
# 基本的なリサイズ（幅280、高さ300）
$ curl http://localhost/images/foods/w=280,h=300/1.jpg

# 品質指定付きリサイズ
$ curl http://localhost/images/foods/w=280,h=300,q=85/1.jpg

# 自動トリミング＋リサイズ
$ curl http://localhost/images/foods/w=280,h=300,c=true/1.jpg

# 手動トリミング
$ curl http://localhost/images/foods/w=200,h=150,mc=true,wo=10,ho=20,cw=200,ch=150,aw=300/1.jpg

# オリジナル画像を取得
$ curl http://localhost/images/foods/o=true/1.jpg

# 中間画像サイズを指定して取得
$ curl http://localhost/images/foods/m=1000/1.jpg

# Data URI形式で取得（拡張子を .data に）
$ curl http://localhost/images/foods/w=280,h=300/1.data
```

### Environment variables

| name                           | required | default value               | valid value type                                                                      | note                                                                               |
| ------------------------------ | -------- | --------------------------- | ------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| KINU_BIND                      | ☓        | 127.0.0.1:80                | IP:PORT / unix domain socket path / Einhorn(einhorn@[num]) / FileDescripter(fd@[num]) | Compliance with the specifications of the goji/bind package.                       |
| KINU_DEBUG                     | ☓        | none                        | true                                                                                  | enable pprof                                                                       |
| KINU_RESIZE_ENGINE             | ◯        | none                        | ImageMagick                                                                           |                                                                                    |
| KINU_LOG_LEVEL                 | ◯        | none                        | panic / fatal / error / warning / info / debug                                        |                                                                                    |
| KINU_LOG_FORMAT                | ☓        | text                        | ltsv / json / text                                                                    |                                                                                    |
| KINU_RESIZE_WORKER_MODE        | ☓        | none                        | true                                                                                  |                                                                                    |
| KINU_RESIZE_WORKER_MAX_SIZE    | ☓        | cpu num * 10                | Integer                                                                               |                                                                                    |
| KINU_RESIZE_WORKER_WAIT_BUFFER | ☓        | KINU_RESIZE_WORKER_SIZE * 3 | Integer                                                                               |                                                                                    |
| KINU_STORAGE_TYPE              | ◯        | none                        | File / S3                                                                             |                                                                                    |
| KINU_FILE_DIRECTORY            | ☓        | none                        | directory path                                                                        | When the `File` has been set in a `KINU_STORAGE_TYPE`\ you must set this variable. |
| KINU_S3_REGION                 | ☓        | none                        | AWS Region                                                                            | When the `S3` has been set in a `KINU_STORAGE_TYPE`\ you must set this variable.   |
| KINU_S3_BUCKET                 | ☓        | none                        | Amazon S3 bucket                                                                      | When the `S3` has been set in a `KINU_STORAGE_TYPE`\ you must set this variable.   |
| KINU_S3_BUCKET_BASE_PATH       | ☓        | none                        |                                                                                       |                                                                                    |
| AWS_ACCESS_KEY_ID              | △        | none                        |                                                                                       | Compliance with the specifications of the aws-sdk-go package.                      |
| AWS_SECRET_ACCESS_KEY          | △        | none                        |                                                                                       | Compliance with the specifications of the aws-sdk-go package.                      |

### Directory structure of the image storage.

now writing

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/tokubai/kinu

## Thanks

- [Issei Naruta](https://github.com/mirakui)
  - The [Tofu](http://www.slideshare.net/mirakui/ss-8150494), which inspired Kinu, is created by him.

## Author

- [Takatoshi Maeda](https://github.com/TakatoshiMaeda)
- [Shinya Watanabe](https://github.com/wata-gh)
