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

    # get a terminal of a node:
    docker exec -it node1 bash

    # consul web interface at:
    http://192.168.99.100:8500/ui

Running outside of docker-compose (probably outdated):

.. code-block:: bash

    docker run --rm --name node1 -h node1 -v ${PWD}/data:/data svenkreiss/shoveling-worker /bin/consul agent -data-dir /data

    JOINIP="$(docker inspect -f '{{.NetworkSettings.IPAddress}}' node1)"
    docker run --rm --name node2 -h node2 -v ${PWD}/data:/data svenkreiss/shoveling-worker /bin/consul agent -data-dir /data -join $JOINIP


DNS test
--------

Running a DNS service request on the Docker host IP:

.. code-block:: bash

    dig @192.168.99.100 -p 8600 shoveling-worker.service.consul. SRV

.. code-block:: none

    dig @192.168.99.100 -p 8600 shoveling-worker.service.consul. SRV

    ; <<>> DiG 9.8.3-P1 <<>> @192.168.99.100 -p 8600 shoveling-worker.service.consul. SRV
    ; (1 server found)
    ;; global options: +cmd
    ;; Got answer:
    ;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 53458
    ;; flags: qr aa rd; QUERY: 1, ANSWER: 3, AUTHORITY: 0, ADDITIONAL: 4
    ;; WARNING: recursion requested but not available

    ;; QUESTION SECTION:
    ;shoveling-worker.service.consul. IN    SRV

    ;; ANSWER SECTION:
    shoveling-worker.service.consul. 0 IN   SRV 1 1 5060 node4.node.dc1.consul.
    shoveling-worker.service.consul. 0 IN   SRV 1 1 5060 node3.node.dc1.consul.
    shoveling-worker.service.consul. 0 IN   SRV 1 1 5060 node1.node.dc1.consul.

    ;; ADDITIONAL SECTION:
    node4.node.dc1.consul.  0   IN  A   172.18.0.5
    node3.node.dc1.consul.  0   IN  A   172.18.0.3
    node1.node.dc1.consul.  0   IN  A   172.18.0.2
    node2.node.dc1.consul.  0   IN  A   172.18.0.4

    ;; Query time: 0 msec
    ;; SERVER: 192.168.99.100#8600(192.168.99.100)
    ;; WHEN: Thu Mar 10 08:59:51 2016
    ;; MSG SIZE  rcvd: 413
