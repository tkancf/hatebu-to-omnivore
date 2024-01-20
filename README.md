# hatebu-to-omnivore

はてなブックマークからエクスポートしたデータを Omnivore へインポートするためのプログラム

## Install

```sh
$ go install github.com/tkancf/hatebu-to-omnivore@latest
```

## Usage

### 1. はてなブックマークからデータをエクスポート

1. はてなブックマークにログイン
2. [はてなブックマークのデータ管理 設定ページ](https://b.hatena.ne.jp/-/my/config/data_management) にアクセス
3. Atomフィード形式でエクスポート

### 2. Omnivore のAPIキーを取得

1. Omnivore にログイン
2. [APIキーの設定ページ](https://omnivore.app/settings/api) にアクセス

### 3. Omnivore へインポート

#### `hatebu-to-omnivore` にオプションを指定してインポートを実行

- 下記は必要最低限のオプションを指定した例
    - OmnivoreのAPI URLはデフォルトの (https://api-prod.omnivore.app/api/graphql) を指定
    - インポートされた記事はinboxに入る

```sh
$ hatebu-to-omnivore -i <エクスポートしたAtomフィードデータ> -k <OmnivoreのAPIキー>
```

- 下記は全てのオプションを指定した例
    - `-i`: エクスポートしたAtomフィードデータを指定
    - `-a`: インポートされた記事をARCHIVEする
    - `-k`: OmnivoreのAPIキーを指定
    - `-u`: OmnivoreのAPI URLを指定 (Omnivoreをセルフホストしている場合に指定する)

```sh
$ hatebu-to-omnivore -i <エクスポートしたAtomフィードデータ> -a -k <OmnivoreのAPIキー> -u <OmnivoreのAPI URL>
```
