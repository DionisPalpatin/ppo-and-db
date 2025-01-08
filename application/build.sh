#!/bin/zsh

# Определяем пути
SOURCE_STATIC="/servers/notebook-app/static"
SOURCE_DOCS="/servers/notebook-app/docs"
CONTEXT_PATH="/home/daemo/SpecialForCP/ppo-and-db/application"
TARGET_STATIC="${CONTEXT_PATH}/static"
TARGET_DOCS="${CONTEXT_PATH}/docs"
DOCKERFILE_PATH="${CONTEXT_PATH}/Dockerfile"

# Проверяем существование исходных директорий
if [[ ! -d "$SOURCE_STATIC" ]]; then
  echo "Ошибка: Папка $SOURCE_STATIC не найдена!"
  exit 1
fi

if [[ ! -d "$SOURCE_DOCS" ]]; then
  echo "Ошибка: Папка $SOURCE_DOCS не найдена!"
  exit 1
fi

# Копируем папки в контекст Docker
echo "Копирование папок в контекст сборки..."
cp -r "$SOURCE_STATIC" "$TARGET_STATIC"
cp -r "$SOURCE_DOCS" "$TARGET_DOCS"

# Запускаем сборку Docker
echo "Запуск сборки Docker compose..."
docker compose build

# Сохраняем статус сборки
BUILD_STATUS=$?

# Удаляем временные папки
echo "Удаление временных папок..."
rm -rf "$TARGET_STATIC" "$TARGET_DOCS"

# Проверяем статус сборки
if [[ $BUILD_STATUS -eq 0 ]]; then
  echo "Сборка Docker успешно завершена!"
else
  echo "Ошибка при сборке Docker!"
  exit $BUILD_STATUS
fi
