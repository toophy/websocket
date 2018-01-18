# 文字MUD游戏

---

## 仓库

存储玩家道具和标签货币

## 好友

好友多也是一种标签货币:)

## 邮件

系统发放使用

## 战队

影响房间战斗收益, 聊天频道

## 房间

游戏核心, 快乐的战斗之后, 产出各种标签货币、道具碎片, 战斗中只能和队友聊天(语音),
禁止一切机会和敌方聊天, 禁止串通作弊, 一旦出现奇怪的行动, 不攻击等, 视为作弊
对消极战斗, 全面给与惩罚, 全面给出打酱油, 评分低, 而且信用降低, 产生一次信用问题记录

## 大厅

各种系统的入口, 统计中心

## 标签货币

非常多, 从开服第一场战斗开始, 建立标签货币产出总值, 现存总值
货币的用途就是兑换钻石, 不管什么货币, 只要比上一个出价高(依据当时汇率), 就成功竞价
竞拍成功后, 规定时间内汇款, 否则流拍钻石回归拍卖场所
当前汇率是全服可见, 从最值钱的到最便宜的一次排列, 也就是总值从小到大排列, 每1分钟统计一次均值,
对于前面几种罕见货币, 进行录像审查, 审查后, 才能作为货币反馈
标签货币就是说 : 用游戏的过程赚钱, 付出就有回报, 不管是胜利还是失败, 大家都牛逼了, 失败就是最大赢家

## 下列是标签货币

货币名称|描述
--|--
集邮王|碎片收集足够多
吃货|一场战斗中吃的最多
好基友|一场战斗中助攻最多
技术达人|一场战斗中杀人最多
打酱油|一场战斗中, 啥成就也不突出
天使|一场战斗中, 死亡次数最少
开场秀|开战第一滴血
超神|10杀
我是第一|团队第一名

## 游戏系统

```mermaid
graph LR

充值-->|获取|钻石((真钻石))
钻石-->|竞拍|通货
战斗((房间战斗))-->|产出|通货
战斗-->|产出|碎片
碎片-->|合成/维修|装备
通货-->|合成/维修|装备
装备-->|消耗|战斗
```

### 装备升阶

```mermaid
graph LR

钻石-->|增加幸运|装备((装备升阶))
好基友-->|有益孔链接|装备
开场秀-->|有益第一属性出攻击|装备
我是第一名-->|有益加成所有属性|装备 

```

### 匹配流程

```mermaid
graph LR

游戏模式-->|选择后进入二级界面|选择队友
选择队友-->|等待|提交匹配
提交匹配-->|等待 \ 最多30s|提交匹配
提交匹配-->|成功|选择英雄
选择英雄-->|成功|选择装备
选择装备-->|成功|进入游戏

```

### 钻石矿 && 拍卖会

强调RMB玩家的钻石地位, 特别开通钻石矿, 强调钻石的作用和地位

```mermaid
graph LR

钻石矿-->|系统变速定期产出|钻石
钻石-->|系统定期开启|拍卖会
被承认的通货-->|参与竞价|拍卖会
每种通货唯一拍卖场-->拍卖会

```

---

## 主题

### 冒险 + 竞技

- 冒险
    > 场景随机, 各种随机触发陷阱, 有利和不利的都有, 可怕的草丛, 乌黑的山洞

- 竞技
    > 英雄自带装备+技能占1/3权重
    > 场景提供的道具和陷阱占1/3权重
    > 英雄走位占1/3权重
    > 英雄配合就是全部

---

## 主界面

- 帐号相关
  基本信息和统计系统

- 冒险模式
> 1. 组队/单排乱斗
>    背包     : 临时背包, 开始冒险为空, 玩家只能携带自身装备, 临时背包很小(6格), 拾取的道具放入临时背包
>    宠物     : 战斗前携带宠物, 冒险中不能切换宠物, 宠物可以携带装备
>    死亡惩罚 : 随机掉落1件装备, 优先掉落比赛中拾取的装备, 该装备耐久掉1, 队友拾取, 靠近自动归还掉落者, 敌人拾取就归敌人本场战斗使用
>    冒险结束 : 依据各种标记获取战利品, 归还拾取的玩家装备(装备损坏提醒玩家), 各种通货和奖牌(也是道具) 
> 2. 训练场
>    单机训练模式
> 3. 开房间
>    组队训练模式, 不能刷分, 不掉耐久, 不能带走产出道具, 只能产出"房卡"
>
 


- 钻石矿+拍卖会
    > 系统馈赠的钻石, 非R需要用游戏币竞价获取

- 仓库
    > 没有格子概念, 分页分类保存, 只有重量限制, 玩家可以购买重量

- 装备强化
    > 用通货和冒险获取的装备进行强化, 过程是个迷, 没有固定配方, 一切都是迷
    > 升阶强化非常重要, 将影响装备的最高耐久, 之后每次修理最高耐久就下降
    > 再好的装备用的多很快就报废, 无法使用

- 市场
    > 用通货交易装备和通货, 不能用同一类通货交换自身, 每个交易物品单独标价
    > 市场都是一口价
    > 查询物品( 物品, 通货, 价位, 品质, )


## 新手引导
- 英雄引导
```mermaid
graph LR
进入游戏-->|选择英雄|训练场
训练场-->|游戏结束|掉落1件10耐久白装
```

- 装备引导
```mermaid
graph LR
选择训练场模式-->|选择英雄|选择装备
选择装备-->训练场战斗
训练场战斗-->|游戏结束|掉落1只白宠物
```

- 宠物引导
```mermaid
graph LR
选择训练场模式-->|选择英雄|选择装备
选择装备-->|选择宠物|训练场战斗
训练场战斗-->|游戏结束|掉落多种装备强化道具
```

- 装备强化引导
```mermaid
graph LR
选择装备强化分页-->选择装备
选择装备-->选择强化道具
选择强化道具-->增减道具
增减道具-->点击锤子强化
```

- 实战引导  *可以省略装备和宠物的选择, 使用默认配置* 
```flow
st=>start: 开始
e=>end: 结束

cond1=>condition: 冒险?
cond2=>condition: 训练?
cond3=>condition: 开房间?

op1=>operation: 匹配
op2=>operation: 进入场景
op3=>operation: 选择一名英雄
op4=>operation: 选择队友
op5=>operation: 选择武器
op6=>operation: 选择宠物

xop1=>operation: 匹配
xop2=>operation: 进入场景
xop3=>operation: 选择一名英雄
xop4=>operation: 选择队友
xop5=>operation: 选择武器
xop6=>operation: 选择宠物

kop1=>operation: 匹配
kop2=>operation: 进入场景
kop3=>operation: 选择一名英雄
kop4=>operation: 选择队友
kop5=>operation: 选择武器
kop6=>operation: 选择宠物


st->cond1
cond1(yes)->op4->op3->op1->op5->op6->e
cond1(no)->cond2

cond2(yes)->xop4->xop3->xop5->xop6->e
cond2(no)->cond3

cond3(yes)->kop4->kop3->kop5->kop6->e


```

- **市场购买流程**
```sequence
User->Hall: 请求购买商品
Hall->Market: 转发请求购买商品
Market->Shops: 购买商品
Shops-->Market: 购买失败,反馈原因
Shops->MailSys: 购买成功,邮寄商品
Market-->Hall: 转发购买失败,反馈原因
Hall-->User: 提示购买失败,反馈原因
MailSys->Hall: 用户在线,执行脚本邮件
Hall-->User: 提示商品购买成功
Hall-->MailSys: 收货成功,同意打款给商铺
MailSys-->Shops: 商铺在线,执行打款
```

- **编号规则**
> 类型  :  第一字节
> 区服  :  第二字节和第三字节, 这里其实是线程编号, 每个区服就是一个线程
> 编号  :  第四字节到第八字节

> 类型  :  1 帐号(邮箱)
> 类型  :  2 装备
> 类型  :  3 玩家角色
> 类型  :  4 智能角色


- **帐号仓库**
  按照道具分类显示, 可以自定义标签分页, 防止不同的道具, 没有标签的放入
  临时标签页

- **整服构架**
```mermaid
graph TD

Acc1-->Conn1
Acc2-->Conn2

subgraph one
Hall-->ConnY((大厅网络))
Market-->ConnY
Store-->ConnY
Room-->ConnY
Match-->ConnY
Mails-->ConnY
钻石矿-->ConnY
数据库-->ConnY
end

subgraph two
Hall-->BigC((大服中心))
Market-->BigC
Store-->BigC
Room-->BigC
Match-->BigC
Mails-->BigC
钻石矿-->BigC
数据库-->BigC
end

```
```plantuml
@startsalt
{+
{* File | Edit | Source | Refactor 
 Refactor | New | Open File | - | Close | Close All }
{/ General | Fullscreen | Behavior | Saving }
{
	{ Open image in: | ^Smart Mode^ }
	[X] Smooth images when zoomed
	[X] Confirm image deletion
	[ ] Show hidden images 
}
[Close]
}
@endsalt
```
```plantuml
@startuml
sprite $bProcess jar:archimate/business-process
sprite $aService jar:archimate/application-service
sprite $aComponent jar:archimate/application-component

rectangle "Handle claim"  as HC <<$bProcess>> #yellow 
rectangle "Capture Information"  as CI <<$bProcess>> #yellow
rectangle "Notify\nAdditional Stakeholders" as NAS <<$bProcess>> #yellow
rectangle "Validate" as V <<$bProcess>> #yellow
rectangle "Investigate" as I <<$bProcess>> #yellow
rectangle "Pay" as P <<$bProcess>> #yellow

HC *-down- CI
HC *-down- NAS
HC *-down- V
HC *-down- I
HC *-down- P


CI -right->> NAS
NAS -right->> V
V -right->> I
I -right->> P



rectangle "Scanning" as scanning <<$aService>> #A9DCDF
rectangle "Customer admnistration" as customerAdministration <<$aService>> #A9DCDF
rectangle "Claims admnistration" as claimsAdministration <<$aService>> #A9DCDF
rectangle Printing  <<$aService>> #A9DCDF
rectangle Payment  <<$aService>> #A9DCDF

scanning -up-> CI
customerAdministration  -up-> CI
claimsAdministration -up-> NAS
claimsAdministration -up-> V
claimsAdministration -up-> I
Printing -up-> V
Printing -up-> P
Payment -up-> P

rectangle "Document\nManagement\nSystem" as DMS <<$aComponent>> #A9DCDF
rectangle "General\nCRM\nSystem" as CRM <<$aComponent>> #A9DCDF
rectangle "Home & Away\nPolicy\nAdministration" as HAPA <<$aComponent>> #A9DCDF
rectangle "Home & Away\nFinancial\nAdministration" as HFPA <<$aComponent>> #A9DCDF

DMS .up.|> scanning
DMS .up.|> Printing
CRM .up.|> customerAdministration
HAPA .up.|> claimsAdministration
HFPA .up.|> Payment

legend left
Example from the "Archisurance case study" (OpenGroup).
See 
==
<$bProcess> :business process
==
<$aService> : application service
==
<$aComponent> : appplication component
endlegend
@enduml
```

```plantuml
@startuml

package "中心服" {
  中心连接 - [中心核心]
}

package "大厅服" {
  连接1 - [帐号1]
  大厅连接 - [大厅核心]
  [大厅核心] --> 中心服
}

package "用户" {
    [用户核心] -> 连接1
}

cloud "邮件服" {
  [邮件核心] --> 大厅服
  [邮件核心] --> 中心服
}

cloud "房间服" {
  [房间核心] --> 大厅服
  [房间核心] --> 中心服
}

cloud "匹配服" {
  [匹配核心] --> 大厅服
  [匹配核心] --> 中心服
}

cloud "矿场服" {
  [矿场核心] --> 大厅服
  [矿场核心] --> 中心服
  [矿场]
  [拍卖场]
}

cloud "市场服" {
  [市场核心] --> 大厅服
  [市场核心] --> 中心服
}

database "共享服" {
  [共享核心] --> 大厅服
  [共享核心] --> 中心服
}



@enduml
```
