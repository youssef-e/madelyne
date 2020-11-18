# Madelyne advanced moveset 

 - [Json responses files](#json-responses-files)
 - [Patterns](#patterns)
 - [Capturing patterns](#capturing-patterns)
 - [Partial files](#partial-files)
 - [Optional fields](#optional-fields)
 - [Env file](#include-dynamic-content)

## Json responses files


In the yaml files, you provided a `out` parameter, this is a file that contains the expected response. A json response could look like this:

```json
{
	"email": "raphael.alves@everycheck.fr",
	"age": 25
}
```

You want to validate that email is really an email and age a number, but you sometimes or often don't care about the content. That's why we can write patterns in those files.

This is what looks like the pattern file:

```json
{
	"email": "@string@.isEmail()",
	"age": "@number@"
}
```

This file will validate the content showed above.
Patterns are composed of types (such as `@string@`) and functions (like `isEmail()`) that can be chained.

Here are the types you can use:

|Type|Description|
|--|--|
|`@boolean@`| Match booleans |
|`@number@`| Match numbers |
|`@string@`| Match strings |
|`@uuid@`| Match version 4 uuid strings |


Now, let's see what functions you can use :

|Function|Usage example|
|--|--|
|`greaterThan($boundry)`|`@number@.greaterThan(10)`|
|`lowerThan($boundry)`|`@number@.lowerThan(10)`|
|`contains($string)`|`@string@.contains('Hello')`|
|`notContains($string)`|`@string@.notContains('Hello')`|
 `endsWith($stringEnding)`|`@string@.endsWith('llo')`|
|`startsWith($stringBeginning)`|`@string@.startsWith('Hel')`|
|`isDateTime()`|`@string@.isDateTime()` (Should be a JSON (ISO 8601 standard) date)|
|`isEmail()`|`@string@.isEmail()`|
|`isUrl()`|`@string@.isUrl()`|
|`isEmpty()`|`@string@.isEmpty()`|
|`isNotEmpty()`|`@string@.isNotEmpty()`|
|`matchRegex($regex)`|`@string@.matchRegex('^\d+(\.\d+)?')`|
|`oneOf(...$expanders)`|`@number@.oneOf(greaterThan(10), lowerThan(0))`|

An example of chained functions:

`@string@.startsWith('You').endsWith('hello').contains('say goodbye')` --> With this pattern the string `"You say goodbye and I say hello hello, hello"` will match.

### Capturing patterns

Using the patterns explained above, you can "capture" some values to use it later. It's very useful with scenarios. 
Imagine you send a request that respond you with an ID, you want to send this ID to the next request. 
If the mathcing pattern is `@pattern@`  you will hate to write your pattern like this:

```
#var_name={{@pattern@}}
```

In the next call you will be able to inject  `var_name` like this:

```yaml
- { url: "/user?id=#var_name#", ... }
```

### Partial files

When you write arrays in your json files, it could be easier to write the content in a separate file. Here is an example of how it works:

```json
{
	"limit": 10,
	"offset": 0,
	"count": 4,
	"entities": "relative/path/to/entities"
}
```

In the file `relative/path/to/entities` you have two options. You can write a single json object, all entities will have to match this object, or you can also write an array.

### Optional fields

You can make a field optional by prefixing it by `"?"`

## Env file

The env file should be a simple key->value json file like this:

```json
{
    "key1": "value1",
    "key2": "value2"
}
```
