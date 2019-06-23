Split large files into smaller ones using the same Content Defined Chunking
algorithm the [restic][1] backup program uses.

Build (using Go >= 1.11):

    $ go build

Sample usage:

    $ ./split --verbose --output /tmp --input /tmp/data
    next chunk offset 0, 814244 bytes written to /tmp/split-000
    next chunk offset 814244, 1649886 bytes written to /tmp/split-001
    next chunk offset 2464130, 3332485 bytes written to /tmp/split-002
    next chunk offset 5796615, 1996103 bytes written to /tmp/split-003
    [...]
    next chunk offset 101940538, 700441 bytes written to /tmp/split-069
    next chunk offset 102640979, 533829 bytes written to /tmp/split-070
    next chunk offset 103174808, 537761 bytes written to /tmp/split-071
    next chunk offset 103712569, 1145031 bytes written to /tmp/split-072
    wrote 104857600 bytes to 73 files

Using `cat`, we can put the file back together (and use `sh256sum` to verify it's the same data):

    $ cat /tmp/split-* | sha256sum
    66b9d2c5de34170b93f387988a80fb600717da5e437dcc4da1025343fb9019a1  -

    $ sha256sum /tmp/data
    66b9d2c5de34170b93f387988a80fb600717da5e437dcc4da1025343fb9019a1  /tmp/data

Check out the help for other options:

    $ ./split -h
    Usage of split:
      -i, --input file     Read from file instead of stdin
      -u, --max-size n     Set maximal chunk size to n bytes (default 8388608)
      -l, --min-size n     Set minimal chunk size to n bytes (default 524288)
      -o, --output dir     Write files to directory dir instead of the current directory (default ".")
      -p, --polynomial p   Use polynomial p for splitting (hex notation, no prefix) (default "3DA3358B4DC173")
      -t, --template s     Use s as the (printf-style) template for output files (default "split-%03d")
      -v, --verbose        Be verbose

The library used for this program is https://github.com/restic/chunker

If you're interested in the mathematical foundation for Content Defined Chunking with [Rabin Fingerprints][2], head over to [the restic blog][3] which has an [introductory article][4].

[1]: https://restic.net
[2]: https://en.wikipedia.org/wiki/Rabin_fingerprint
[3]: https://restic.net/blog
[4]: https://restic.net/blog/2015-09-12/restic-foundation1-cdc
