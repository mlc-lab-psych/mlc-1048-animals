# MLC-1048 "Animals" ![Static Badge](https://img.shields.io/badge/weather-%23b16ded?style=flat&logo=github&logoColor=black&labelColor=0%2C0%2C0&link=https%3A%2F%2Fgithub.com%2FweatherGod3218%2F)

Source code for the RIT MLC Pysch Labs experiment MLC-1048. Nicknamed "Animals"

This project uses Go, [Gin](https://gin-gonic.com/en/), HTML/CSS, Javascript, and the jsPysch library.


## Env Variables
These projects use alot of environmental secrets in order to hide database and airtable data. The repo includes an attached .env.template file for guranteeing env variables are named properly

```
TRAIL_NAME=mlc-1048 # Used For Internal naming schemes in databases

FIREBASE_CREDENTIALS_JSON= # JSON Credentials for authenticating into the Firebase
FIREBASE_DATABASE_URL= # URL for connecting to the Firebase

REDIS_URL= # URL to connect to the redis database

AIRTABLE_API_KEY= # API Key for authenticating with the experiments Airtable
AIRTABLE_BASE= # Base for airtables in the trail

AIRTABLE_TABLES= # JSON of every table and their URL. Key is used as the "name" of the table
```

