# P9 本地 Docker 环境

本目录用于维护 P9 本地 Docker 基础环境，当前主要用于数据库设计、本地开发、基础组件联调及后续部署配置参考。

后续可根据实际部署需要，将本目录配置迁移到对应业务或部署仓库，并按 Dev、Alpha、Beta、Prod 等环境调整实际配置。

## 1. 环境划分

| 环境 | Compose Project | 目录 |
|---|---|---|
| 总网 | `p9-platform` | `platform/` |
| 调度中心 | `p9-node-manager` | `node-manager/` |
| 分站节点 01 | `p9-operator-node-01` | `operator/node-01/` |
| 分站节点 02 | `p9-operator-node-02` | `operator/node-02/` |
| 分站节点 03 | `p9-operator-node-03` | `operator/node-03/` |

> `node-manager` 为调度中心当前暂定技术名称。

## 2. 端口

| 环境 | PostgreSQL | Redis | Etcd |
|---|---:|---:|---:|
| Platform | `15432` | `16379` | `12379` |
| Node Manager | `15440` | `16380` | - |
| Operator Node 01 | `15441` | `16381` | - |
| Operator Node 02 | `15442` | `16382` | - |
| Operator Node 03 | `15443` | `16383` | - |

容器内部端口统一：

| 服务 | 端口 |
|---|---:|
| PostgreSQL | `5432` |
| Redis | `6379` |
| Etcd | `2379` |

## 3. 数据库

| 环境 | 数据库 |
|---|---|
| Platform | `p9_platform` |
| Node Manager | `p9_node_manager` |
| Operator Node 01 | `p9_operator` |
| Operator Node 02 | `p9_operator` |
| Operator Node 03 | `p9_operator` |

PostgreSQL 本地账号：

```text
Username: root
Password: root
```

## 4. 当前镜像版本

```text
PostgreSQL: postgres:18.6
Redis:      redis:8.8.0
Etcd:       gcr.io/etcd-development/etcd:v3.6.14
```

当前先使用较新的稳定版本进行本地开发，正式版本后续讨论确认。

版本讨论参考：

```text
PostgreSQL: 18.x / 17.x
Redis:      8.8.x / 8.2.x
```

## 5. 常用命令

进入对应环境目录后执行：

```powershell
docker compose up -d
docker compose ps
docker compose logs -f
docker compose down
```

重新拉取镜像：

```powershell
docker compose pull
```

删除当前环境及数据卷：

```powershell
docker compose down -v
```

查看 P9 相关容器：

```powershell
docker ps --filter "name=p9-"
```

查看 P9 相关数据卷：

```powershell
docker volume ls --filter "name=p9-"
```