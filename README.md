# dbs - a simple program for remote diagnosis
**dbs** is a simple program for remote diagnosis written in Go.

The program uses so-called *info providers*. It is just a function that provides some information in a free format. **dbs** runs it once in a certain period of time. The returned information is converted to JSON and sent to a server.

*Info providers* may, for example, provide the information about computer resources usage, tabs opened in a browser, etc.

There are already two info implemented providers - ProcessesInfoProvider and RAMInfoProvider.

## Installation and usage
You should install PostgreSQL.

You should create a .env file in the `client` directory with this information:
* `API_URL` - the API URL used to send the information provided by *info providers* to a server.

And in the `server` server with this information:
* `DEFAULT_DATABASE_URL` - the URL of the database used to create a new database (e.g. postgres://username:password@localhost:5432/postgres)
* `DBS_DATABASE_URL` - the URL of the database used by **dbs**
* `DBS_DATABASE_NAME` - the name of the database used by **dbs**
* `DBS_TABLE_URL` - the name of the table used to store the information returned by *info providers*.
* `PORT` - the port used by the server.
