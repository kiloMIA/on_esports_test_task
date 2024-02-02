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

Reminder:
!remindme 1m Take a break!

Weather:
!weather New York