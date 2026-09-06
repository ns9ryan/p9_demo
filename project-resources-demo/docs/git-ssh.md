# P9 Git SSH 配置

## 1. Git 地址

```text
Git：https://oa.98ent.com
SSH：oassh.98ent.com:2023
组织：p9
```

## 2. hosts

```text
13.250.16.160 oa.98ent.com
13.250.16.160 oassh.98ent.com
```

## 3. SSH Key

建议为 P9 单独创建 SSH Key：

```bash
ssh-keygen -t ed25519 -C "<Git账号>@p9" -f /c/ssh/certs/id_ed25519_<Git账号>_p9
```

生成后，将公钥添加到 `oa.98ent.com` 的 SSH / GPG 密钥。

## 4. SSH 配置

```sshconfig
Host oassh.98ent.com
    HostName 13.250.16.160
    Port 2023
    User git
    ServerAliveInterval 30
    ServerAliveCountMax 3
    PreferredAuthentications publickey
    IdentityFile C:/ssh/certs/id_ed25519_<Git账号>_p9
    IdentitiesOnly yes
```

测试：

```bash
ssh -T oassh.98ent.com
```
