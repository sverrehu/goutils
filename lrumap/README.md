# lrumap

Implements a map with a maximum size and a time-to-live for each entry.
This map is suitable as a cache.
When a new entry is added that would exceed the maximum size of the cache, either the least recently used (LRU) item is evicted, or an already expired item is evicted, whichever is found first.
