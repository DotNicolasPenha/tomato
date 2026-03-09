To start a project, first go into a folder and create a ```manifest.json``` file, and write this:

````
{
    "port":3000,
    "nameApplication":"your name application"
}

````

The ```manifest.json``` file is an input file for a Tomato application, it is crucial.

Now, let's create routes. First, create a folder with your preferred name, I'll name it "hello". Next, inside the "hello" folder, we'll create an ``index.json`` file.

The `index.json` file will be the input file for declaring a route; in it, you will define the route path, the method, and the base (pre-created controller).

See this example:

````
/hello/index.json

{
    "method":"get",
    "path":"/hello",
    "base":"response-json",
    "base-configs":{
        "msg":"hi!" // this is a base parameter
    }
}

````
this example configure a simple route then response a json:

````
{
    "msg":"hi!"
}
````

To add logic to your route, you can use the base functions, which are functions already created by 'tapi', a built-in distribution of your setup. To use the functions, first call them using the ``base`` property. To call them, simply enter their name. You can also give them parameters using the ``base-configs`` property.


## Running a project

Now in your terminal use the binare file of this repository or the latest release to run your project:

````
./tm.exe lta run
````

With this command you can init a ``LTA``
(Local Tomato Application) in your terminal.

## Environment variables in JSON

You can use environment variables in **any JSON file** by writing strings in this format:

````
"@env:YOUR_ENV_KEY"
````

Example:

````
{
  "dsn": "@env:DB_DSN"
}
````

## Domain schemas (`@domain/entitys` and `@domain/dtos`)

Tomato now loads schemas from:

- `./@domain/entitys`
- `./@domain/dtos`

Example entity:

````
./@domain/entitys/User.json

{
  "name": "User",
  "table": "users",
  "fields": {
    "id": {"type": "int", "primaryKey": true, "autoIncrement": true},
    "name": {"type": "string", "required": true, "min": 3, "max": 100}
  }
}
````

Example DTO:

````
./@domain/dtos/CreateUserDTO.json

{
  "name": "CreateUserDTO",
  "fields": {
    "name": {"type": "string", "required": true, "min": 3, "max": 100}
  }
}
````

You can use these schemas as types in route body validation using:

- `@entity:SchemaName`
- `@dto:SchemaName`

## SQL CRUD base

Use base `sql-crud` to execute CRUD operations based on entities:

````
{
  "method": "post",
  "path": "/users",
  "base": "sql-crud",
  "base-configs": {
    "operation": "create",
    "entity": "User",
    "driver": "postgres",
    "dsn": "@env:DB_DSN"
  },
  "request-requiredFormat": {
    "body-type": "@dto:CreateUserDTO"
  }
}
````

Supported operations:

- `create`
- `read`
- `update`
- `delete`
- `list`
