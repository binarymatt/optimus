definitions:
tracker - the tracker keeps track of files being tailed via a single fsnotify watcher.
Primary Workflow:

1. initialize tracker
2. add path(s)
3. watch for fsnotify events
4. exit when context is cancelled.

Directory Path Workflow:

1. add path to watcher

file path workflow:

1. add path to watcher
2. initialize tailer struct
3. start tailer

event workflow:

- if event is create, add new file path.
- if event is write, stop waiting
- if event is delete, stop associated tailer (if it was file)

tail workflow:

- read all lines from saved position
- if eof, wait for changes
