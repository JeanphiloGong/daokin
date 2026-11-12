# DaoKin

这是从空项目开始的新产品开发入口。当前已建立本地Git仓库、API的Go模块，以及业务设计与教程规划；当前业务代码尚未实现，没有可运行的业务接口。代码远端为 [JeanphiloGong/daokin](https://github.com/JeanphiloGong/daokin)。

## 当前只做什么

核心故事是：用户创建社区和作品，作者选择作品使用规则；需要资产的作品再发行作品份额，读者购买或接收份额后，按规则取得使用资格。

先从稳定用户身份、社区与作品模型开始，后续逐步实现文章发布、有限阅读模板、资产供应与持仓、测试货币购买和转让。作品和资产分别建模，普通免费作品无需发行资产。版权和作者管理权的转移需要另外明确约定。

商城、Console、商品订单、旧积分充值及旧论坛迁移不属于当前新产品MVP的依赖。旧系统按原目录保留；将来若明确需要复用旧账号或内容，再独立设计集成与迁移。

## 当前目录

```text
daokin/                    # 一个本地Git仓库
├── README.md
├── .gitignore
└── api/
    └── go.mod             # 已初始化的Go模块
```

教程独立维护于 [daokin-tutorial](https://github.com/JeanphiloGong/daokin-tutorial)；本地建议将两个仓库放在同级目录。教程仓库目前为私有，需要访问权限。

前端和独立货币服务在对应实现阶段才创建。届时可在本仓库新增web/和currency/，Currency拥有自己的Go模块、进程和数据，服务通过明确接口协作。当前不预建空的多层业务目录，不引入go.work或共享包。

## 从这里开始

在本目录执行：

```bash
cd api
go list -m
```

输出为 `github.com/JeanphiloGong/daokin/api`。这是本仓库API的Go模块名称。当前没有业务Go文件，暂时没有可运行的服务或业务测试；写出首个完整业务能力后再运行相应检查。

教程与代码分开存放。教程已编写V0–V6正文及逐篇源码，运行检查在隔离副本中进行；本仓库仍是初始化状态。入口：[新项目教程索引](https://github.com/JeanphiloGong/daokin-tutorial/blob/main/README.md)、[V0完整业务模型](https://github.com/JeanphiloGong/daokin-tutorial/blob/main/v0-domain-model.md)、[V0–V6实施方案](https://github.com/JeanphiloGong/daokin-tutorial/blob/main/plan.md)。从V1.1的身份与主动成员规则开始，按教程逐篇在本仓库实现。

V1–V6实现正文位于同级教程目录，本仓库的业务代码尚未实现。旧V1–V8教程针对已有项目，不能把其patch直接应用到这里，也不能把旧教程测试结果算作本项目通过。商城、Console和旧论坛迁移不在当前主线。

## 验证

- 当前main从初始化项目开始实现；旧版原型保存在 [legacy/prototype-2026-11-13](https://github.com/JeanphiloGong/daokin/tree/legacy/prototype-2026-11-13)，旧提交和议题继续保留。
- 已确认API模块声明可解析，且未引入第三方依赖。
- 已提供环境文件、密钥、依赖目录与构建输出的忽略规则。
- 本次未添加业务抽象或兼容层；没有复制缓存或启动临时服务。
