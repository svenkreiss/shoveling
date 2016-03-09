shoveling
=========

* Go backend
* Shared C library that exports the Go interface using cgo (separate repo)
* Python module as one of the possible interfaces (separate repo)
* orchestration with consul: worker discovery and healthcheck through "consul services"
* configure workers with consul key/value store (e.g. max-memory)
* bulk data communication with ZMQ?


Orchestration
-------------

* workers announce their presence with hostname and port in ``consul``
* workers directly communicate to each other using the consul info


Use Cases
---------

* manager for distributed shuffle operations (e.g. for distributed sort)
    * need to be able to create a larger number of cache ids that are uniformly distributed across the cluster even though they are of size 0 when created (round robin?)
* distributed cache
    * compress data


API
---

* ``set(id, data)``
* ``get(id)``
* ``append(id, data)``
* ``init(list_ids, strategy='uniform_nodes')`` (create cache ids uniformly distributed across the cluster where strategy can be uniform in nodes, uniform in cpu or uniform in available memory)
* ``node_id(cache_id, strategy=default)`` (suggestion for best node to process this cache id)
* ``node_ids(cache_id)`` (return the nodes where the data is)
* ``_node_cache_score(node_id, cache_id)`` (a score of "distance" between node and data -- use consul's Network Coordinates)


Build
-----

.. code-block:: bash

    go build -o bin/interactive interactive/interactive.go


Build and run a worker:

.. code-block:: bash

    docker build -t svenkreiss/shoveling-worker worker/
    docker-compose up


Running outside of docker-compose (probably outdated):

.. code-block:: bash

    docker run --rm --name node1 -h node1 -v ${PWD}/data:/data svenkreiss/shoveling-worker /bin/consul agent -data-dir /data

    JOINIP="$(docker inspect -f '{{.NetworkSettings.IPAddress}}' node1)"
    docker run --rm --name node2 -h node2 -v ${PWD}/data:/data svenkreiss/shoveling-worker /bin/consul agent -data-dir /data -join $JOINIP
