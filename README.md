<div align="center">
    <h1>Sora (空) - Base Bot WhatsApp Golang</h1>
    <img
        src="https://i.pinimg.com/originals/3b/a1/b1/3ba1b1bec48437f074074ce70e95a06f.jpg"
        alt="Sora"
        style="display:block; margin:-100px 20px; width:65%;"
    />
</div>


# **Sora (空)** adalah sebuah **Base Bot WhatsApp** yang dibangun menggunakan **Go (Golang)** dan library **Whatsmeow**.

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

## 🛠️ Prerequisites

* **Go v1.24.6**: Download & Install ini di https://go.dev/dl (disono ada cara installnya) 

## 🚀 Instalasi & Menjalankan

1.  **Clone repository ini:**
    ```bash
    git clone https://github.com/hllstr/sora.git
    cd sora
    ```

2.  **Copy .env.example dan edit file `.env`:**
    ```bash
    cp .env.example .env
    ```
    Edit isinya sesuai keinginanmu
    ```env
    # Nomor WhatsApp yang akan dijadikan bot (contoh: 6281234567890)
    NUMBER="6281234567890"

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

### Setup Script (kalo mau install otomatis)
**Setelah Clone, Kalian bisa langsung pake script `setup.sh` yang udah gw siapin untuk install di panel, termux, ataupun vps**
```bash
bash setup.sh
```


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
* **blog.eikarna.my.id**: Atas arahan pas awal belajar bikin bot pake whatsmeow.

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
