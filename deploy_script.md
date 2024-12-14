## Описание скрипта запуска

Этот скрипт автоматизирует сборку, запуск, остановку и очистку ресурсов Docker для вашего VPN-приложения. Он также обрабатывает системные сигналы для корректной остановки контейнера.

1. Функция cleanup

    - Что делает:
        - Останавливает и удаляет контейнер с именем vpn_container.
        - Выводит сообщение о завершении работы.
        - Завершает выполнение скрипта с кодом 0 (успешно).

    - Зачем: Для обеспечения корректной очистки ресурсов при завершении работы приложения, например, после нажатия `Ctrl+C`.

2. Обработка сигналов с помощью trap

    - Что делает:
        - Перехватывает сигналы `INT` (прерывание, например, `Ctrl+C`) и `TERM` (остановка).
        - При получении этих сигналов автоматически вызывает функцию `cleanup`.

    - Зачем: Позволяет избежать зависания контейнера в активном состоянии или оставления ненужных ресурсов после завершения работы приложения.
3. Режим 1: `start`
     - Использует vpn.dockerfile для сборки образа с тегом vpn_ubuntu_all_files_include
     - Запрашивает порт у пользователя. По умолчанию порт — 12345, если пользователь нажимает Enter без ввода.
     - Запускает контейнер ```docker run --cap-add=NET_ADMIN --cap-add=NET_RAW --network host --device /dev/net/tun -e PORT=$PORT --name vpn_container -it vpn_ubuntu_all_files_include```:
          - Передаёт порт через переменную среды PORT.
          - Использует:
              - cap-add=NET_ADMIN Этот флаг добавляет в контейнер способность управления сетевыми настройками, которая по умолчанию запрещена из соображений безопасности. Позволяет изменять настройки сети, например:
                - Создавать и удалять интерфейсы (TUN/TAP, виртуальные сети и т.д.).
                - Изменять таблицы маршрутизации.
                - Управлять правилами фильтрации IP (iptables).
                - Включать или отключать сетевые устройства.
                - Разрешает выполнение таких команд, как:
                  - ip link set (управление интерфейсами).
                  - ip route (управление маршрутами).
                  - ifconfig (конфигурация интерфейсов).

              - cap-add=NET_RAW — Этот флаг добавляет контейнеру доступ к работе с "сырыми" (raw) сокетами. Позволяет отправлять и получать "сырые" сетевые пакеты, минуя обработку протоколов на уровне операционной системы.
                  Используется для:
                
                  - Работы с нестандартными протоколами.
                  - Изучения сетевого трафика (анализатор пакетов).
                  - Разработки или тестирования собственных протоколов передачи данных.


              - network host — подключает контейнер к сети хоста.
              - device /dev/net/tun — разрешает доступ к TUN-устройству.
    - Ожидает завершения контейнера. `docker wait` останавливает выполнение скрипта, пока контейнер не завершит работу.
    - Вызывает `cleanup` после завершения контейнера.
  
4. Режим 2: `clear-cache`
   - Удаляет образ с тегом vpn_ubuntu_all_files_include
   - Очищает промежуточные образы. Удаляет все неиспользуемые Docker-образы, сети и тома
  
## `run_vpn.sh`
```
#!/bin/bash

cleanup() {
    echo "Stopping Docker container..."
    docker stop vpn_container
    docker rm vpn_container
    echo "Container stopped."
    exit 0
}

trap cleanup INT TERM

if [[ $1 == "start" ]]; then
    echo "Building Docker image..."
    docker build -f vpn.dockerfile --tag vpn_ubuntu_all_files_include .

    read -p "Сhoose a port to open later: " port
    port=${port:-12345}
    
    export PORT=$port
    echo "Starting Docker container..."
    docker run  --cap-add=NET_ADMIN --cap-add=NET_RAW --network host --device /dev/net/tun -e PORT=$PORT --name vpn_container -it vpn_ubuntu_all_files_include
    
    docker wait vpn_container
    
    cleanup

elif [[ $1 == "clear-cache" ]]; then
    echo "Removing Docker image and clearing cache..."
    docker rmi vpn_ubuntu_all_files_include >/dev/null 2>&1 
    
    docker image prune -a -f
    
    echo "Cache cleared."

else
    echo "Usage: $0 {start|clear-cache}"
fi

```
