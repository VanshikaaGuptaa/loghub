──────────────────────────────────────────────────────────────────────────────
LOGHUB  –  A tiny log-aggregator for local development
──────────────────────────────────────────────────────────────────────────────

Loghub watches multiple log files (or standard-input) and streams them in a
single colourised view.  It also offers live filtering and time-range export
so developers can spot errors quickly or share a short slice with a teammate.

Example live view
-----------------
![image](https://github.com/user-attachments/assets/4a9739c5-1367-4f89-ab56-d5f19308cb61)


Line format expected by Loghub
------------------------------
<ISO-timestamp>  <LEVEL>  <service>  "message text"  key=value …

Example:
2025-06-26T12:45:01.123Z  INFO  backend  "Server started"  port=8000

Build instructions
------------------
Requirements:  Go 1.24 or newer.

1.  Clone the repository and build a static binary:

      git clone https://github.com/yourname/loghub
      cd loghub
      go build -o loghub .          (Windows produces loghub.exe)

2.  Optionally move the binary to a folder on your PATH.

Quick start (basic)
-------------------
Create a folder called “logs” in your project and write log lines in the
format above, e.g.

   echo 2025-06-26T12:45:01.123Z INFO backend \"hello\" >> logs/backend.log

Then run

   loghub watch --path ./logs

Loghub colours each service and level automatically and reloads files after
log-rotate.

Quick start (reading standard-input)
------------------------------------
You can pipe any process directly into Loghub:

   npm start | loghub watch --stdin

If you want both files and stdin, provide both flags:

   npm start | loghub watch --stdin --path ./logs

Main commands
-------------
loghub watch    --path DIR             live tail of *.log in DIR
                 --services a,b        (optional) comma list of log file stems
                 --stdin               include lines from standard-input

loghub filter   --file backend.log     print only lines that match pattern
                 ERROR                 (pattern can be level or regexp)

loghub export   --path DIR             gather recent log lines into JSON
                 --since 10m           anything newer than 10 minutes
                 --out slice.json      output file

Typical integration
-------------------
 • Node / Express backend – use Winston and point its file transport at
   ../logs/backend.log

 • React / Vite frontend – pipe Vite through a tiny node script that prefixes
   each stdout line, e.g.:
       node tools/prefix-log.mjs frontend vite dev

 • Background workers – redirect stdout/stderr to logs/worker.log

Run everything with a process manager such as “concurrently”:

   concurrently -k "
     npm --prefix backend run dev" "
     npm --prefix frontend run dev" "
     ./loghub watch --path ./logs --services frontend,backend,worker"

Useful patterns
---------------
View only errors in real time
   loghub watch --path ./logs | loghub filter ERROR

Save last five minutes for a bug report
   loghub export --path ./logs --since 5m --out bug-slice.json

Build static binary for all platforms (requires GoReleaser installed)
   goreleaser release --snapshot --clean


