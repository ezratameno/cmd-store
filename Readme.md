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
