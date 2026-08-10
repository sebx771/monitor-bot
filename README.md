# Monitor Bot

Bot de monitoreo y control de servidores de Minecraft, escrito en Go.

Proyecto de práctica que automatiza tareas de administración de servidores mediante un navegador controlado (Playwright + Chromium) y verificación de estado por ping. Soporta el arranque automático de servidores de **Aternos** y la gestión de servicios en la nube de **Aiven**, con arquitectura preparada para nuevas implementaciones.

## Funcionalidad actual

- Verifica si el servidor de Minecraft está online mediante ping.
- Lanza un navegador Chromium controlado con Playwright.
- Persiste y reutiliza la sesión del navegador mediante `StorageState` (cookies/localStorage).
- Automatiza el arranque del servidor en Aternos:
  - Navegación al panel de servidores.
  - Selección del servidor por ID.
  - Clic en el botón de inicio y manejo de diálogos (EULA, avisos, etc.).
  - Detección de sesión no válida con error descriptivo.
- Espera (polling) a que el servidor responda al ping antes de finalizar.
- Worker programable para ejecución periódica de tareas con cooldown ante fallos.
- Gestión de servicios en la nube de **Aiven**:
  - Verifica el estado de servicios por proyecto.
  - Inicia automáticamente servicios apagados (`POWEROFF`).
  - Soporta múltiples proyectos Aiven en paralelo.

## Tecnologías

- Go 1.24
- [playwright-go](https://github.com/mxschmitt/playwright-go)
- Chromium
- [go-mcping](https://github.com/iverly/go-mcping)

## Arquitectura

```
cmd/                    Aplicación CLI principal
internal/
├── ports/              Interfaz BrowserManager (abstracción del navegador)
├── adapters/           Implementación con Playwright
├── automation/         Automatización del arranque del servidor
├── minecraft/          Verificación de estado por ping
├── services/           Capa de servicios
├── worker/             Ejecución periódica de tareas (intervalo + cooldown)
└── aiven/              Cliente y checker para la API de Aiven
storage/                Persistencia de sesión del navegador (state.json)
```

## Requisitos

- Go 1.24+
- Node.js y `npx playwright install chromium` (para el driver de Playwright)

## Configuración

Copia `.env.example` a `.env` y completa los valores:

```
HOST= tuservidor.aternos.me
PORT= puerto_del_servidor
SERVER_ID= id_del_servidor
STORAGE_PATH= storage/state.json

# Aiven — una credencial por índice
AIVEN_TOKEN_1= token_api_aiven_1
AIVEN_PROJECT_1= nombre_proyecto_aiven_1
AIVEN_TOKEN_2= token_api_aiven_2
AIVEN_PROJECT_2= nombre_proyecto_aiven_2
```

## Uso

```sh
go run ./cmd
```

En la primera ejecución el navegador se abrirá en modo visible: inicia sesión en Aternos manualmente, cierra y vuelve a ejecutar. La sesión queda persistida en `storage/state.json` y se reutilizará en ejecuciones posteriores.

El bot ejecuta en paralelo:
1. **Worker de Minecraft**: revisa cada 24 horas si el servidor está online y lo inicia automáticamente en Aternos si está apagado.
2. **Worker de Aiven**: revisa cada 60 minutos los servicios de los proyectos configurados y enciende los que estén apagados.

## Próximos pasos

- Comandos de control del bot
- API HTTP
- Soporte para otras plataformas de hosting
- Integración con Discord/Telegram
