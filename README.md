# Cyber Security News Scraper 

Go dilinde yazılmış, çeşitli siber güvenlik haber sitelerinden otomatik olarak haber çeken web scraper uygulaması.

## Özellikler 

- **Üç farklı siber güvenlik kaynağından haber çekme:**
  - The Hacker News
  - Krebs on Security
  - Cyber Security News
  
- **Chromedp ile dinamik sayfa yükleme** - JavaScript ile yüklenen içeriği destekler
- **Dosyaya Kaydetme** - Tüm haberler otomatik olarak `SiteName_DD-MM-YYYY.txt` formatında kaydedilir
- **Filtreleme:**
  - `--date` - Tarihleri gizler
  - `--description` - Açıklamaları gizler

## Kurulum 

### Gereksinimler
- Go 1.21 veya üstü
- Chrome/Chromium tarayıcı (chromedp tarafından kullanılır)

### Adımlar

```bash
# Repo'yu klonla
git clone https://github.com/taro544/newscraper.git
cd go_scraper

# Bağımlılıkları yükle
go mod download

# Kodu derle
go build -o scraper .
```

## Kullanım 

```bash
# The Hacker News'ten haberler çek
go run scraper.go -1

# Krebs on Security'den haberler çek
go run scraper.go -2

# Cyber Security News'den haberler çek
go run scraper.go -3

# Tarih bilgisini gizleme
go run scraper.go -3 --date

# Açıklamaları gizleme
go run scraper.go -3 --description

# Hem tarih hem açıklama gizleme
go run scraper.go -1 --date --description

# Uygulamadan çık
go run scraper.go -4
```

## Terminal Çıktısı Örneği 

```
==================================================
Site: Cyber Security News
Cekilme Tarihi: 16-04-2026 21:05:33
==================================================

1.Haber
Baslik: SpankRAT Exploits Windows Explorer Processes
Aciklama:
A newly identified two-component Remote Access Trojan (RAT)...
Tarih:
April 16, 2026
--------------------------------------------------
```


## Çıktı Dosyaları 

Her çalıştırmada haberler şu formatta kaydedilir:
- `The_Hacker_News_16-04-2026.txt`
- `Krebs_on_Security_16-04-2026.txt`
- `Cyber_Security_News_16-04-2026.txt`



### Kullanılan Kütüphaneler
- **chromedp** - Headless Chrome/Chromium otomasyonu
- **goquery** - HTML parsing ve CSS selectors


---

