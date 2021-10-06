# `fac`, A Thing-Doer

`fac` is Latin for "do", and this program is a little program very similar to many others that just does things. Specifically, it reads a YAML task file for tasks and their dependencies, and then executes as many of them as it can in parallel until they're all done.

The task file will look something like this:

```yaml
Update Repo:
  command: git
  args:
    - pull
  expectedReturnCode: 0
Update JS Deps:
  command: npm
  args:
    - install
  dependencies:
    - Update Repo
Run Grunt:
  command: npx
  args:
    - grunt
  dependencies:
    - Update JS Deps
Update Gems:
  command: bundle
  args:
    - install
  dependencies:
    - Update Repo
Precompile Assets:
  command: rake
  args:
    - assets:precompile
  dependencies:
    - Update Gems
    - Run Grunt
Load DB Dump:
  command: pgrestore
  args:
    - dumps/last_dump.db
  environment:
    PGPASSWD: Swordfish
  dependencies:
    - Update Gems
  expectedStdOutRegex: successfully
Migrate DB:
  command: rake
  args:
    - db:migrate
  dependencies:
    - "! Load DB Dump"
```

In this example, the job `Update Repo` will run first all by itself (because it's the only task without dependencies). If it returns with a status code of 0 (unnecessarily specified: 0 is the default), the two tracks will proceed in parallel: a JS track and a Ruby track.

The Ruby track will attempt to migrate the database (`Migrate DB`) _only_ if the `Load DB Dump` fails, which it will do if the word `successfully` is not printed to `STDOUT`.

The `Precompile Assets` may run in parallel with the `Migrate DB` or `Load DB Dump` task, but only if both the JS track has completed successfully. So dependency diamonds are allowed (but not dependency loops).

Each task has the a name and the following anatomy:

| Field | Type | Meaning |
| ----- | ---- | ------- |
| `command` | string | The shell command to run |
| `args` | array of strings | Arguments to pass to the command |
| `environment` | dictionary of strings to strings | Environment variables to set |
| `dependencies` | array of strings | The names of other tasks that should be completed first. If the name starts with a `!` or a `-`, then the dependency is negated: the task will only run if the dependency fails. |
| `expectedReturnCode` | integer | The return code from the executable that indicates success. Defaults to 0 |
| `expectedStdOutRegex` | string | A regular expression pattern to look for in `STDOUT` that indicates success. |
| `expectedStdErrRegex` | string | A regular expression pattern to look for in `STDERR` that indicates success. |

To run the program, assuming your tasks are specified in a file called `facenda.yaml`, simple run `fac facenda.yaml`. The program uses a text-based UI (`gocui`) to display its progress. The `STDOUT` and `STDERR` of any task currently in progress is displayed in its own window. After all the tasks have been completed (successfully or not), you may use the arrow keys to scroll through them and examine their output.