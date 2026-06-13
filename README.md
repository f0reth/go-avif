# go-avif

## 実行に必要なファイル

このライブラリは [purego](https://github.com/ebitengine/purego) を使って libavif の共有ライブラリを動的に読み込みます。事前に以下のファイルをインストールしてください（libavif 1.1.0 以上が必要です）。

| OS | 必要なファイル |
|---|---|
| Windows | `libavif.dll` |
| Linux | `libavif.so` |
| macOS | `libavif.dylib` |

## 使い方

### インポート

```go
import "github.com/f0reth/go-avif"
```

### デコード（AVIF → image.Image）

```go
f, _ := os.Open("input.avif")
defer f.Close()

img, err := avif.Decode(f)
```

### アニメーション AVIF のデコード

```go
f, _ := os.Open("anim.avifs")
defer f.Close()

ret, err := avif.DecodeAll(f)
// ret.Image: 各フレームの画像スライス
// ret.Delay: 各フレームの表示時間（秒）
```

### エンコード（image.Image → AVIF）

```go
f, _ := os.Create("output.avif")
defer f.Close()

err := avif.Encode(f, img)
```

### エンコードオプション

```go
err := avif.Encode(f, img, avif.Options{
    Quality:           80,  // 画質 0〜100（100 で可逆、デフォルト: 60）
    QualityAlpha:      80,  // アルファの画質 0〜100
    Speed:             8,   // エンコード速度 0（高品質）〜10（速い、デフォルト: 10）
    ChromaSubsampling: image.YCbCrSubsampleRatio420, // クロマサブサンプリング 444|422|420
})
```

### ライブラリの読み込み確認

```go
if err := avif.Dynamic(); err != nil {
    log.Fatal("共有ライブラリが見つかりません:", err)
}
```
