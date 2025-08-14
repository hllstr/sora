# Sora (空) - Base Bot WhatsApp Golang

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue.svg)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Sora (空)** adalah sebuah **Base Bot WhatsApp** yang dibangun menggunakan **Go (Golang)** dan library **Whatsmeow**.

Proyek ini dirancang sebagai "kerangka dasar" yang simpel dan super ringan. Tujuannya adalah menyediakan fondasi yang bersih dan mudah dipahami, sehingga kamu bisa langsung fokus mengembangkan fitur-fitur unikmu sendiri tanpa harus pusing dengan kerumitan setup awal.

## ✨ Filosofi & Fitur

* **Simplicity**: Kode ditulis sesederhana mungkin, biar gampang dioprek.
* **High Performance**: Dirancang untuk sekencang mungkin dengan konsumsi memori yang minim.
* **Extensibility**: Sangat mudah dikembangkan. Tinggal tambah file baru untuk menambah command baru (*plug-and-play*).
* **Dukungan Grup LID**: Bot dapat beroperasi dengan normal di dalam grup yang anggotanya sudah menggunakan format JID LIDs (`@lid`), memastikan kompatibilitas dengan pembaruan terbaru WhatsApp.

## 💻 Bisa Jalan di Mana Saja!

Salah satu keunggulan utama Sora adalah portabilitasnya. Asalkan Go sudah terinstall, script ini bisa berjalan lancar di berbagai platform:
* **Termux**: Langsung di HP Android kamu.
* **Panel Pterodactyl**: Cocok untuk hosting panel yang speknya "kentang" sekalipun.
* **VPS**: Tentu saja, dari VPS murah sampai yang paling gacor.

## 🚀 Instalasi & Menjalankan

1.  **Clone repository ini:**
    ```bash
    git clone https://github.com/hllstr/sora.git
    cd sora
    ```

2.  **Buat file konfigurasi `.env`:**
    Buat file bernama `.env` di direktori utama, lalu isi seperti contoh di bawah.
    ```env
    # Nomor WhatsApp yang akan dijadikan bot (contoh: 6281234567890)
    NUMBER=""

    # Prefix yang ingin digunakan, pisahkan dengan koma
    PREFIXES="/,.,!"

    # Mode bot: "self" (hanya merespons nomormu) atau "public"
    MODE="self"
    ```

3.  **Install dependensi:**
    ```bash
    go mod tidy
    ```

4.  **Jalankan bot:**
    ```bash
    go run .
    ```
    Pada kali pertama, kamu akan diminta untuk melakukan pairing (memasukkan nomor telepon untuk Pairing Code atau mengetik `qr` untuk Kode QR). Sesi akan disimpan di dalam folder `session/`.

## 🧩 Menambah Command Baru

Menambah fitur baru itu gampang banget.

1.  Buat file baru di dalam folder `commands/`, misalnya `halo.go`.
2.  Isi dengan kodemu, gunakan `func init()` untuk mendaftarkan command.

**Contoh: `commands/halo.go`**
```go
package commands

import "fmt"

func init() {
    Plugin(Cmd{
        Name:  "halo",
        Alias: []string{"hi"},
        Desc:  "Bot akan menyapamu kembali.",
        Exec: func(ctx *CommandContext) {
            namaPengirim := ctx.Message.Info.PushName
            balasan := fmt.Sprintf("Halo juga, %s!", namaPengirim)
            ctx.Reply(balasan)
        },
    })
}
```
**Selesai!** Tanpa mengubah file lain, command `.halo` atau `.hi` akan otomatis aktif saat kamu menjalankan bot.

## 🙏 Thanks To

* **Ginko/Soursop**: Atas inspirasi dan arahan di awal-awal belajar `whatsmeow`. Cek profilnya di **[GitHub/ginkohub](https://github.com/ginkohub)**.
* **blog.eikarna.my.id**: Atas bimbingan di fase awal pengembangan, terutama soal praktik terbaik dalam menggunakan package.

### Check Out These Bots!
Jika kamu mencari referensi atau bot lain, pastikan untuk melihat proyek-proyek ini:
* **[Mushi](https://github.com/ginkohub/mushi)** by Ginko
* **[Katsumi](https://github.com/nat9h/Katsumi)** by nat9h

## Hosting Rekomendasi ☁️

Butuh tempat untuk menjalankan bot ini 24/7? Gw pribadi pake dan sangat merekomendasikan panel Pterodactyl murah dan gacor dari **Caliph Cloud**.

🔗 **Link: [https://caliph.cloud/](https://caliph.cloud/)**

## 🤝 Kontribusi

Kontribusi sangat diterima! Jika kamu punya ide untuk fitur baru, perbaikan, atau menemukan bug, silakan buat *issue* atau *pull request* di repository **[github.com/hllstr/sora](https://github.com/hllstr/sora)**.

## 📜 Lisensi

Proyek ini dilisensikan di bawah **MIT License**.
