# sni-ip-sync

Go script that syncs SNI IP range data from [kaeferjaeger.gay](https://kaeferjaeger.gay/?dir=sni-ip-ranges) every 24 hours via cron.

## What it does

1. Downloads `ipv4_merged_sni.txt` from each provider (amazon, digitalocean, google, microsoft, oracle)
2. Stores them in `data/{provider}/`
3. Merges all files into a single `final.txt`

## Setup on Linux

### 1. Install Go

```bash
wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.2.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile
```

### 2. Clone and build

```bash
cd /root
git clone https://github.com/YOUR_USERNAME/kaeferyeager-sni-ip-sync.git
cd sni-ip-sync
go build -o sni-ip-sync sync.go
```

### 3. Set up cron

```bash
crontab -e
```

Add this line to run daily at 3:00 AM:

```bash
0 3 * * * cd /root/kaeferyeager/sni-ip-sync && ./sni-ip-sync >> /var/log/sni-ip-sync.log 2>&1
```

### 4. Verify cron is running

```bash
crontab -l
```

## Usage

### Sync (download and merge data)

```bash
./sni-ip-sync
```

This must be run first to generate `final.txt` before extracting.

### Extract subdomains for a domain

```bash
./extract.sh dell.com
```

This reads `final.txt`, finds all entries matching `.dell.com`, extracts the unique subdomains, and writes them to `output/dell_com_<timestamp_ms>.txt`.

### Full workflow

```bash
# 1. Sync the latest data
./sni-ip-sync

# 2. Extract subdomains for a domain
./extract.sh dell.com

# 3. Check the output
cat output/dell_com_*.txt
```
