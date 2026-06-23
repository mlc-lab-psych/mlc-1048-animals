# mlc-labs-template ![Static Badge](https://img.shields.io/badge/weather-%23b16ded?style=flat&logo=github&logoColor=black&labelColor=0%2C0%2C0&link=https%3A%2F%2Fgithub.com%2FweatherGod3218%2F)
A Boilerplate for easier creation of MLC Pysch Experiment Websites

This project uses Golang, [Gin](https://gin-gonic.com/en/), HTML/CSS, Javascript, and the jsPysch library.


## Env Variables
These projects use alot of environmental secrets in order to hide database and airtable data. The repo includes an attached .env.template file for guranteeing env variables are named properly
```
FIREBASE_CREDENTIALS_JSON=
FIREBASE_DATABASE_URL=

AIRTABLE_API_KEY=
AIRTABLE_BASE=

AIRTABLE_TABLE1=
```
*Note: AIRTABLE_TABLE1 serves as the link for the first table, these can be added as needed. The naming convention is AIRTABLE_TABLE{$TABLE}*