### The Artificer

Sometimes you know, you don't like shell scripts for such trivial tasks right? This tool is your recursive collector for lambda executables and archives.
It is 100% maintainable and readable. Open to PRs.

Thought this may be useful to have this tool due to GoReleaser requiring Go entrypoints for each lambda.
And AWS requires them to be zipped under our artifact folder.

The default folder structure it understands;
```
\project
-> \artifacts

-> \infrastucture
--> \stacks
---> \integrationName1
----> \prd
----> \uat
----> \dev
----> terragrunt.hcl
----> vars.tf
----> \lambda
-----> \makeHappy
------> \main.go
---> \integrationName2
----> \prd
----> \uat
----> \dev
----> terragrunt.hcl
----> vars.tf
----> \lambda
-----> \doSomething
------> \main.go

-> \tools
--> \artificer
---> \main.go
```

You have to provide 2 arguments for it, which is -artifacts_path and -stacks_path. If not given,
it will assume stacks path as ../../infrastructure/stacks and artifacts path as ../../artifacts.

 ```
 _______  _______ __________________ _______ _________ _______  _______  _______ 
(  ___  )(  ____ )\__   __/\__   __/(  ____ \\__   __/(  ____ \(  ____ \(  ____ )
| (   ) || (    )|   ) (      ) (   | (    \/   ) (   | (    \/| (    \/| (    )|
| (___) || (____)|   | |      | |   | (__       | |   | |      | (__    | (____)|
|  ___  ||     __)   | |      | |   |  __)      | |   | |      |  __)   |     __)
| (   ) || (\ (      | |      | |   | (         | |   | |      | (      | (\ (   
| )   ( || ) \ \__   | |   ___) (___| )      ___) (___| (____/\| (____/\| ) \ \__
|/     \||/   \__/   )_(   \_______/|/       \_______/(_______/(_______/|/   \__/
```

### Usage

```sh
cd tools/artificer
go run main.go
```

So, for the custom directories of artifacts and stacks, use flags
```sh
go run tools/artificer/main.go -artifacts_path=artifacts -stacks_path=infrastructure/stacks  
```

Have a look into your artifacts folder now!