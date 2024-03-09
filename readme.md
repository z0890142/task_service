# 啟動方式
於根目錄下執行 make build，進行 docker compose build，在進行 make start
如沒有安裝 make 環境可直接參考 makefile 內的指令。

## makefile 啟動方式
```
make start
```

## docker 啟動方式
```
docker-compose up -d  
```

## 產生文件

於根目錄執行以下指令，執行完後位於根目錄下的 /docs 資料夾裡
```
make generate_doc
```

---

### list task api order 參數說明
如需根據 ID 由小排到大請輸入 => id 
如需根據 ID 由大排到小請輸入 => id desc

**範例**
```
curl --location 'http://127.0.0.1:8080/task-service/api/v1/tasks?order=id%20desc&limit=1&offset=1'
```