# Jefill

Jefill ni tool CLI simpel pake Golang untuk simpan deng kasih jalan ko pu pattern `awk` andalan. Ini cocok skali untuk anak-anak bug bounty supaya tra usah ketik ulang atau copas command `awk` yang panjang macam gerbong kereta tiap kali lagi recon.

```
       _       __ _  _ _ 
      (_)     / _(_)| | |
       _  ___| |_ _ | | |
      | |/ _ \  _| || | |
      | |  __/ | | || | |
      | |\___|_| |_||_|_|
     _/ |                
    |__/       v1.0.0    
```

## Inspirasi & Keamanan (Anti RCE)

Jefill ni terinspirasi dari tool mantap puang **[tomnomnom/gf](https://github.com/tomnomnom/gf)** (*A wrapper around grep*). Tapi Jefill ambil jalan beda, Jefill ni **khusus untuk pattern `awk` saja** tanpa campur tool lain.

Kenapa sa bikin begini? Supaya super aman, kawan! Kalo ada hacker atau orang luar yang tawarkan file pattern Jefill gratis di internet, dorang tra bisa selipkan script berbahaya. File `pattern.txt` di Jefill tu murni *hanya argumen pattern awk saja* (tra usah ketik kata `awk` lagi di dalamnya). Jadi dorang tra bisa tambah pipe (`|`) atau chain command (macam `&& reverse_shell`). Jefill pastikan ko pu mesin tetap aman 100% waktu pake pattern dari luar!

## Fitur

* Atur ko pu file pattern `awk` langsung di direktori `~/.config/jefill/`.
* Eksekusi cepat skali, cuma modal satu command saja (contoh: `jefill uniqueparam`).
* Support full Linux pipeline (`stdin` deng `stdout`), mantap untuk ko chain.
* Super aman dari RCE karena file konfigurasi cuma baca argumen spesifik eksekusi awk.

## Cara Install

Ko pastikan su install Go di ko pu mesin (lokal atau VPS). Kalo su beres, ko hantam command ni:

```bash
go install github.com/Bernic777/jefill@latest
```

*(Note: pastikan `$GOPATH/bin` su masuk di ko pu environment `PATH` e, supaya command-nya bisa dipanggil dari mana saja).*

## Konfigurasi & Nambah Pattern

Pas ko kasih jalan Jefill untuk pertama kali, tool ni nanti otomatis bikin folder `~/.config/jefill/` sekalian generate file `uniqueparam.txt` di dalam situ.

Isi default dari `~/.config/jefill/uniqueparam.txt` tu pattern ni:

```text
-F'[?&]' '{split($1,d,"/");dom=d[3];for(i=2;i<=NF;i++){split($i,a,"=");k=dom":"a[1];if(!seen[k]++){print $1"?"a[1]"="}}}'
```
*(Ingat e: tra pake kata `awk` di depan, dan tra pake `| sort -u` di belakang, murni cuma argumen pattern saja!)*

Kalo ko mau tambah pattern `awk` lain, ko tinggal bikin file text baru saja di folder tu. 

## Contoh Pattern Awk Lain yang Mantap

Ko bisa bikin pattern `awk` apa saja bebas. Ini ada beberapa contoh pattern sakti yang biasa anak-anak bug bounty pake. Ko cukup bikin filenya di folder `~/.config/jefill/` pake nama yang ko suka.

**1. `~/.config/jefill/getjs.txt`** (Ambil URL JavaScript saja)
```text
'/\.js(\?|$)/'
```
*Cara panggil: `cat urls.txt | jefill getjs`*

**2. `~/.config/jefill/xss.txt`** (Saring parameter yang rawan kena XSS)
```text
'/q=|query=|search=|id=|name=|keyword=/'
```
*Cara panggil: `cat urls.txt | jefill xss`*

**3. `~/.config/jefill/lfi.txt`** (Saring parameter yang rawan kena LFI/SSRF)
```text
'/file=|path=|page=|dir=|folder=|doc=|url=/'
```
*Cara panggil: `cat urls.txt | jefill lfi`*

**4. `~/.config/jefill/admin.txt`** (Saring parameter sensitif macam admin/debug/config)
```text
'tolower($0) ~ /access=|admin=|dbg=|debug=|edit=|grant=|test=|alter=|clone=|create=|delete=|disable=|enable=|exec=|execute=|load=|make=|modify=|rename=|reset=|shell=|toggle=|adm=|root=|cfg=|config=/'
```
*(Catatan: `tolower($0)` dipake supaya dia case-insensitive mirip macam flag `-iE` di grep)*
*Cara panggil: `cat urls.txt | jefill admin`*

## Cara Pake

Jefill ni su didesain untuk terima input dari *pipeline* (stdin) terus eksekusi pattern sesuai deng nama filenya.

**Liat list pattern yang ko su save:**

```bash
jefill list
```

**1. Baca input dari file text (contoh pake pattern 'uniqueparam'):**

```bash
cat katana.txt | jefill uniqueparam
```

**2. Di-chain langsung deng tool recon lain (contoh dari Katana lempar pi Dalfox):**

```bash
katana -u https://target.com -d 5 | jefill uniqueparam | dalfox pipe
```
