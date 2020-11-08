# IoT Gateway

httpプロトコルを用いてIoT機器を排他制御するためのサーバー。
socketサーバーに接続したIoT機器に対して、web側から管理、制御を行うことができる。

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

```json
GET /devices
    接続済みのIoT機器一覧を表示する。
    Operables以下の情報は削除予定
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
    
```
	r.HandleFunc("/devices",DeviceList)
	r.HandleFunc("/devices/{name}",DeviceDetail)
	r.HandleFunc("/units",UnitList).Methods("GET")
	r.HandleFunc("/units",MakeUnit).Methods("POST")
	r.HandleFunc("/units/{name}",UnitDetail).Methods("GET")
	r.HandleFunc("/units/{name}",MakeBooking).Methods("POST")
	r.HandleFunc("/units/{name}/{operable}",Operate)

