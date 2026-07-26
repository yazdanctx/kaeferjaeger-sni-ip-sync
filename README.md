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
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.24.2.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile
```

### 2. Clone and build

```bash
cd ~/kaeferyeager
git clone https://github.com/YOUR_USERNAME/kaeferyeager-sni-ip-sync.git
cd sni-ip-sync
go build -o sni-ip-sync .
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

```
0 3 * * * cd /home/YOUR_USER/kaeferyeager/sni-ip-sync && ./sni-ip-sync >> /var/log/sni-ip-sync.log 2>&1
```

Replace `YOUR_USER` with your actual Linux username.

### 5. Verify cron is running

```bash
crontab -l
```

## Manual run

```bash
cd ~/kaeferyeager/sni-ip-sync
./sni-ip-sync
```
