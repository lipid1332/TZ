# Тех задание ITK academy
## ⚙️ Установка и запуск
1. git clone https://github.com/lipid1332/TZ
2. cd yourprojectdir
---
### Windows (PowerShell)
3. Выполнить данную команду в оболочке:

```
  Get-Content .\config.env | Where-Object { $_ -notmatch '^#' } | ForEach-Object {
      $parts = $_ -split '=', 2
      if ($parts.Length -eq 2) {
          [System.Environment]::SetEnvironmentVariable($parts[0], $parts[1])
      }
    }
```
4. В Docker Desktop разрешить запуск файла миграции (либо добавить папку в File Sharing)
5. Запустить docker-compose up
### Linux / macOS (bash)
3. Выполнить команду в терминале:

```
export $(grep -v '^#' config.env | xargs)
```
4. Запустить docker-compose up
---
# 🚀 Тесты
Запустить легковесный функциональный тест:

```
go test -v -run TestPingRoute test/functional_test.go
```
