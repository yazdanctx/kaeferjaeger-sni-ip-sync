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
go build -o extract extract.go
```

### 3. Test manually

```bash
./sni-ip-sync
```

### 4. Set up cron

```bash
crontab -e
```

Add this line to run daily at 3:00 AM:

```bash
0 3 * * * cd /root/kaeferyeager/sni-ip-sync && ./sni-ip-sync >> /var/log/sni-ip-sync.log 2>&1
```
### 5. Verify cron is running

```bash
crontab -l
```

## Manual run

```bash
cd /root/kaeferyeager/sni-ip-sync
./sni-ip-sync
```

## Extract domains

The `extract` tool pulls unique subdomains for a given domain from `final.txt`.

### Usage

```bash
# Extract domains for a specific domain (reads from default path)
./extract example.com

# Extract with custom input file
./extract -f /path/to/final.txt example.com
```

### Output

Files are saved in the current directory with the format:

```
x_example.com_2026-07-26_143052.txt
```
