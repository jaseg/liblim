liblim
======

``liblim`` is a microframework for database synchronization in Python and Go based on Convergent Replicated Data Types
(CRDTs). CRDTs are a technique for constructing well-defined merge semantics for diverging states from different parts
of a network that are guaranteed to succeed without conflicts.

In terms of the CAP-Theorem_, a ``liblim``-based system is highly available and partition-tolerant, but not guaranteed
to be consistent at any time. ``liblim`` can only guarantee eventual consistency with an decay factor that is dependent
on the concrete network conditions. For a ``liblim``-based system, however, this limited consistency is guaranteed to
not affect the execution of operations. Any operation performed on a ``liblim``-based system is guaranteed to succeed
and under no conditions any data is lost or operation is dropped due to version conflicts.

Use cases
---------
``liblim`` can be used in any system requiring multi-party synchronization. This does apply to a classical server-client
architecture as well as a fully peer-to-peer system. Your author thinks ``liblim`` is especially suited for database
synchronization for cross-device applications since it may use, but is not dependent on, a centralized server.

Protocol
--------
At the moment, this prototype is using a simplistic text-based protocol on top of TCP. We are working on fully
specifying this protocol and will port it to different transports as soon as our prototype implementations yields first
results.

Compatibility
-------------
Two implementations of liblim, one in Python and one in Go, are being developed in parralel. Both are tested against
each other to ensure compatibility.

Current status
--------------
This project is under active development by Matti Möll and Sebastian Götte. Any feedback is appreciated and you may feel
free to scrutinize and improve our work. Currently, this is not even alpha-quality software. Though we very much
appreciate any input, if your goal is to just use liblim in your application, please bear with us some more until we
have finished our first prototype work. As soon as the protocol reaches a somewhat stable state and the two
implementations are working, we will release a preliminary user guide for those interested in experimenting with this.

.. _CAP-Theorem: http://en.wikipedia.org/wiki/CAP_theorem
