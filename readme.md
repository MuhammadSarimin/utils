# Utils



## Intro

The "utils" I’m referring to are small functions that can address minor issues we might encounter. Maybe colleagues have something to contribute. Feel free to fork the repo, and make a pull request, and share it with others.

### Utils List
#### - Unmarshal
The unmarshal function created is camelCase-sensitive, meaning it ignores JSON tags that are not in camelCase. This situation typically arises in request bodies, where using json.Unmarshal can result in fields in snake_case or lowercase being interpreted as camelCase in the struct.