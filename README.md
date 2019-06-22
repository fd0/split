Split large files into smaller ones using the same Content Defined Chunking[1]
algorithm the restic[2] backup program uses.

Build (using Go >= 1.11):

    $ go build

Sample usage:

    $ ./split -v -o /tmp < /tmp/data

The library used for this program is https://github.com/restic/chunker

[1] https://restic.net
[2] https://restic.net
