# cmd-store

A lightweight tool to organize and execute domain-based Bash commands.

## 📁 Config Location

All configurations are stored under: `~/.cmd-store`


## 🗂️ Folder Structure

All domain-specific command configurations are stored in the `~/.cmd-store` directory as separate `.yaml` files. Each file represents a domain and includes the commands relevant to that domain.

### Example
``` bash
├── ~/.cmd-store
│   ├── aws.yaml
│   ├── docker.yaml

```

- `aws.yaml`: Contains AWS CLI-related shortcuts.
- `docker.yaml`: Stores shortcuts for Docker commands.

Each of these files follows the configuration format described in the [Configuration Structure](#-configuration-structure) section above.

> You can add or remove domains by simply adding or deleting YAML files in this directory.


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


## 🧩 Shell Completion

`cmd-store` supports shell autocompletion to make your workflow smoother and faster.

You can generate a completion script using:

```bash
cmd-store completion
```
> Note: You’ll need to re-run this command any time you add a new domain or command to update the completions.