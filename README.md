# go-bank
An experimental project that mimics a small banking service with a go server w/ Postgres DB containerized by docker

### Dirty Write
This occurs when two or more transactions try to overwrite each others' modifications over same resource.
For example, in the case for a transaction in which a thread reads then writes onto account table, a thread could reads the data while the other thread reads and then writes the very same data. The problem is that the time the value read by the first thread is stays behind the modification made by the second thread and when the first thread writes the data, it overwrites the change.
I fixed the problem by using FOR UPDATE for SELECT query to obtain a lock to prevent other thread to read the very data.

### Deadlock caused by foreign key constaints
A dead lock could caused from many reasons and one of them occurs when a thread creates an entry row associated with an account id as a foreign key while the other thread is creating/updating the same account row.
For this case, I added NO KEY to the SELECT FOR UPDATE query not to prevent from other thread to be blocked
