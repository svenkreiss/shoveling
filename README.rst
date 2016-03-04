shoveling
=========

* C backend
* Python one of the possible frontends
* orchestration with consul
* communication with ZMQ?


Orchestration
-------------

* workers announce their presence with hostname and port in ``consul``
* workers directly communicate to each other using the consul info


Use Cases
---------

* manager for distributed shuffle operations (e.g. for distributed sort)
* distributed cache
