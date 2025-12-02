# AWS MCP Gateway

<p align="center">
  <img src="/assets/images/go-gohper.png" alt="Gophard, the Lambda MCP Wizard" width="300">
</p>

![Docker Version](https://img.shields.io/docker/v/bwalheim1205/aws-mcp-gateway?sort=semver)

```aws-mcp-gateway``` is a lightweight Go gateway that exposes AWS Lambda functions as MCP-compatible tools, inspired by the API Gateway Lambda proxy. It lets developers focus on defining their Lambda functions as tools, which can be discovered, listed, and invoked programmatically via MCP, with input/output mapping handled automatically. Allows developers to focus on the fun part - the tools!

# Features
- **MCP SSE Server**: Tools are served over SSE for real-time interaction.
- **Configuration Driven**: Define tool input schema, name, description, and target Lambda in config.
- **IAM Authentication**: System uses standard IAM Authentication process. Will use IAM role or environemnt variables to authenticate to lambda

### Prerequisites
*   Go 1.21 or higher
*   Git

# Installation

### Clone Repository

First clone down the repository then follow instructions for either Go executable or docker build

```sh
git clone https://github.com/bwalheim1205/aws-mcp-gateway.git
cd aws-mcp-gateway
```

### Go

```sh
go build -o aws-mcp-gateway ./cmd/mcp-server
```

### Docker

```sh
docker build . -t aws-mcp-gateway
```

# Configuratoin

The `aws-mcp-gateway` uses a YAML configuration file to define the Lambda functions you want to expose as MCP tools. Each tool includes metadata, the Lambda ARN, and an optional input schema to describe its parameters.

### **Example Configuration**

```yaml
name: LambdaMCPGateway
version: v1.0.0
port: 8080
endpoint: /mcp/sse

tools:
  - name: getWeather
    description: Gets the weather for a specific city
    lambdaArn: arn:aws:lambda:us-east-1:123456789012:function:getWeather
    inputSchema:
      type: object
      properties:
        city:
          type: string
          description: City for which to get the weather
      required: ["city"]
```

### **Fields**

* **name**: Name of the MCP server broadcasted. Defaults to LambdaMCPGateway
* **version**: Version of the MCP server it will broadcast. Defaults to v1.0.0
* **endpoint**: The request path to host mcp server at. Defaults to /mcp/sse
* **port**: The port MCP server is hosted at. Defaults to 8080
* **tools**: A list of Lambda functions to expose as MCP tools.

  * **name**: Unique name for tool as it will appear in MCP.
  * **description**: Short description of the tool’s functionality.
  * **lambdaArn**: The ARN of the Lambda function to invoke.
  * **inputSchema**: JSON Schema describing the tool’s input parameters.
    This allows MCP clients to validate input and provide better tooling for users.

# Usage

Once you've completed either build you can run the aws-mcp-gatewat using executable or docker image. The MCP server will then be available at http://localhost:8080/mcp/sse

### Go
```sh
./aws-mcp-gateway -f /path/to/tools.yaml
```

### Docker
```sh
docker run -v /path/to/tools.yaml:/app/tools.yaml -p 8080:8080 aws-mcp-gateway
```

# Roadmap

If there's something you'd like to see implemented we'd love to her from you just open an issue. Here are some of the functionality next on the horizon:

- **Lambda Auto Discovery**: Coming up with potential tagging approach
- **Assume Role**: Allow configuration to AssumeRole for cross-account Lambda access or temporary elevated permissions, making the gateway more flexible and secure for multi-account setups
