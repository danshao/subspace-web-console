# 開發階段注意事項
## 在 Local 開發，但是資料連到遠端 EC2 時

1. 修改 `conf/app.conf`
   - RunMode 從 `prod` 改為 `dev`
   - [dev] 區塊內的 host 從 `PUBLIC_INSTANCE_IP` 改為要連的遠端的 EC2 public IP

2. Initial Setup 的 Instance ID 填 12345678。

3. AWS 上 EC2 的 security group 要開以下 port 允許自己的 IP 連入。

   | 用途                            | Protocol | Port |
   | ----------------------------- | -------- | ---- |
   | SSH                           | TCP      | 22   |
   | MySQL                         | TCP      | 3306 |
   | Redis                         | TCP      | 6379 |
   | Soft ether Server Manage Tool | TCP      | 992  |

4. 若使用 release 版本，EC2 要修改 MySQL 和 Redis 設定去接受外部連線

  - MySQL

     1. 在 EC2 上修改 `vi /etc/mysql/mysql.conf.d/mysqld.cnf`，將 `bind-address           = 127.0.0.1` 這行前面加上 `#` 註解起來
     2. 開啟 Security Group 的 3306 port 允許自己 IP

  - Redis

     1. 在 EC2 上修改 `vi /etc/redis/redis.conf`，將 `bind 127.0.0.1` 這行前面加上 `#` 註解起來
     2. 開啟 Security Group 的 3306 port 允許自己 IP

5. 在 Local [安裝 Beego 環境](#安裝 Beego 環境)




# EC2 上 Subspace 開發環境安裝方式

## 安裝 Go 1.8

### 下載 Go

https://github.com/golang/go/wiki/Ubuntu

```sh
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt-get update
sudo apt-get install golang-go
```

### 設定 Go 環境變數

執行 `vi ~/.profile` 編輯 bash 的 profile，將 path 該段內容改為以下：

```shell
# set PATH so it includes user's private bin directories
GOPATH="$HOME/go"
PATH="$GOPATH/bin:$HOME/bin:$HOME/.local/bin:$PATH"
```

完成以後執行以下指令讓設定生效。

```sh
source ~/.profile
```



## 主程式安裝

### 下載 Subspace 相關 package

```sh
mkdir ~/go/src/gitlab.ecoworkinc.com/Subspace && cd ~/go/src/gitlab.ecoworkinc.com/Subspace
git clone https://gitlab.ecoworkinc.com/Subspace/web-console
git clone https://gitlab.ecoworkinc.com/Subspace/subspace-utility
git clone https://gitlab.ecoworkinc.com/Subspace/softetherlib
git clone https://gitlab.ecoworkinc.com/Subspace/vpn-profile-generator
git clone https://gitlab.ecoworkinc.com/Subspace/server-status-api.git
```

### 下載 Web Console 需要的 go package

```sh
cd ~/go/src/gitlab.ecoworkinc.com/Subspace/web-console
go get
```



## 安裝 Beego 環境

```sh
go get -u github.com/astaxie/beego
go get -u github.com/beego/bee
```

### 測試執行

如果在 EC2 上要執行，要先執行 `sudo service subspace stop` 把原本在跑的 subspace 停掉，否則 port 要調開。

```sh
cd ~/go/src/gitlab.ecoworkinc.com/Subspace/web-console
bee run
```

## 個人習慣的 git 設定

```
git config --global alias.co checkout
git config --global alias.ci commit
git config --global alias.st status
git config --global alias.br branch
git config --global alias.lgo 'log --graph --decorate'
git config --global core.editor "vim"
```

