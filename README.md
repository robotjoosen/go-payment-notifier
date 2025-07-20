# Payment notifier

Play a sounds when a payment is done

## Build

| type  | command  |
|---|---|
| docker | `make docker-build` |
| binary | `make binary-build` |

## Run

1. Create `.env` by copying `.env.dist`
  ``` shell
  cp .env.dist .env
  ```
2. Start notifier
  - Docker: 
    ``` shell
    make docker-run
    ```
  - Binary:
    ``` shell
    ./bin/notifier
    ```

## Interact

| endpoint  | description |
|---|---|
| `/health` | health endpoint that can be used for docker or k8s |
| `/shutdown` | gracefuly shutdown service through endpoint |
| `/callbacks/payment` | callback endpoint for handling bunq payment callbacks |
| `/callbacks/mutation` | callback endpoint for handling bunq mutation callbacks |

