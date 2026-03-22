package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func printBanner() {
	banner := `
       _       __ _  _ _ 
      (_)     / _(_)| | |
       _  ___| |_ _ | | |
      | |/ _ \  _| || | |
      | |  __/ | | || | |
      | |\___|_| |_||_|_|
     _/ |                
    |__/       v2.0.0    
`
	fmt.Fprintln(os.Stderr, banner)
}

// parseArgs berfungsi memecah string dari file text menjadi array argumen
// layaknya shell, tapi murni di dalam Go. Ini memastikan tool aman 
// dari Command Injection/RCE.
func parseArgs(input string) []string {
	var args []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	escapeNext := false

	for _, r := range input {
		if escapeNext {
			current.WriteRune(r)
			escapeNext = false
			continue
		}
		if r == '\\' {
			escapeNext = true
			continue
		}
		if r == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}
		if r == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}
		if r == ' ' || r == '\n' || r == '\t' {
			if !inSingleQuote && !inDoubleQuote {
				if current.Len() > 0 {
					args = append(args, current.String())
					current.Reset()
				}
				continue
			}
		}
		current.WriteRune(r)
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}

func main() {
	// [DOUBLE CHECK 1] Pastikan 'awk' terinstall di sistem
	_, err := exec.LookPath("awk")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[!] Error: 'awk' tidak ditemukan di sistem ko. Pastikan awk su terinstall!")
		os.Exit(1)
	}

	// Dapatkan direktori home user
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "[!] Error dapet direktori home:", err)
		os.Exit(1)
	}

	// Tetapkan laluan folder config
	configDir := filepath.Join(homeDir, ".config", "jefill")

	// Cipta folder ~/.config/jefill/ jika belum wujud
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "[!] Error bikin folder config:", err)
		os.Exit(1)
	}

	// Cipta fail corak lalai "uniqueparam.txt" jika tidak wujud
	defaultPatternFile := filepath.Join(configDir, "uniqueparam.txt")
	if _, err := os.Stat(defaultPatternFile); os.IsNotExist(err) {
		defaultAwk := `-F'[?&]' '{split($1,d,"/");dom=d[3];for(i=2;i<=NF;i++){split($i,a,"=");k=dom":"a[1];if(!seen[k]++){print $1"?"a[1]"="}}}'`
		os.WriteFile(defaultPatternFile, []byte(defaultAwk), 0644)
	}

	// Paparkan banner dan bantuan jika argumen kurang
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		printBanner()
		fmt.Fprintln(os.Stderr, "Cara Pake:")
		fmt.Fprintln(os.Stderr, "  cat urls.txt | jefill <nama_pattern>")
		fmt.Fprintln(os.Stderr, "  jefill <nama_pattern> urls.txt")
		fmt.Fprintln(os.Stderr, "  jefill -list   (untuk liat list pattern yg ada)")
		fmt.Fprintln(os.Stderr, "  jefill -update (untuk download pattern terbaru dari GitHub)")
		os.Exit(1)
	}

	patternName := os.Args[1]

	// Jika pengguna mahu melihat list pattern
	if patternName == "-list" {
		printBanner()
		fmt.Fprintln(os.Stderr, "List Pattern di ~/.config/jefill/:")
		files, err := os.ReadDir(configDir)
		if err != nil {
			fmt.Fprintln(os.Stderr, "[!] Error baca folder config:", err)
			os.Exit(1)
		}

		count := 0
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
				name := strings.TrimSuffix(file.Name(), ".txt")
				fmt.Fprintf(os.Stderr, "  - %s\n", name)
				count++
			}
		}
		if count == 0 {
			fmt.Fprintln(os.Stderr, "  (Belum ada file pattern .txt)")
		}
		os.Exit(0)
	}

	// Jika pengguna mahu update pattern dari repo Github tanpa butuh GIT
	if patternName == "-update" {
		printBanner()
		fmt.Fprintln(os.Stderr, "[i] Menarik pattern sakti dari github.com/Bernic777/jefill-pattern tanpa git...")

		// URL zip dari Github branch main
		zipURL := "https://github.com/Bernic777/jefill-pattern/archive/refs/heads/main.zip"

		// Download file zip
		resp, err := http.Get(zipURL)
		if err != nil {
			fmt.Fprintln(os.Stderr, "[!] Error waktu download pattern. Pastikan internet aman:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Fprintf(os.Stderr, "[!] Error: GitHub kasih status %d. Pastikan repo su public dan branch main wujud.\n", resp.StatusCode)
			os.Exit(1)
		}

		// Bikin file temporary untuk simpan zip
		tempZip, err := os.CreateTemp("", "jefill-pattern-*.zip")
		if err != nil {
			fmt.Fprintln(os.Stderr, "[!] Error bikin file temp zip:", err)
			os.Exit(1)
		}
		defer os.Remove(tempZip.Name()) // Bersihkan file zip kalau su selesai

		// Tulis hasil download ke file temp zip
		if _, err := io.Copy(tempZip, resp.Body); err != nil {
			fmt.Fprintln(os.Stderr, "[!] Error simpan isi zip:", err)
			os.Exit(1)
		}
		tempZip.Close()

		// Buka zip untuk di-ekstrak
		r, err := zip.OpenReader(tempZip.Name())
		if err != nil {
			fmt.Fprintln(os.Stderr, "[!] Error buka file zip:", err)
			os.Exit(1)
		}
		defer r.Close()

		// Proses ekstrak khusus file .txt
		count := 0
		for _, f := range r.File {
			// Kita tra ambil folder, dan cuma ambil yang belakangnya .txt
			if !f.FileInfo().IsDir() && strings.HasSuffix(f.Name, ".txt") {
				// Github zip biasanya ada folder utama (misal: jefill-pattern-main/xss.txt)
				// Kita cuma ambil nama filenya saja (xss.txt) pakai filepath.Base
				baseName := filepath.Base(f.Name)
				dstPath := filepath.Join(configDir, baseName)

				// Buka file di dalam zip
				rc, err := f.Open()
				if err != nil {
					continue
				}

				// Cipta file di lokal ~/.config/jefill/
				dstFile, err := os.Create(dstPath)
				if err == nil {
					io.Copy(dstFile, rc) // Salin isinya
					dstFile.Close()
					count++
				}
				rc.Close()
			}
		}

		fmt.Fprintf(os.Stderr, "[v] Mantap kawan! %d pattern su berhasil ditambah/diupdate ke ~/.config/jefill/\n", count)
		os.Exit(0)
	}

	// Baca fail corak (pattern) yang diminta
	patternFile := filepath.Join(configDir, patternName+".txt")
	contentBytes, err := os.ReadFile(patternFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Error: Pattern '%s' tra ketemu kawan.\nBikin dulu file %s\n", patternName, patternFile)
		os.Exit(1)
	}

	// Ambil string dari file pattern
	patternStr := strings.TrimSpace(string(contentBytes))
	
	// [DOUBLE CHECK 2] Pastikan file pattern tidak kosong
	if patternStr == "" {
		fmt.Fprintf(os.Stderr, "[!] Error: File %s kosong, tra ada isinya!\n", patternFile)
		os.Exit(1)
	}

	// Pecah jadi argumen
	args := parseArgs(patternStr)

	// Laksanakan perintah AWK SECARA LANGSUNG (Anti-RCE)
	cmd := exec.Command("awk", args...)

	// [DOUBLE CHECK 3] Support Ekstra Argumen.
	// Kalo ko ketik: jefill uniqueparam urls.txt
	// Maka 'urls.txt' bakal dioper langsung ke awk.
	if len(os.Args) > 2 {
		cmd.Args = append(cmd.Args, os.Args[2:]...)
	}

	// Sambungkan IO
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// [DOUBLE CHECK 4] Peringatan Jika Nunggu Input Terminal
	stat, _ := os.Stdin.Stat()
	// Jika stdin berasal dari terminal (bukan pipe/file) DAN tidak ada file argumen tambahan
	if (stat.Mode()&os.ModeCharDevice) != 0 && len(os.Args) == 2 {
		fmt.Fprintln(os.Stderr, "[i] Menunggu input dari terminal... (Tekan Ctrl+D kalau su selesai ketik)")
	}

	// Jalankan perintah
	if err := cmd.Run(); err != nil {
		// Abaikan error exit status 1 atau 2 dari awk jika memang pattern sengaja di-stop
		// Tapi kalau error lain, baru kita print
		if _, ok := err.(*exec.ExitError); !ok {
			fmt.Fprintf(os.Stderr, "[!] Error jalanin awk: %v\n", err)
			os.Exit(1)
		}
	}
}
