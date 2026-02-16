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
