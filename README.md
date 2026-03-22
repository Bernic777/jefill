JefillJefill ni tool CLI simpel pake Golang untuk simpan deng kasih jalan ko pu pattern awk atau shell command andalan. Ini cocok skali untuk anak-anak bug bounty supaya tra usah ketik ulang atau copas command awk yang panjang macam gerbong kereta tiap kali lagi recon.       _       __ _  _ _ 
      (_)     / _(_)| | |
       _  ___| |_ _ | | |
      | |/ _ \  _| || | |
      | |  __/ | | || | |
      | |\___|_| |_||_|_|
     _/ |                
    |__/       v2.0.0    
FiturAtur ko pu file pattern awk langsung di direktori ~/.config/jefill/.Eksekusi cepat skali, cuma modal satu command saja (contoh: jefill uniqueparam).Support full Linux pipeline (stdin deng stdout), mantap untuk ko chain.Gampang skali tambah pattern baru—ko tinggal bikin file .txt saja di folder config.Cara InstallKo pastikan su install Go di ko pu mesin (lokal atau VPS). Kalo su beres, ko hantam command ni:go install [github.com/username-github-anda/jefill@latest](https://github.com/username-github-anda/jefill@latest)
(Note: pastikan $GOPATH/bin su masuk di ko pu environment PATH e, supaya command-nya bisa dipanggil dari mana saja).Konfigurasi & Nambah PatternPas ko kasih jalan Jefill untuk pertama kali, tool ni nanti otomatis bikin folder ~/.config/jefill/ sekalian generate file uniqueparam.txt di dalam situ.Isi default dari ~/.config/jefill/uniqueparam.txt tu command ni:awk -F'[?&]' '{split($1,d,"/");dom=d[3];for(i=2;i<=NF;i++){split($i,a,"=");k=dom":"a[1];if(!seen[k]++){print $1"?"a[1]"="}}}' | sort -u
Kalo ko mau tambah command awk atau script bash lain, ko tinggal bikin file text baru saja di folder tu.Contohnya, ko bikin file ~/.config/jefill/getjs.txt yang isinya grep "\.js$", nanti ko tinggal panggil pake command jefill getjs.Cara PakeJefill ni su didesain untuk terima input dari pipeline (stdin) terus eksekusi pattern sesuai deng nama filenya.Liat list pattern yang ko su save:jefill list
1. Baca input dari file text (contoh pake pattern 'uniqueparam'):cat katana.txt | jefill uniqueparam
2. Di-chain langsung deng tool recon lain (contoh dari Katana lempar pi Dalfox):katana -u [https://target.com](https://target.com) -d 5 | jefill uniqueparam | dalfox pipe
