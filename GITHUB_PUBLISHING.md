# Как опубликовать HOI4 Focus Studio на GitHub

## 1. Репозиторий

1. Распакуй `HOI4_Focus_Studio_GitHub_Full_Project.zip`.
2. Создай репозиторий на GitHub и загрузи **содержимое** архива в корень репозитория (не сам ZIP).
3. Убедись, что `index.html`, `README.md`, папки `.github/`, `desktop-src/`, `windows-src/`, `scripts/` и `release-files/` лежат на верхнем уровне.
4. В `Settings → Pages` выбери `Deploy from a branch`, `main`, `/(root)` и сохрани. Сайт появится по адресу `https://USERNAME.github.io/REPOSITORY/`.

## 2. GitHub Releases

1. Распакуй отдельный архив `HOI4_Focus_Studio_All_Platforms_Release_v1.1.0.zip` (либо открой папку `release-files/` в полном проекте).
2. В репозитории открой `Releases → Draft a new release`.
3. Создай тег `v1.1.0`, заголовок `HOI4 Focus Studio v1.1.0`.
4. Вставь описание из `RELEASE_NOTES.md`.
5. В поле Attach binaries перетащи **все 7 дистрибутивов** и `SHA256SUMS.txt` (не нужно загружать `RELEASE_NOTES.md` как файл, если вставил описание).
6. Нажми `Publish release`.

На GitHub Releases файлы приложений доступны пользователям отдельно от ZIP с исходным кодом.

## 3. Автоматическая пересборка из GitHub Actions

В репозитории уже есть workflow `Build Windows EXE and installer` и `Build macOS and Linux applications`. В разделе `Actions` их можно запустить кнопкой `Run workflow`. Также они собираются при push тега вида `v*`. Артефакты сборки появятся в конкретных прогонах Actions — **в Releases они автоматически не публикуются**. Для нового выпуска скачай и приложи актуальные файлы в релиз вручную.

**Важно:** `release-files/` содержит уже собранный выпуск v1.1.0 и создана ради пакета «всё в одном». На будущих версиях не используй эти бинарники без пересборки. Их можно не коммитить, если предпочитаешь хранить дистрибутивы только в GitHub Releases.

## Где хранятся данные

Проекты сохраняются только в браузере/локальном профиле пользователя. GitHub Pages и GitHub Releases не являются хранилищем пользовательских проектов. Для переноса и восстановления используй JSON-копии из редактора.

Подробнее о редакторе и требованиях к ОС см. `README.md`.
