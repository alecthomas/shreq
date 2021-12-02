# Find external commands required by shell scripts

When writing shell scripts that need to run portably across multiple hosts and
platforms, it's useful to know what external commands are required. This utility
does exactly that.

## Example

```
$ shreq ./testdata/*.sh
testdata/001.sh:3:20: unsupported external command: parallel
testdata/002.sh:3:1: unsupported external command: parallel
testdata/002.sh:5:3: unsupported external command: mc
$ shreq -c parallel,mc testdata/*
$
```

## Usage

```text
Usage: shreq <script> ...

This utility verifies all commands used by a shell script against an allow list:

      . : [ admin alias ar asa at awk basename bash batch bc bg bind break
      builtin c99 cal case cat cd cflow chgrp chmod chown cksum cmp comm
      command compgen complete compress continue cp crontab csplit ctags cut
      cxref date dd declare delta df diff dirname dirs disown du echo ed
      enable env eval ex exec exit expand export expr false fc fg file find
      fold fort77 fuser gencat get getconf getopts grep hash head help
      history iconv id if ipcrm ipcs jobs join kill let lex link ln local
      locale localedef logger logname logout lp ls m4 mailx make man mesg
      mkdir mkfifo more mv newgrp nice nl nm nohup od paste patch pathchk
      pax popd pr printf prs ps pushd pwd qalter qdel qhold qmove qmsg
      qrerun qrls qselect qsig qstat qsub read readonly renice return rm
      rmdel rmdir sact sccs sed set sh shift shopt sleep sort source split
      strings strip stty suspend tabs tail talk tee test time times touch
      tput tr trap true tsort tty type typeset ulimit umask unalias uname
      uncompress unexpand unget uniq unlink unset until uucp uudecode
      uuencode uustat uux val vi wait wc what while who write xargs yacc

Arguments:
  <script> ...    Shell scripts to validate.

Flags:
  -h, --help              Show context-sensitive help.
  -a, --allow=none,...    Enable optional features (none,relative,var-relative).
  -c, --cmds=CMD,...      Extra commands to allow.
```
