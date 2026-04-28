# In-Memory Key-Value Database

Учебная in-memory база данных формата «ключ-значение», написанная на Go. Поддерживает клиент-серверное взаимодействие по TCP и механизм Write-Ahead Logging (WAL) для предотвращения потери данных при перезагрузке.

## 🚀 Что реализовано

* **Хранилище:** Потокобезопасная хэш-таблица в оперативной памяти.
* **Язык запросов:** Собственный парсер для команд `SET`, `GET` и `DEL`.
* **Сеть:** TCP-сервер с многопоточной обработкой клиентов (с лимитом соединений) и отдельное приложение CLI-клиента.
* **Надежность (WAL):** Запись мутирующих операций на диск батчами (с `fsync`), сегментация файлов лога и автоматическое восстановление данных при старте.
* **Конфигурация:** Настройка сети, параметров WAL и логирования через `config.yaml`.
* **Архитектура:** Послойная структура (компоненты разделены на network, compute, storage).

## 🛠 Запуск сервера

Для запуска сервера с настройками по умолчанию:
```bash
make run-server
# или стандартной командой:
go run cmd/server/main.go
Для запуска с конкретным файлом конфигурации:

Bash
make run-server CONFIG_FILE_NAME=config.yaml
# или стандартной командой:
go run cmd/server/main.go --config=config.yaml
💻 Запуск CLI клиента
Для запуска клиента и подключения к базе:

Bash
make run-cli
# или с явным указанием адреса сервера:
go run cmd/client/main.go --address=127.0.0.1:3223
📝 Пример конфигурации (config.yaml)
YAML
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