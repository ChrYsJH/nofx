# 如何在 Linux 服务器上部署 NOFX (支持 AsterDex 测试网)

本指南将帮助你将 NOFX 服务部署到 Linux 服务器，并配置 AsterDex 测试网以进行低风险测试。

## 1. 准备工作 (Prerequisites)

请确保你的服务器安装了以下软件：

*   **Linux OS**: Ubuntu 22.04 LTS 或其他主流发行版
*   **Go**: 1.21 或更高版本
*   **Node.js**: 18.0 或更高版本 (用于构建前端)
*   **Git**: 用于克隆代码

### 安装 Go (如果未安装)
```bash
wget https://go.dev/dl/go1.22.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version
```

### 安装 Node.js (如果未安装)
```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs
node -v
```

## 2. 获取代码并构建

1.  **克隆仓库**
    ```bash
    git clone https://github.com/your-username/nofx.git
    cd nofx
    ```

2.  **配置环境变量**
    复制示例配置文件：
    ```bash
    cp .env.example .env
    ```
    编辑 `.env` 文件 (使用 `nano .env` 或 `vim .env`)：
    *   修改 `JWT_SECRET` 为一个强大的随机字符串。
    *   如果需要远程访问 API，确保 `TRANSPORT_ENCRYPTION=false` (除非你配置了 HTTPS)。
    *   (可选) 配置数据库路径或 PostgreSql 连接信息。默认使用 SQLite。

3.  **构建后端**
    ```bash
    go mod download
    go build -o nofx main.go
    ```
    构建成功后，你会看到可以在当前目录下看到可执行文件 `nofx`。

4.  **构建前端**
    ```bash
    cd web
    npm install
    npm run build
    cd ..
    # 前端构建产物通常在 web/dist 目录，后端会服务于该目录或你需要配置 Nginx 反向代理
    ```

## 3. 设置 Systemd 服务 (推荐)

为了让服务在后台运行并开机自启，建议创建一个 Systemd 服务文件。

1.  **创建服务文件**
    ```bash
    sudo nano /etc/systemd/system/nofx.service
    ```

2.  **粘贴以下内容** (请修改 `User`, `Group`, `WorkingDirectory` 为你的实际路径)：
    ```ini
    [Unit]
    Description=NOFX Trading System
    After=network.target

    [Service]
    User=root
    Group=root
    WorkingDirectory=/root/nofx
    ExecStart=/root/nofx/nofx
    Restart=always
    RestartSec=5
    EnvironmentFile=/root/nofx/.env

    [Install]
    WantedBy=multi-user.target
    ```
    *注意：请将 `/root/nofx` 替换为你实际的代码存放路径，User 替换为实际用户。*

3.  **启动服务**
    ```bash
    sudo systemctl daemon-reload
    sudo systemctl enable nofx
    sudo systemctl start nofx
    ```

4.  **查看状态**
    ```bash
    sudo systemctl status nofx
    # 查看日志
    journalctl -u nofx -f
    ```

## 4. 配置 AsterDex 测试网

服务启动后，你需要通过 Web 界面进行配置。

1.  **访问 Web 界面**
    打开浏览器访问 `http://<服务器IP>:<端口>` (默认端口是 8080)。

2.  **登录/注册**
    如果启用了注册功能，创建一个管理员账号。

3.  **配置交易所 (Exchange Config)**
    *   进入 "Exchange" 或 "Settings" 页面。
    *   点击 "Add Exchange" (添加交易所)。
    *   **Type**: 选择 `aster`。
    *   **Use Testnet (测试网)**: **勾选此选项 (Enable Testnet)**。这非常重要，勾选后系统将连接到 `https://testnet-fapi.asterdex.com`。
    *   **API Wallet Address (Signer)**: 输入你的 AsterDex API 钱包地址。
    *   **Main Wallet Address (User)**: 输入你的主钱包地址。
    *   **Private Key**: 输入 API 钱包的私钥 (Hex 格式)。

4.  **创建 AI 交易员 (Create Trader)**
    *   进入 "Traders" 页面。
    *   点击 "Create Trader"。
    *   选择刚才配置的 Aster 交易所账户。
    *   选择 AI 模型 (如 DeepSeek, Qwen 等)。
    *   启动交易员，AI 将会在测试网环境中开始运行。

## 5. 常见问题

*   **端口无法访问**: 检查服务器防火墙 (UFW/iptables) 是否开放了 8080 端口。
    ```bash
    sudo ufw allow 8080
    ```
*   **数据库错误**: 确保 `data/` 目录存在且有写入权限 (如果是 SQLite)。

现在，你应该可以在安全的环境下使用 AI 在 AsterDex 测试网上进行交易测试了！
