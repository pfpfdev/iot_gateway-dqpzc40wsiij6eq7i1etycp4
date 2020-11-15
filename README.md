# IoT Gateway

httpプロトコルを用いてIoT機器を排他制御するためのサーバー。
socketサーバーに接続したIoT機器に対して、web側から管理、制御を行うことができる。

ソースコードのコメント/リファクタリングは基本しない予定です。
使い捨てのコードで適当に書いた感じなので

## 使い方

### IoT機器側

規定のポート番号に対して(8081)、TCPコネクションを立てる。
後述するプロトコルにより、http側からの制御を受けることができる。

#### プロトコル

TCPコネクションを確立した機器は固有の名前持つ。

通信は`\n`をデリミタとする一行を単位とする。
通信の中では語を` `で区切るCSV方式で通信する。

実際の通信は以下のような流れで行う。

```
--- Initial phase

IoT > <NAME>\n
IoT > <CMD> ...<ARGS>\n
IoT > FIN\n

--- Normal phase

Server > <CMD> <ARG>\n
IoT > <RES>\n
...

```

Initial phaseで実行することができるのは、以下のコマンド

+ ADD [OPERABLE_NAME]
  + 制御対象の機器(operable)を追加する
+ REG [OPERABLE_NAME] [COMMAND_NAME] [ARG_TYPE]
  + 追加済みのoperableにたいして、コマンドを定義する
  + [ARG_TYPE]には`OnOff`,`Hundred`の二種類があり、前者は引数として`On|Off`,`[0-9]{2}`の文字列が送信される

たとえば次のように送信することで、{LED,speed}と{power}を制御することができるIoT機器が接続されているとみなす。

```
Device1 
ADD ope1
REG ope1 LED OnOff    
REG ope1 speed Hundred
ADD ope2
REG ope2 power OnOff  
FIN
```

Normal phaseではIoT機器から情報を受け取ることはせず、http側からの要求に基づいて定義されたコマンドが呼び出される。

上の場合では

```
ope1 LED On
ope1 speed 50
ope2 Off
```

などが送信されるので、それをハンドリングしレスポンスとして1行を送信する。

Normal phaseでは一定時間(10秒)以上通信がないIoT機器は落ちたと判断してコネクションを切断する。
それを避けるためには定期的に(3秒程度)`\n`だけを送信する必要がある。

### HTTP側

HTTPサーバーは規定のポート(8080)で実行される。
APIのエンドポイントは以下のようにマッピングされている。

```js
GET /devices
    接続済みのIoT機器一覧を表示する。
    ~~Operables以下の情報は削除予定~~
    削除しません。詳細表示が意味なくなるけど、意外と削除するのがめんどくさそうだったので今回はこのままいきます
    {
        "DeviceName":{
            "Name":"DeviceName",
            "Operables":{
                "OperableName":{
                    "Name":"OperableName",
                    "Operations":{
                        "CommandName":{
                            "Cmd":"CommandName",
                            "Type":"OnOff|Hundred"
                        },...
                    }
                },...
            }
        },...
    }
GET /device/{name}
    デバイスの詳細情報を表示する
    {
        "Name":"DeviceName",
        "Operables":{
            "OperableName":{
                "Name":"OperableName",
                "Operations":{
                    "CommandName":{
                        "Cmd":"CommandName",
                        "Type":"OnOff|Hundred"
                    },...
                }
            },...
        }
    }
GET /units
    ユニットの情報を一覧表示する
    {
        "UnitName":{
            "Name":"UnitName",
            "Operables":{
                "ope1":{
                    "Name":"ope1",
                    "Operations":{
                        "LED":{
                            "Cmd":"LED",
                            "Type":"OnOff"
                        },
                        "speed":{
                            "Cmd":"speed",
                            "Type":"Hundred"
                        }
                    }
                },
                "ope2":{
                    "Name":"ope2",
                    "Operations":{
                        "power":{
                            "Cmd":"power",
                            "Type":"OnOff"
                        }
                    }
                },...
            },
            "Queue":[
                {
                    "LastTime":"2020-11-13T21:33:16.9206603+09:00",//最終アクセス(有効化)時刻
                    "Until":"2020-11-13T21:34:28.7914384+09:00",//制御権の有効期限の目安(割り当てられてないならば意味なし)
                    "IsAlive":true
                },...
            ],
        },...
    }
POST /units
    ユニットの構成を設定する
    以下のフォーマットをBODYで送信する
    BASIC認証が必要
    {
        "UnitName":{
            "DeviceName of operables":["OperableName1","OperableName1"],
            "Other Device":["OperableName1","OperableName1"]
        }
    }
POST /units/{unitName}
    整理番号(制御に必要なトークン)の確保
    ここで確保したら/units/{unitName}で順番確認
    定期的にtokenの有効化を行う
    {
        "Token":randomuint64
    }
GET /units/{unitName}?token={token}
    tokenの更新を行う
    定期的にここにアクセスしTokenの有効化を行う(10~30秒に一度程度)
    返り値は現在の順番(1が先頭)
    {
        "Order":1
    }    
    トークンが指定されなければユニットに関する情報を取得する
    この場合の返り値は/unitsの限定的なものにつき省略
GET /units/{unitName}/{operableName}?cmd={cmdName}&arg={arg}
    操作を行う
    unitsのUserのIdと同じtokenの人が操作できる
GET /log?offset={offset}
    ログの取得
    Offsetを指定すると続きを読み取れる
    ラウンドロビンとか考えてないので定期的に再起動したほうがいいかもしれない{要検証}
    BASIC認証が必要
    5秒に一度くらい読み取って、logファイルへのアクセスログをフィルタリングすれば良い気が
    {
        "Log":"LogText\n",
        "Offset":123
    }
```

## 実装

### 順番管理の戦略

1. 整理番号を取得した人を並べる
2. 先着順で制御権を与える
3. 一定時間後に制御を制御権を没収
4. 次の順番の人に制御権を渡す
   1. ただし一定時間の間に最低一回tokenの有効化を行った人に限る

#### Token

unit64の疑似乱数

暗号的に全く持って安全でないので、ばれるリスクはあるもののそこまで気にする必要はなさげ

~~というか現在の実装だと/unitsにアクセスしてしまえば権限横取りし放題~~改善済み

それは改善予定だが、URL queryはhttpsでも暗号化されないので、パケット見れば読み取り放題

Go言語だとjsonのparseが(丁寧にやらないと)比較的面倒なので、元気があれば直すかも

## TODO

+ テストの作成
