# DaoKin

这是 DaoKin 的开发仓库，远端为 [JeanphiloGong/daokin](https://github.com/JeanphiloGong/daokin)。当前完成顶层目录、Go 模块和 V1 身份分层目录初始化，业务功能仍需按教程逐步实现。

## 当前只做什么

核心故事是：用户证明钱包身份后，直接创建并发布自己的作品。免费作品可以公开阅读；选择持仓阅读的版本通过钱包发行作品资产，读者购买或接收份额后，按当前可信持仓取得后续阅读资格。

作品直接关联作者，不要求先创建或加入 Dao。社区另行设计。作品资产不自动转移版权或作者管理权；供应、余额、持仓与交割以链上合约为准，Currency 负责核验与索引。

本轮 MVP 使用测试币 TST；商城、Console、旧积分充值、旧论坛迁移、USDT 兑换与任意作者合约均不在当前主线。

## 当前目录

```text
daokin/                              # 一个 Git 仓库
├── README.md
├── .gitignore
├── api/                             # 身份、作品、业务意图与正文交付
│   ├── go.mod
│   └── internal/
│       └── identity/                # V1 从身份模块开始
│           ├── domain/              # 用户、挑战与会话规则
│           ├── application/         # 用例与所需接口
│           ├── infrastructure/      # 存储与验签实现
│           └── interfaces/http/     # 请求、响应与身份入口
├── currency/
│   └── go.mod                       # 独立链上核验、索引服务的模块
├── contracts/
│   └── go.mod                       # Proto 与生成代码的模块
├── chain/                           # V3 添加 Solidity 与 Foundry
├── web/                             # V1.3 添加 Next.js + React + TypeScript 页面
└── scripts/                         # 随章节添加生成、部署与验收程序
```

空目录用 `.gitkeep` 保留在 Git 中，开始添加实际文件时可删除相应占位文件。API 的启动命令、迁移和具体基础设施在 V1.2 添加，work 模块在 V2 添加；Currency 内部分层与 Proto 在 V3 添加。

教程独立维护于 [daokin-tutorial](https://github.com/JeanphiloGong/daokin-tutorial)，本地目录是同级 `../daokin-tutorial/`。架构说明见其 [架构图集](../daokin-tutorial/architecture/README.md)。本仓库按章节增加实际源码。

## 模块与运行边界

| 目录 | 是否运行服务 | 通信与数据职责 |
| --- | --- | --- |
| `web/` | Next.js Web 服务 | 同源页面、Cookie、HTTP/JSON 调用 API；浏览器请求钱包签名 |
| `api/` | Go HTTP API | 拥有身份与作品数据库；通过 gRPC 查询 Currency |
| `currency/` | Go gRPC 服务，内含索引循环 | 拥有独立索引数据库；读取、核验链上事实 |
| `contracts/` | 否 | API 与 Currency 共用的 Proto、生成类型和 RPC 定义 |
| `chain/` | 合约由 EVM 执行 | TST、作品份额、市场交割的权威状态 |

三个 Go 模块都使用 Go 1.25.5，模块名为 `github.com/JeanphiloGong/daokin/{api,currency,contracts}`。根目录不是 Go 模块，各模块内运行 Go 命令；V3 使用协议时按教程增加 `require` 和 `replace ... => ../contracts`。当前不需要根 `go.mod` 或 `go.work`。

API、Currency 的 domain/application 只包含业务规则与用例所需接口。SQL、HTTP/gRPC 实现、生成 Proto 与 Ethereum SDK 留在具体边界；接口与模型在出现实际调用者时添加。钱包负责签署和广播，服务不持有用户私钥。

## 从这里开始

先完成 [V1.1：稳定身份、单次挑战与会话期限](../daokin-tutorial/v1-01-stable-identity.md)，将模型放入 `api/internal/identity/domain/`，再实现 V1.2 的存储、验签与 HTTP。本地未跟踪的 `api/internal/user.go` 是未完成的手写草稿，不属于本次初始化提交；迁入正确位置并完成实现后，应删除旧位置的草稿。

初始化阶段，在项目根检查已有模块声明：

```bash
for module in api currency contracts; do
    (cd "$module" && go list -m)
done
```

该检查只证明模块声明有效。提交的骨架尚无业务 Go 包，没有可运行的服务或业务测试；若本地保留上述草稿，其未使用的 imports 会使 `cd api && go test ./...` 失败。

| 阶段 | 增加的实际内容 |
| --- | --- |
| V1.1 | 身份、挑战、会话领域模型及检查 |
| V1.2–V1.3 | 登录事务、HTTP、Web 登录与退出 |
| V2 | work 模块、不可变发布版本、创作与阅读页面 |
| V3 | 三合约、Proto、Currency、发行绑定与持仓阅读 |
| V4–V6 | 购买、转让、本人工作台与完整验收 |

每完成一个检查点，在本仓库执行对应检查再继续。教程隔离副本的测试结果不代表本仓库已实现，公开测试网与真实钱包扩展也需单独验收。

## 验证

- 初始化检查：三个模块声明、目录与教程路径一致，以及 `git diff --check`。
- V1.1 起执行领域检查；V1.2 起执行真实 PostgreSQL 的 API 检查；后续按各章加入 Web、合约和链上场景。
- 当前 `main` 从新项目初始化开始实现；旧版原型保存在 [legacy/prototype-2026-11-13](https://github.com/JeanphiloGong/daokin/tree/legacy/prototype-2026-11-13)。
