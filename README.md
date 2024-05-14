[![codecov](https://codecov.io/gh/nestjam/goph-keeper/branch/main/graph/badge.svg?token=4UDX8BV3G7)](https://codecov.io/gh/nestjam/goph-keeper)

GophKeeper keeps secrets.

## Запуск сервера и клиента

1. Запустить сервер с указанием файла конфигурации и мастер ключа хранилища (32 символа).
    - Для работы сервера по HTTPs необходимо в файле конфигурации указать путь к файлу сертификата и файлу приватного ключа.

    ```sh
    go run main.go -c ../../internal/config/config.yml -k=N3SaEN8k2z3?DCf_4_8j+Yc92pTrFt6W
    ```

2. Собрать и запустить клиент с указанием адреса сервера

    ```sh
    go build -o client.exe -ldflags "-X main.BuildVersion=v0.0.1 -X 'main.BuildDate=$(date +'%Y/%m/%d')'" main.go

    start client.exe -s https://localhost:8080
    ```