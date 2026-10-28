# HOI4 Focus Studio — v1.1.0

Выпуск готовых приложений и веб-редактора национальных фокусов Hearts of Iron IV.

## Скачать для своей системы

| ОС | Файл | Установка |
|---|---|---|
| Windows 10/11 (x64) | `HOI4-Focus-Studio-Setup.exe` | Запустить установщик (нужен Edge) |
| Windows 10/11 (x64) | `HOI4-Focus-Studio.exe` | Портативная версия: запустить без установки |
| macOS 11+, Intel / Apple Silicon | `HOI4-Focus-Studio-macOS-universal.zip` | Распаковать и перенести `.app` в Applications |
| Ubuntu / Debian / Mint x86-64 | `HOI4-Focus-Studio-linux-amd64.deb` | `sudo apt install ./HOI4-Focus-Studio-linux-amd64.deb` |
| Linux x86-64 | `HOI4-Focus-Studio-linux-amd64-portable.tar.gz` | Распаковать и запустить `HOI4-Focus-Studio` |
| Ubuntu / Debian / Mint ARM64 | `HOI4-Focus-Studio-linux-arm64.deb` | `sudo apt install ./HOI4-Focus-Studio-linux-arm64.deb` |
| Linux ARM64 | `HOI4-Focus-Studio-linux-arm64-portable.tar.gz` | Распаковать и запустить `HOI4-Focus-Studio` |

## Возможности

- Визуальный редактор фокусов HOI4 и связей И/ИЛИ/взаимоисключений.
- Несколько проектов, локальное автосохранение, откат действий и JSON-резервные копии.
- Экспорт фокусов в скрипты HOI4, локализации и описания для ИИ.
- Версия для GitHub Pages в исходном репозитории.

**Важно:** редактор использует установленный браузер и локальный веб-сервер. Windows требует Edge, macOS и Linux используют Chrome/Chromium/Edge/Brave или браузер по умолчанию. Между устройствами проекты переносятся через JSON, автоматической облачной синхронизации нет.

**Сборки без подписи:** Windows SmartScreen и macOS Gatekeeper могут показывать предупреждения. Приложения собраны в Linux-среде; запуск на реальных Windows/macOS не подтверждён. Файл `.dmg` не включён: его можно получить через GitHub Actions на macOS runner.

Сверяйте файлы с `SHA256SUMS.txt` перед публикацией. Приложение не является официальным продуктом Paradox Interactive.
