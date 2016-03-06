shoveling
=========

* Go backend
* Python one of the possible frontends
* orchestration with consul: worker discovery through "consul services"
* configure workers with consul key/value store (e.g. max-memory)
* "distributed namenode" using consul key/value store, or updates propagate too slow?
* bulk data communication with ZMQ?


Orchestration
-------------

* workers announce their presence with hostname and port in ``consul``
* workers directly communicate to each other using the consul info


Use Cases
---------

* manager for distributed shuffle operations (e.g. for distributed sort)
* distributed cache


API
---

* ``set(id, data)``
* ``get(id)``
* ``node_id(cache_id, strategy=default)`` (suggestion for best node to process this cache id)
* ``_node_cache_score(node_id, cache_id)`` (a score of "distance" between node and data -- use consul's Network Coordinates)
