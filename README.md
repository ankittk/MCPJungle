<h1 align="center">
  :deciduous_tree: MCPJungle :deciduous_tree:
</h1>
<p align="center">
  Self-hosted MCP Gateway for your private AI agents
</p>
<p align="center">
  <a href="https://discord.gg/TSrUCTw9" style="text-decoration: none;">
    <img src="https://img.shields.io/badge/Discord-MCPJungle-5865F2?style=flat-square&logo=discord&logoColor=white" alt="Discord" style="max-width: 100%;">
  </a>
</p>

MCPJungle is a single source-of-truth registry for all [Model Context Protocol](https://modelcontextprotocol.io/introduction) Servers running in your Organisation.

🧑‍💻 Developers use it to register & manage MCP servers and the tools they provide from a central place.

🤖 MCP Clients use it to discover and consume all these tools from a single "Gateway" MCP Server.

![diagram](./assets/mcpjungle-diagram/mcpjungle-diagram.png)

<p align="center">MCPJungle is the only MCP Server your AI agents need to connect to!</p>

# Who should use MCPJungle?
1. **Developers** using MCP Clients like Claude & Cursor that need to access MCP servers for tool-calling
2. **Developers** building production-grade AI Agents that need to access MCP servers with built-in security, privacy and Access Control.
3. **Organisations** wanting to view & manage all MCP client-server interactions from a central place. Hosted in their own datacenter 🔒

# 📋 Table of Contents

- [Quick Start guide](#quickstart-guide)
- [Installation](#installation)
- [Usage](#usage)
  - [Server](#server)
    - [Running mcpjungle server inside Docker](#running-inside-docker)
    - [Running mcpjungle server directly on the host machine](#running-directly)
  - [Client](#client)
    - [Adding Streamable HTTP-based MCP servers](#registering-streamable-http-based-servers)
    - [Adding STDIO-based MCP servers](#registering-stdio-based-servers)
    - [Removing MCP servers](#deregistering-mcp-servers)
  - [Connect to mcpjungle from Claude](#claude)
  - [Connect to mcpjungle from Cursor](#cursor)
  - [Enabling/Disabling Tools globally](#enablingdisabling-tools)
  - [Authentication](#authentication)
  - [Enterprise features](#enterprise-features-)
    - [Access Control](#access-control)
- [Limitations](#current-limitations-)
- [Contributing](#contributing-)

# Quickstart guide
This quickstart guide will show you how to:
1. Start the MCPJungle server locally using docker-compose
2. Register a simple MCP server in mcpjungle
3. Connect your Claude to mcpjungle to access your MCP tools

## Start the server
```bash
curl -O https://raw.githubusercontent.com/mcpjungle/MCPJungle/refs/heads/main/docker-compose.yaml
docker-compose up -d
```

## Register MCP servers
Download the client binary either using brew or from the [Releases](https://github.com/mcpjungle/MCPJungle/releases).
```bash
brew install mcpjungle/mcpjungle/mcpjungle
```

Add the [context7](https://context7.com/) remote MCP server to mcpjungle
```bash
mcpjungle register --name context7 --url https://mcp.context7.com/mcp
```

## Connect to mcpjungle

Add the following configuration in your Claude MCP Servers:
```json
{
  "mcpServers": {
    "mcpjungle": {
      "command": "npx",
      "args": [
        "mcp-remote",
        "http://localhost:8080/mcp",
        "--allow-http"
      ]
    }
  }
}
```

Try asking Claude for the following:
```text
Use context7 to get the documentation for `/lodash/lodash`
```

Claude will then attempt to call the `context7__get-library-docs` tool via MCPJungle, which will return the documentation for the Lodash library.

<p align="center">
  <img src="./assets/quickstart-claude-call-tool.png" alt="claude calls context7 tool via mcpjungle" height="400">
</p>

Congratulations 🎉 You have successfully registered a remote MCP server in MCPJungle and called one of its tools via Claude

# Installation

> [!WARNING]
> MCPJungle is **BETA** software.
>
> We're actively working to make it production-ready.
> You can provide your feedback by [starting a discussion](https://github.com/mcpjungle/MCPJungle/discussions) in this repository.

MCPJungle is shipped as a stand-alone binary.

You can either download it from the [Releases](https://github.com/mcpjungle/MCPJungle/releases) Page or use [Homebrew](https://brew.sh/) to install it:

```bash
brew install mcpjungle/mcpjungle/mcpjungle
```

Verify your installation by running

```bash
mcpjungle version
```

> [!IMPORTANT]
> On MacOS, you will have to use homebrew because the compiled binary is not [Notarized](https://developer.apple.com/documentation/security/notarizing-macos-software-before-distribution) yet.

MCPJungle provides a Docker image which is useful for running the registry server (more about it later).

```bash
docker pull mcpjungle/mcpjungle
```

# Usage
MCPJungle has a Client-Server architecture and the binary lets you run both the Server and the Client.

## Server
The MCPJungle server is responsible for managing all the MCP servers registered in it and providing a unified MCP gateway for AI Agents to discover and call tools provided by these registered servers.

The gateway itself runs over streamable http transport and is accessible at the `/mcp` endpoint.

### Running inside Docker
For running the MCPJungle server locally, docker compose is the recommended way:
```shell
curl -O https://raw.githubusercontent.com/mcpjungle/MCPJungle/refs/heads/main/docker-compose.yaml
docker-compose up -d
```

This will start the MCPJungle server along with a persistent Postgres database container.

You can quickly verify that the server is running:
```bash
curl http://localhost:8080/health
```

If you plan on registering stdio-based MCP servers that rely on `npx` or `uvx`, use mcpjungle's `stdio` tagged docker image instead.
```bash
MCPJUNGLE_IMAGE_TAG=latest-stdio docker-compose up -d
```

This image is significantly larger. But it is very convenient and recommended for running locally when you rely on stdio-based MCP servers.

For example, if you only want to register remote mcp servers like context7 and deepwiki, you can use the standard (minimal) image.

But if you also want to use stdio-based servers like `filesystem`, `time`, `github`, etc., you should use the `stdio`-tagged image instead.

> [!NOTE]
> If your stdio servers rely on tools other than `npx` or `uvx`, you will have to create a custom docker image that includes those dependencies along with the mcpjungle binary.

**Production Deployment**

The default [MCPJungle Docker image](https://hub.docker.com/r/mcpjungle/mcpjungle) is very lightweight - it only contains a minimal base image and the `mcpjungle` binary.

It is therefore suitable and recommended for production deployments.

For the database, we recommend you deploy a separate Postgres DB cluster and supply its endpoint to mcpjungle (see [Database](#database) section below).

You can see the definitions of the [standard Docker image](./Dockerfile) and the [stdio Docker image](./stdio.Dockerfile).

### Running directly
You can also run the server directly on your host machine using the binary:

```bash
mcpjungle start
```

This starts the main registry server and MCP gateway, accessible on port `8080` by default.


### Database
The mcpjungle server relies on a database and by default, creates a SQLite DB in the current working directory.

This is okay when you're just testing things out locally.

Alternatively, you can supply a DSN for a Postgresql database to the server:

```bash
export DATABASE_URL=postgres://admin:root@localhost:5432/mcpjungle_db

#run as container
docker run mcpjungle/mcpjungle:latest

# or run directly
mcpjungle start
```

## Client
Once the server is up, you can use the mcpjungle CLI to interact with it.

MCPJungle currently supports MCP servers using [stdio](https://modelcontextprotocol.io/specification/2025-03-26/basic/transports#stdio) and [Streamable HTTP](https://modelcontextprotocol.io/specification/2025-03-26/basic/transports#streamable-http) Transports.

Let's see how to register them in mcpjungle.

### Registering streamable HTTP-based servers
Let's say you're already running a streamable http MCP server locally at `http://127.0.0.1:8000/mcp` which provides basic math tools like `add`, `subtract`, etc.

You can register this MCP server with MCPJungle:
```bash
mcpjungle register --name calculator --description "Provides some basic math tools" --url http://127.0.0.1:8000/mcp
```

If you used docker-compose to run the server, and you're not on Linux, you will have to use `host.docker.internal` instead of your local loopback address.
```bash
mcpjungle register --name calculator --description "Provides some basic math tools" --url http://host.docker.internal:8000/mcp
```

The registry will now start tracking this MCP server and load its tools.

![register a MCP server in MCPJungle](./assets/register-mcp-server.png)

You can also provide a configuration file to register the MCP server:
```bash
cat ./calculator.json
{
  "name": "calculator",
  "transport": "streamable_http",
  "description": "Provides some basic math tools",
  "url": "http://127.0.0.1:8000/mcp"
}

mcpjungle register -c ./calculator.json
```

All tools provided by this server are now accessible via MCPJungle:

```bash
mcpjungle list tools

# Check tool usage
mcpjungle usage calculator__multiply

# Call a tool
mcpjungle invoke calculator__multiply --input '{"a": 100, "b": 50}'
```

![Call a tool via MCPJungle Proxy MCP server](./assets/tool-call.png)

> [!NOTE]
> A tool in MCPJungle must be referred to by its canonical name which follows the pattern `<mcp-server-name>__<tool-name>`.
> Server name and tool name are separated by a double underscore `__`.
>
> eg- If you register a MCP server `github` which provides a tool called `git_commit`, you can invoke it in MCPJungle using the name `github__git_commit`.
> 
> Your MCP client must also use this canonical name to call the tool via MCPJungle.

The config file format for registering a Streamable HTTP-based MCP server is:
```json
{
  "name": "<name of your mcp server>",
  "transport": "streamable_http",
  "description": "<description>",
  "url": "<url of the mcp server>",
  "bearer_token": "<optional bearer token for authentication>"
}
```

### Registering STDIO-based servers

Here's an example configuration file (let's call it `filesystem.json`) for a MCP server that uses the STDIO transport:
```json
{
  "name": "filesystem",
  "transport": "stdio",
  "description": "filesystem mcp server",
  "command": "npx",
  "args": ["-y", "@modelcontextprotocol/server-filesystem", "."]
}
```

You can register this MCP server in MCPJungle by providing the configuration file:
```bash
mcpjungle register -c ./filesystem.json
```

The config file format for registering a STDIO-based MCP server is:

```json
{
  "name": "<name of your mcp server>",
  "transport": "stdio",
  "description": "<description>",
  "command": "<command to run the mcp server, eg- 'npx', 'uvx'>",
  "args": ["arguments", "to", "pass", "to", "the", "command"],
  "env": {
    "KEY": "value"
  }
}
```

**Limitation** 🚧

MCPJungle creates a new connection when a tool is called. This means a new sub-process for a STDIO mcp server is started for every tool call.

This has some performance overhead but ensures that there are no memory leaks.

But it also means that currently MCPJungle doesn't support stateful connections with your MCP server.

We want to hear your feedback to improve this mechanism, feel free to create an issue, start a discussion or just reach out on Discord.


### Deregistering MCP servers
You can remove a MCP server from mcpjungle.

```bash
mcpjungle deregister calculator
mcpjungle deregister filesystem
```

Once removed, this mcp server and its tools are no longer available to you or your MCP clients.

## Integration with other MCP Clients
Assuming that MCPJungle is running on `http://localhost:8080`, use the following configurations to connect to it:

### Claude
```json
{
  "mcpServers": {
    "mcpjungle": {
      "command": "npx",
      "args": [
        "mcp-remote",
        "http://localhost:8080/mcp",
        "--allow-http"
      ]
    }
  }
}
```

### Cursor
```json
{
  "mcpServers": {
    "mcpjungle": {
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

## Enabling/Disabling Tools
You can enable or disable a specific tool or all the tools provided by an MCP Server.

If a tool is disabled, it is not available via the MCPJungle Proxy, so no MCP clients can view or call it.

```bash
# disable the `get-library-docs` tool provided by the `context7` MCP server
mcpjungle disable context7__get-library-docs

# re-enable the tool
mcpjungle enable context7__get-library-docs

# disable all tools provided by the `context7` MCP server
mcpjungle disable context7

# re-enable all tools of `context7`
mcpjungle enable context7
```

A disabled tool is still accessible via mcpjungle's HTTP API, so humans can still manage it from the CLI (or any other HTTP client).

> [!NOTE]
> When a new server is registered in MCPJungle, all its tools are **enabled** by default.

## Authentication
MCPJungle currently supports authentication if your Streamable HTTP MCP Server accepts static tokens for auth.

This is useful when using SaaS-provided MCP Servers like HuggingFace, Stripe, etc. which require your API token for authentication.

You can supply your token while registering the MCP server:
```bash
# If you specify the `--bearer-token` flag, MCPJungle will add the `Authorization: Bearer <token>` header to all requests made to this MCP server.
mcpjungle register --name huggingface --description "HuggingFace MCP Server" --url https://huggingface.co/mcp --bearer-token <your-hf-api-token>
```

Or from your configuration file
```bash
{
  "name": "huggingface",
  "transport": "streamable_http",
  "url": "https://huggingface.co/mcp",
  "description": "hugging face mcp server",
  "bearer_token": "<your-hf-api-token>"
}
```

Support for Oauth flow is coming soon!

## Enterprise Features 🔒

If you're running MCPJungle in your organisation, we recommend running the Server in the `production` mode:
```bash
# enable enterprise features by running in production mode
mcpjungle start --prod

# you can also specify the server mode as environment variable (valid values are `development` and `production`)
export SERVER_MODE=production
mcpjungle start

# this also works when running the server via docker-compose
SERVER_MODE=production docker-compose up
```

By default, mcpjungle server runs in `development` mode which is ideal for individuals running it locally.

In Production mode, the server enforces stricter security policies and will provide additional features like Authentication, ACLs, observability and more.

After starting the server in production mode, you must initialize it by running the following command on your client machine:
```bash
mcpjungle init-server
```

This will create an admin user in the server and store its API access token in your home directory (`~/.mcpjungle.conf`).

You can then use the mcpjungle cli to make authenticated requests to the server.

### Access Control

In `development` mode, all MCP clients have full access to all the MCP servers registered in MCPJungle Proxy.

`production` mode lets you control which MCP clients can access which MCP servers.

Suppose you have registered 2 MCP servers `calculator` and `github` in MCPJungle in production mode.

By default, no MCP client can access these servers. **You must create an MCP Client in mcpjungle and explicitly allow it to access the MCP servers.**

```bash
# Create a new MCP client for your Cursor IDE to use. It can access the calculator and github MCP servers
mcpjungle create mcp-client cursor-local --allow "calculator, github"

MCP client 'cursor-local' created successfully!
Servers accessible: calculator,github

Access token: 1YHf2LwE1LXtp5lW_vM-gmdYHlPHdqwnILitBhXE4Aw
Send this token in the `Authorization: Bearer {token}` HTTP header.
```

Mcpjungle creates an access token for your client.
Configure your client or agent to send this token in the `Authorization` header when making requests to the mcpjungle proxy.

For example, you can add the following configuration in Cursor to connect to MCPJungle:

```json
{
  "mcpServers": {
    "mcpjungle": {
      "url": "http://localhost:8080/mcp",
      "headers": {
        "Authorization": "Bearer 1YHf2LwE1LXtp5lW_vM-gmdYHlPHdqwnILitBhXE4Aw"
      }
    }
  }
}
```

A client that has access to a particular server this way can view and call all the tools provided by that server.

> [!NOTE]
> If you don't specify the `--allow` flag, the MCP client will not be able to access any MCP servers.

# Current limitations 🚧
We're not perfect yet, but we're working hard to get there!

### 1. MCPJungle doesn't maintain any long-running connections to the registered MCP Servers
When you call a tool in a Streamable HTTP server, mcpjungle creates a new connection to the server to serve the request.

When you call a tool in a STDIO server, mcpjungle creates a new connection and starts a new sub-process to run this server.

After servicing your request, it terminates this sub-process.

So a new stdio server process is started for every tool call.

This has some performance overhead but ensures that there are no memory leaks.

It also means that if you rely on stateful connections with your MCP server, mcpjungle can currently not provide that.

We plan on improving this mechanism in future releases and are open to ideas from the community!

### 2. MCPJungle does not support OAuth flow for authentication.
This is a work in progress.

We're collecting more feedback on how people use OAuth with MCP servers, so feel free to start a Discussion or open an issue to share your use case.

# Contributing 💻

If you're interested in contributing to MCPJungle, see [Developer Docs](./docs/developer.md).
