# Learning NATS.io with Golang

NATS adalah open source data layer yang simple, aman, dan berkinerja tinggi untuk aplikasi cloud native, IoT, dan arsitektur microservice.

NATS.io (https://nats.io/) memunkinkan aplikasi berkomunikasi dengan aman di semua kombinasi vendor cloud, on-premise, edge, web dan mobile, serta perangkat. ATS terdiri dari rangkaian produk open source yang terintegrasi dengan erat namun dapat digunakan dengan mudah dan mandiri.

Server NATS bertindak sebagai sistem saraf pusat untuk membangun aplikasi terdistribusi. Official client tersedia dalam bahasa Go, Rust, JavaScript (Node dan Web), TypeScript (Deno), Python, Java, C#, C, Ruby, Elixir, dan CLI di samping lebih dari 30 klien yang dikontribusikan oleh komunitas. Streaming data secara real time, penyimpanan data yang sangat tangguh, dan pengambilan data yang fleksibel didukung melalui JetStream, platform streaming generasi berikutnya yang terpasang pada server NATS. Lihat daftar lengkap klien NATS.

Karena kita memakai bahasa pemrograman Go, jadinya bisa kita embed NATS server ke aplikasi Go. Karena NATS support embedding ke bahasa Go. Tapi mending pisah deh

Komponen NATS:
- NATS server / monitoring
- NATS client
- NATS cluster

NATS pub-sub vs request-reply:
- Publish-Subscribe: cocok untuk komunikasi many-to-many atau one-to-many, di mana semua penerima yang relevan mendapatkan pesan yang sama. Contoh penggunaanya broadcast promosi, iklan, notifikasi, event, logginh, dll.
- Request-Reply: cocok untuk komunikasi one-to-one, di mana pengirim membutuhkan jawaban langsung dari penerima tertentu. Contoh penggunaanya, kominkasi chat, transaksi, dll.

## Setup dan running

### Turn on NATS server

```bash
docker compose up -d
docker exec -it CONTAINER nats-server
```

### Run location tracker app

```bash
cd pub-sub/location-track && go run .
```

### Run chat app

```bash
cd request-reply/chat-app && go run .
```