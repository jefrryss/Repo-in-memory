# In-Memory Key-Value Database

 in-memory база данных формата «ключ-значение» на Go. Поддерживает клиент-серверное взаимодействие по TCP и механизм Write-Ahead Logging (WAL) для предотвращения потери данных при перезагрузке.

## Что реализовано

* **Хранилище:** Потокобезопасная хэш-таблица в оперативной памяти.
* **Язык запросов:** Собственный парсер для команд `SET`, `GET` и `DEL`.
* **Сеть:** TCP-сервер с многопоточной обработкой клиентов (с лимитом соединений) и отдельное приложение CLI-клиента.
* **Надежность (WAL):** Запись мутирующих операций на диск батчами (с `fsync`), сегментация файлов лога и автоматическое восстановление данных при старте.
* **Конфигурация:** Настройка сети, параметров WAL и логирования через `config.yaml`.

## Запуск сервера
Для запуска сервера с настройками по умолчанию:
```bash
go run cmd/server/main.go
```


```Bash
go run cmd/client/main.go --address=127.0.0.1:3223
```
Пример конфигурации (config.yaml)
```YAML
engine:
  type: "in_memory"

wal:
  flushing_batch_size: 100
  flushing_batch_timeout: "10ms"
  max_segment_size: "10MB"
  data_directory: "/data/wal"

network:
  address: "127.0.0.1:3223"
  max_connections: 100
  max_message_size: "4KB"
  idle_timeout: "5m"

logging:
  level: "info"
  output: "/log/output.log"
```