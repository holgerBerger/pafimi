# pafimi
WORK IN PROGRESS - NOT USABLE 

parallel file migration

tool to migrate many large files in parallel from one parallel filesystem to
another, using a pool of worker machines having both filesystems mounted.

architecture:

a client submits a copy job to a initial worker.
A copy job is defined by a source and target directory.
The initial worker is picked from a pool of workers at random,
and does some of the work and the orchestration of the work pool.

Idea a)
The initial worker is traversing the source tree, and creates the directory
structure in the target filesystem as he traverses the filesystem. 
In directories with sufficient files, he distributes the copying of the 
files to the other workers, and continues to traverse deeper into directories.


Idea b)
The initial worker is traversing all the source tree, and searches for subtrees
of equal size and complexity which are delegated to the workers.

