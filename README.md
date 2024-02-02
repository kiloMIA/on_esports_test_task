# User Guide

This Discord Bot provides you such functionalities: Setting Reminders, Getting Weather Forecast, Making Polls

## Installation
Install and configure Docker [Guide](https://docs.docker.com/engine/install/)

Clone this repository
```bash
git clone https://github.com/kiloMIA/on_esports_test_task.git
```

To start this application you can run this command
```bash
make docker-up
```
Or directly via docker-compose
```bash
docker compose up --build
```
## Examples

Polls:
!poll What's your favorite color? | Red | Blue | Green
![poll_screen](https://github.com/kiloMIA/on_esports_test_task/assets/97970527/6797358c-3b2c-4d92-8324-d2e99a55c83f)

Reminder:
!remindme 1m Take a break!
![reminder_screen](https://github.com/kiloMIA/on_esports_test_task/assets/97970527/bbf0259b-cee8-461e-8ffb-25167cb685f0)

Weather:
!weather New York
![weather_screen](https://github.com/kiloMIA/on_esports_test_task/assets/97970527/f51e0b82-f9f6-49fd-8c5e-aa6547552024)
