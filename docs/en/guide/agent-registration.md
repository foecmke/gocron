# Agent Auto-Registration

gocron supports one-click generation of installation commands through the web interface. Simply execute the command on the target server to automatically install and register the Agent node. This is the simplest and fastest way to deploy task nodes.

## Features

* ✅ **One-Click Installation** - Copy command, done in one line
* ✅ **Auto Registration** - No manual configuration, automatically connects to gocron server
* ✅ **Secure Token** - One-time token with 3-hour validity
* ✅ **Batch Deployment** - Token can be reused within validity period
* ✅ **System Service** - Automatically creates system service with auto-start
* ✅ **Offline Support** - Supports intranet environments, prioritizes local packages
* ✅ **Cross-Platform** - Supports Linux, macOS, Windows

## Usage

::: tip 💡 Tip: Offline Installation (Optional)
If your environment cannot access GitHub, or you want to speed up installation, you can prepare offline packages first:

1. Create directory in the gocron server installation directory: `mkdir -p gocron-node-package`
2. Download `gocron-node-*.tar.gz` or `gocron-node-*.zip` files for required platforms from [GitHub Releases](https://github.com/gocronx-team/gocron/releases)
3. Place downloaded files in `gocron-node-package/` directory

This way the installation script will download from local server first, no GitHub access needed. See [Offline Installation](#offline-installation) section below for details.
:::

### Step 1: Generate Installation Command

1. Login to gocron web interface
2. Go to "Task Nodes" page
3. Click "Auto Register" button
4. Select target platform (Linux/macOS/Windows)
5. Copy the generated installation command

### Step 2: Execute on Target Server

Choose the installation method based on the target server's operating system:

#### Linux / macOS

```bash
curl -fsSL http://your-server:5920/api/agent/install.sh | bash -s -- <token>
```

**Example**:
```bash
curl -fsSL http://192.168.1.100:5920/api/agent/install.sh | bash -s -- eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Windows

::: warning Note
Windows auto-installation script is not supported for security reasons. Please install manually.
:::

**Manual Installation Steps**:

1. Download Windows package: [gocron-node-windows-amd64.zip](https://github.com/gocronx-team/gocron/releases/latest/download/gocron-node-windows-amd64.zip)
2. Extract to a directory (e.g., `D:\gocron-node`)
3. Open PowerShell
4. Navigate to the extracted directory
5. Run `.\gocron-node.exe` to start the node
6. Manually add the node in the Web Interface (IP and Port)

### Step 3: Verify Installation

After installation completes, return to the "Task Nodes" page in the web interface. You should see the newly registered node.

## Installation Process

For Linux/macOS, the installation script automatically performs the following operations:

1. **Detect System** - Automatically identifies OS and architecture
2. **Download Package** - Downloads corresponding version from server or GitHub
3. **Extract and Install** - Extracts files to designated directory
4. **Register Node** - Registers with server using Token
5. **Create Service** - Creates system service (systemd / launchd)
6. **Start Service** - Automatically starts gocron-node service

### Installation Directory

| OS | Installation Directory |
|----|----------------------|
| Linux | `/opt/gocron-node` |
| macOS | `/usr/local/gocron-node` |

### Configuration File

| OS | Configuration File Path |
|----|------------------------|
| Linux/macOS | `~/.gocron-node/node.ini` |

## Token Mechanism

### Token Features

- **One-time Generation** - New token generated each time "Auto Register" is clicked
- **Validity Period** - Valid for 3 hours
- **Reusable** - Same token can be used for batch installation of multiple nodes
- **Auto Expiration** - Automatically expires after validity period

### Token Security

- Token uses JWT format with encrypted signature
- Token only used for node registration, contains no sensitive information
- Token can be discarded after registration
- Two-way communication between node and server:
  - **Registration**: Node initiates HTTP connection to Server
  - **Task Execution**: Server initiates gRPC connection to Node (ensure Server can access Node's IP:Port)

## Offline Installation

For intranet environments or scenarios where GitHub is inaccessible, gocron supports offline installation.

### Configure Offline Packages

#### 1. Create Package Directory

Create directory in the gocron server installation directory:

```bash
mkdir -p gocron-node-package
```

#### 2. Download Packages

Visit [GitHub Releases](https://github.com/gocronx-team/gocron/releases) to download packages for required platforms.

**Method 1: Manual Download**

Download the following files and place in `gocron-node-package/` directory:
- `gocron-node-linux-amd64.tar.gz`
- `gocron-node-linux-arm64.tar.gz`
- `gocron-node-darwin-amd64.tar.gz`
- `gocron-node-darwin-arm64.tar.gz`
- `gocron-node-windows-amd64.zip`

**Method 2: Download Using curl**

```bash
# Linux amd64
curl -L -o gocron-node-package/gocron-node-linux-amd64.tar.gz \
  https://github.com/gocronx-team/gocron/releases/latest/download/gocron-node-linux-amd64.tar.gz

# Linux arm64
curl -L -o gocron-node-package/gocron-node-linux-arm64.tar.gz \
  https://github.com/gocronx-team/gocron/releases/latest/download/gocron-node-linux-arm64.tar.gz

# macOS amd64
curl -L -o gocron-node-package/gocron-node-darwin-amd64.tar.gz \
  https://github.com/gocronx-team/gocron/releases/latest/download/gocron-node-darwin-amd64.tar.gz

# macOS arm64 (Apple Silicon)
curl -L -o gocron-node-package/gocron-node-darwin-arm64.tar.gz \
  https://github.com/gocronx-team/gocron/releases/latest/download/gocron-node-darwin-arm64.tar.gz

# Windows amd64
curl -L -o gocron-node-package/gocron-node-windows-amd64.zip \
  https://github.com/gocronx-team/gocron/releases/latest/download/gocron-node-windows-amd64.zip
```

### How It Works

The installation script automatically detects whether local packages are available:

- ✅ **Local package available** - Download directly from intranet server, fast and no external network required
- ⚠️ **Local package unavailable** - Automatically redirect to GitHub download (requires external network access)

Users will see clear prompts when executing installation commands:
```
Downloading from local server...
or
Downloading from GitHub...
```

::: tip Tip
- Filenames must strictly follow the format: `gocron-node-{os}-{arch}.{ext}`
- This configuration is optional; without it, packages will be downloaded from GitHub automatically
- It's recommended to update to the latest version regularly
:::

## System Service Management

### Linux (systemd)

Installation script automatically creates systemd service:

```bash
# Check service status
sudo systemctl status gocron-node

# Start service
sudo systemctl start gocron-node

# Stop service
sudo systemctl stop gocron-node

# Restart service
sudo systemctl restart gocron-node

# View logs
sudo journalctl -u gocron-node -f
```

### macOS (Background Process)

The macOS installation script runs gocron-node in the background using `nohup`, rather than creating a launchd service:

```bash
# Check process status
ps aux | grep gocron-node | grep -v grep

# Stop process
pkill -f gocron-node

# Start process
nohup /usr/local/gocron-node/gocron-node > /tmp/gocron-node.log 2>&1 &

# View logs
tail -f /tmp/gocron-node.log
```

::: tip Tip
If you need auto-start on boot, you can manually create a launchd configuration file. See [macOS launchd documentation](https://www.launchd.info/).
:::

### Windows (Windows Service)

For manual installation, Windows Service is not created by default. You can run `gocron-node.exe` directly, or use `sc` command to create a service manually:

```powershell
# Create service
sc.exe create gocron-node binPath= "D:\gocron-node\gocron-node.exe" start= auto

# Start service
Start-Service gocron-node
```

If the service is created, you can manage it using:

```powershell
# Check service status
Get-Service gocron-node

# Stop service
Stop-Service gocron-node

# Restart service
Restart-Service gocron-node
```

## Uninstall Node

### Linux / macOS

```bash
# Stop service/process
sudo systemctl stop gocron-node  # Linux
pkill -f gocron-node  # macOS

# Remove service file
sudo rm /etc/systemd/system/gocron-node.service  # Linux

# Remove installation directory
sudo rm -rf /opt/gocron-node  # Linux
sudo rm -rf /usr/local/gocron-node  # macOS

# Remove configuration file
rm -rf ~/.gocron-node

# Remove log file (macOS)
rm /tmp/gocron-node.log  # macOS
```

### Windows

```powershell
# Stop and remove service
Stop-Service gocron-node
sc.exe delete gocron-node

# Remove installation directory
Remove-Item -Recurse -Force "C:\Program Files\gocron-node"

# Remove configuration file
Remove-Item -Recurse -Force "$env:USERPROFILE\.gocron-node"
```

## Troubleshooting

### Installation Failed

**Problem**: Installation script execution failed

**Solution**:
1. Check network connection
2. Confirm token hasn't expired (3-hour validity)
3. Check for administrator privileges
4. View installation logs for detailed error information

### Node Not Showing

**Problem**: Installation complete but node not visible in list

**Solution**:
1. Check if node service is running properly
2. Check firewall settings
3. Confirm server address is configured correctly
4. View node logs:
   - Linux/macOS: `sudo journalctl -u gocron-node -f`
   - Windows: Event Viewer

### Token Expired

**Problem**: Token has expired and cannot be used

**Solution**:
Regenerate token. Click "Auto Register" in web interface to generate new token.

### Insufficient Permissions

**Problem**: Linux/macOS prompts insufficient permissions

**Solution**:
Use a regular user with `sudo` privileges to execute the installation command. **Note: Do not run the installation script directly as root user**. The script will automatically use `sudo` to elevate privileges for operations that require administrator permissions.

**Problem**: Windows prompts insufficient permissions

**Solution**:
Right-click PowerShell and select "Run as administrator".

## Related Documentation

- [Quick Start](./quick-start) - Learn about other deployment methods
- [Configuration](./configuration) - Node configuration instructions
- [Scheduled Tasks](./scheduled-tasks) - Create and manage tasks
