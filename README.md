## Генерация кода
``` sh
make gen
```

### Если не работает соединение pgAdmin с Postres на WSL: 
1. Узнать новый IP в WSL:
        `hostname -I`
2. Обновить правило portproxy в powershell (запуск через админа):
    `netsh interface portproxy set v4tov4 listenport=5432 listenaddress=<ip wsl> connectport=5432 connect`