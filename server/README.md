there are at least 3 data types saved in the redis db:

- hashes: store blob information
- sortedSet: stores blob's term space (list of all words)
- set: list of all files uuid (which represents their names in the data/)
