## 检测能力



OpenSCA 现在能够解析列出的编程语言和相应的包管理器中的配置文件。该团队现在致力于引入更多语言并逐步丰富相关配置文件的解析。

| 语言         | 包管理器   | 文件                                                         |
| ------------ | ---------- | ------------------------------------------------------------ |
| `Java`       | `Maven`    | `pom.xml`                                                    |
| `Java`       | `Gradle`   | `.gradle` `.gradle.kts`                                      |
| `JavaScript` | `Npm`      | `package-lock.json` `package.json` `yarn.lock`               |
| `PHP`        | `Composer` | `composer.json` `composer.lock`                              |
| `Ruby`       | `gem`      | `gemfile.lock`                                               |
| `Golang`     | `gomod`    | `go.mod` `go.sum` `Gopkg.toml` `Gopkg.lock`                  |
| `Rust`       | `cargo`    | `Cargo.lock`                                                 |
| `Erlang`     | `Rebar`    | `rebar.lock`                                                 |
| `Python`     | `Pip`      | `Pipfile` `Pipfile.lock` `setup.py` `requirements.txt` `requirements.in`（对于后两者，需要pipenv环境和互联网连接） |

## 下载和部署



1. 根据您的系统架构从[发行版](https://github.com/XmirrorSecurity/OpenSCA-cli/releases)下载相应的可执行文件。

2. 或者下载源码并编译（需要 及以上版本）`go 1.18`

   ```
   // github linux/mac
   git clone https://github.com/XmirrorSecurity/OpenSCA-cli.git opensca && cd opensca && go build
   // gitee linux/mac
   git clone https://gitee.com/XmirrorSecurity/OpenSCA-cli.git opensca && cd opensca && go build
   // github windows
   git clone https://github.com/XmirrorSecurity/OpenSCA-cli.git opensca ; cd opensca ; go build
   // gitee windows
   git clone https://gitee.com/XmirrorSecurity/OpenSCA-cli.git opensca ; cd opensca ; go build
   ```

   

   默认选项是生成当前系统架构的程序。如果您想尝试其他系统架构，可以在编译前设置以下环境变量。

   - 禁用`CGO_ENABLED` `CGO_ENABLED=0`
   - 设置操作系统`GOOS=${OS} \\ darwin,liunx,windows`
   - 设置架构`GOARCH=${arch} \\ amd64,arm64`

## 使用 OpenSCA

eg:config.json

### 参数



| 参数     | 类型     | 描述                                                         | 样本                      |
| -------- | -------- | ------------------------------------------------------------ | ------------------------- |
| `config` | `string` | 设置配置文件的路径                                           | `-config config.json`     |
| `path`   | `string` | 设置目标文件或目录的路径                                     | `-path ./foo`             |
| `out`    | `string` | 将结果保存到其格式由后缀定义的指定文件中                     | `-out out.json, out.html` |
| `log`    | `string` | 指定日志文件的路径                                           | `-log my_log.txt`         |
| `token`  | `string` | 来自我们官方网站的云服务验证                                 | `-token xxx`              |
| `proj`   | `string` | SaaS 项目与 [OpenSCA SaaS](https://opensca.xmirror.cn/console) 同步报告`token` | `-proj xxx`               |

从 v3.0.0 开始，除了这些可用于 CMD/CRT 的参数外，还有其他针对不同需求的参数必须在配置文件中设置。

有关每个参数的完整介绍，请参见`config.json`

v3.0.2 及更高版本允许 OpenSCA-cli 使用 OpenSCA SaaS（通过 ）进行 snyc 报告，以便可以一起管理多个项目的所有报告。`proj`

如果配置参数与命令行输入参数冲突，则将采用后者。

当设置的路径中没有配置文件时，将在那里生成一个默认设置的配置文件。

如果未设置配置文件的路径，则检查以下路径：

1. `config.json`在工作目录下
2. `opensca_config.json`在用户目录下
3. `config.json`在目录下`opensca-cli`

从 v3.0.0 开始，已放入配置文件中。默认设置将进入我们的云漏洞数据库。其他符合我们数据库结构的在线数据库也可以通过配置文件进行设置。`url`

使用以前的版本连接云数据库仍然需要 的设置，这可以通过 CMD 和配置文件来完成。例：`url``-url https://opensca.xmirror.cn`

### 报告格式



该参数支持的文件如下：`out`

| 类型      | 格式     | 指定后缀                         | 版本            |
| --------- | -------- | -------------------------------- | --------------- |
| 报告      | `json`   | `.json`                          | `*`             |
|           | `xml`    | `.xml`                           | `*`             |
|           | `html`   | `.html`                          | `v1.0.6`及以上  |
|           | `sqlite` | `.sqlite`                        | `v1.0.13`及以上 |
|           | `csv`    | `.csv`                           | `v1.0.13`及以上 |
|           | `sarif`  | `.sarif`                         |                 |
| SBOM 系列 | `spdx`   | `.spdx` `.spdx.json` `.spdx.xml` | `v1.0.8`及以上  |
|           | `cdx`    | `.cdx.json` `.cdx.xml`           | `v1.0.11`及以上 |
|           | `swid`   | `.swid.json` `.swid.xml`         | `v1.0.11`及以上 |
|           | `dsdx`   | `.dsdx` `.dsdx.json` `.dsdx.xml` | `v3.0.0`及以上  |

### 样本



```
# Use opensca-cli to scan with CMD parameters:
opensca-cli -path ${project_path} -config ${config_path} -out ${filename}.${suffix} -token ${token}

# Start scanning after setting down the configuration file:
opensca-cli
```

- `-path ${project_path}`：指定要扫描的项目路径。`${project_path}` 是一个变量，你需要将其替换为实际的项目路径.
- `-config ${config_path}`：指定配置文件的路径。`${config_path}` 是一个变量，你需要将其替换为实际的配置文件路径。配置文件通常包含扫描的相关设置，如要扫描的依赖项类型等.
- `-out ${filename}.${suffix}`：指定扫描结果的输出文件名和文件扩展名。`${filename}` 和 `${suffix}` 是变量，你需要将它们替换为实际的文件名和扩展名，例如 `-out scan_result.json` 表示将扫描结果输出到名为 `scan_result.json` 的文件中.
- `-token ${token}`：指定用于身份验证的令牌。`${token}` 是一个变量，你需要将其替换为实际的令牌值。令牌通常用于与 OpenSCA 服务器进行身份验证，以便上传扫描结果或获取其他相关服务.

#### 通过Docker容器扫描和报告



```
# 检测当前目录中的依赖项:
docker run -ti --rm -v ${PWD}:/src opensca/opensca-cli

# 连接到云漏洞数据库:
docker run -ti --rm -v ${PWD}:/src opensca/opensca-cli -token ${put_your_token_here}
```

- `-ti`：分配一个伪终端并保持标准输入打开，以便你可以与容器进行交互。
- `--rm`：在容器退出时自动删除容器，避免留下无用的容器实例。
- `-v ${PWD}:/src`：将当前工作目录（`${PWD}`）挂载到容器的 `/src` 目录中。这样，容器可以访问当前目录中的文件。
- `opensca/opensca-cli`：指定要运行的 Docker 镜像。
- 这个命令会启动 `opensca/opensca-cli` 容器，并在当前目录中检测依赖项。
- `-token ${put_your_token_here}`：指定用于连接到云漏洞数据库的 API 令牌。你需要将 `${put_your_token_here}` 替换为你的实际 API 令牌。
- 这个命令会启动 `opensca/opensca-cli` 容器，并使用指定的 API 令牌连接到云漏洞数据库，以便进行漏洞扫描和分析

您还可以使用配置文件进行高级设置。保存到容器的挂载目录或通过 设置容器内的其他路径。在不同终端挂载当前目录的编写方法各不相同，我们在这里列出常见的供参考：`config.json``src``-config`

| 终端         | pwd                   |
| ------------ | --------------------- |
| `bash`       | `$(pwd)`              |
| `zsh`        | `${PWD}`              |
| `cmd`        | `%cd%`                |
| `powershell` | `(Get-Location).Path` |

### 本地漏洞数据库



#### 漏洞数据库文件的格式

eg:db-demo.json

```
[
  {
    "vendor": "org.apache.logging.log4j",
    "product": "log4j-core",
    "version": "[2.0-beta9,2.12.2)||[2.13.0,2.15.0)",
    "language": "java",
    "name": "Apache Log4j2 远程代码执行漏洞",
    "id": "XMIRROR-2021-44228",
    "cve_id": "CVE-2021-44228",
    "cnnvd_id": "CNNVD-202112-799",
    "cnvd_id": "CNVD-2021-95914",
    "cwe_id": "CWE-502,CWE-400,CWE-20",
    "description": "Apache Log4j是美国阿帕奇（Apache）基金会的一款基于Java的开源日志记录工具。\r\nApache Log4J 存在代码问题漏洞，攻击者可设计一个数据请求发送给使用 Apache Log4j工具的服务器，当该请求被打印成日志时就会触发远程代码执行。",
    "description_en": "Apache Log4j2 2.0-beta9 through 2.12.1 and 2.13.0 through 2.15.0 JNDI features used in configuration, log messages, and parameters do not protect against attacker controlled LDAP and other JNDI related endpoints. An attacker who can control log messages or log message parameters can execute arbitrary code loaded from LDAP servers when message lookup substitution is enabled. From log4j 2.15.0, this behavior has been disabled by default. From version 2.16.0, this functionality has been completely removed. Note that this vulnerability is specific to log4j-core and does not affect log4net, log4cxx, or other Apache Logging Services projects.",
    "suggestion": "2.12.1及以下版本可以更新到2.12.2，其他建议更新至2.15.0或更高版本，漏洞详情可参考：https://github.com/apache/logging-log4j2/pull/608 \r\n1、临时解决方案，适用于2.10及以上版本：\r\n\t（1）设置jvm参数：“-Dlog4j2.formatMsgNoLookups=true”；\r\n\t（2）设置参数：“log4j2.formatMsgNoLookups=True”；",
    "attack_type": "远程",
    "release_date": "2021-12-10",
    "security_level_id": 1,
    "exploit_level_id": 1
  }
]
```



#### 漏洞数据库字段说明



| 田                  | 描述                            | 是否需要 |
| ------------------- | ------------------------------- | -------- |
| `vendor`            | 组件的制造商                    | N        |
| `product`           | 组件的名称                      | Y        |
| `version`           | 受漏洞影响的组件版本            | Y        |
| `language`          | 组件的编程语言                  | Y        |
| `name`              | 漏洞的名称                      | N        |
| `id`                | 自定义标识符                    | Y        |
| `cve_id`            | CVE 标识符                      | N        |
| `cnnvd_id`          | cnnvd 标识符                    | N        |
| `cnvd_id`           | CNVD 标识符                     | N        |
| `cwe_id`            | CWE 标识符                      | N        |
| `description`       | 漏洞描述                        | N        |
| `description_en`    | 英文对脆弱性的描述              | N        |
| `suggestion`        | 修复漏洞的建议                  | N        |
| `attack_type`       | 攻击类型                        | N        |
| `release_date`      | 漏洞的发布日期                  | N        |
| `security_level_id` | 漏洞的安全级别（从 1 降低到 4） | N        |
| `exploit_level_id`  | 漏洞的利用级别（0-N/A，1-可用） | N        |

*“language” 字段有几个预设值，包括 java、javascript、golang、rust、php、ruby 和 python。其他语言不限于预设值。

#### 设置漏洞数据库的示例



```
{
  "origin":{
    "json":"db.json",
    "mysql":{
      "dsn":"user:password@tcp(ip:port)/dbname",
      "table":"table_name"
    },
    "sqlite":{
      "dsn":"sqlite.db",
      "table":"table_name"
    }
  }
}
```

## 常见问题



### 使用 OpenSCA 时是否需要环境变量？



不。OpenSCA 解压后可以通过 CLI/CRT 中的命令直接执行。

### 关于漏洞数据库？



OpenSCA 允许配置本地漏洞数据库。它必须根据*漏洞数据库文件的格式*进行排序。

同时，OpenSCA 还提供了一个云漏洞数据库，涵盖包括 CVE/CWE/NVD/CNVD/CNNVD 在内的官方数据库。

### 关于 OpenSCA 扫描的时间成本？



这取决于包的大小、网络状况和语言。

从 v1.0.11 开始，我们将 aliyun 镜像数据库作为备份添加到官方 maven 仓库中，以解决网络连接带来的卡顿问题。

对于 v1.0.10 及以下版本，如果时间异常长，且日志文件中报告了 maven 仓库连接失败的错误信息，v1.0.6 到 v1.0.10 之间的用户可以设置如下字段来修复该问题：`maven``config.json`

```
{
    "maven": [
        {
            "repo": "https://maven.aliyun.com/repository/public",
            "user": "",
            "password": ""
        }
    ]
}   
```



设置完成后，保存到 opensca-cli.exe 的同一文件夹中并执行命令。例如：`config.json`

```
opensca-cli -token {token} -path {path} -out output.html -config config.json
```

- `-token {token}`：指定用于身份验证的令牌。`{token}` 是一个占位符，你需要将其替换为实际的认证令牌值。认证令牌是平台的临时密钥，相当于平台的账号和密码，主要用于作为本地检测引擎工具（OpenSCA-cli）访问云平台知识库的认证 .
- `-path {path}`：指定要扫描的项目路径。`{path}` 是一个占位符，你需要将其替换为实际的项目路径 .
- `-out output.html`：指定扫描结果的输出文件名和文件扩展名。`output.html` 表示将扫描结果输出到名为 `output.html` 的文件中，该文件将以 HTML 格式呈现 .
- `-config config.json`：指定配置文件的路径。`config.json` 是配置文件的文件名，表示配置文件位于当前目录下。配置文件通常包含扫描的相关设置，如要扫描的依赖项类型等 