# cmd-store

A lightweight tool to organize and execute domain-based Bash commands.

## 📁 Config Location

All configurations are stored under: `~/.cmd-store`




## 🛠️ Configuration Structure

Each configuration file represents a domain and follows this YAML structure:

```yaml
description: AWS CLI commands
commands:
  - name: pc
    cmd: echo "This is a placeholder for AWS commands"
    description: "Placeholder command for AWS operations"
```

description: A brief summary of the domain's purpose.

commands: A list of command entries, each including:

name: The shortcut used to trigger the command.

cmd: The actual shell command to be executed.

description: A friendly explanation of what the command does.

🚀 Usage

`cmd-store <domain> <command-name>`

For example: `cmd-store aws pc`


## 📦 Downloading the Latest Release

To get the latest version of `cmd-store`, head over to the [Releases page](https://github.com/ezratameno/cmd-store/releases/latest) on GitHub.

1. Click on the latest release tag.
2. Under **Assets**, download the appropriate file for your system (e.g., `.zip`, `.tar.gz`, or binary).
3. Extract or install it according to your platform's requirements.

Alternatively, you can use `curl` or `wget` to download directly from the command line if a binary is available:

```bash
wget https://github.com/ezratameno/cmd-store/releases/latest/download/cmd_store_linux_amd64.tar.gz

tar -xzf  cmd_store_linux_amd64.tar.gz
chmod +x cmd-store
sudo mv ./cmd-store /usr/local/bin/cmd-store
```