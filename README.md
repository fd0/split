Split large files into smaller ones using the same Content Defined Chunking
algorithm the [restic][1] backup program uses.

Build (using Go >= 1.11):

    $ go build

Sample usage:

    $ ./split -v -o /tmp < /tmp/data

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
