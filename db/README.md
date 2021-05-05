## Dirty Write
This occurs when two or more transactions try to overwrite each others' modifications over same resource.
For example, in the case for a transaction in which a thread reads then writes onto account table, a thread could reads the data while the other thread reads and then writes the very same data. The problem is that the time the value read by the first thread is stays behind the modification made by the second thread and when the first thread writes the data, it overwrites the change.
we can solve the problem by using FOR UPDATE for SELECT query to obtain a lock to prevent other thread to read the very data.

## Deadlock
### Caused by foreign key constaints

A dead lock could caused from many reasons and one of them occurs when a thread creates an entry row associated with an account id as a foreign key while the other thread is creating/updating the same account row.

_Example_:

Thread 1

Create Transfer row for fromAccount and toAccount
Create Entry row for fromAccount -> waiting for thread 2 bc thread 2 is trying to update fromAccount
Create Entry row for toAccount
Update Account of fromAccount 
Update Account of toAccount

Thread 2

Create Transfer row for fromAccount and toAccount
Create Entry row for fromAccount
Create Entry row for toAccount
Update Account of fromAccount -> waiting for thread 1 bc thread 1 is trying to create an entry for fromAccount
Update Account of toAccount

For this case, we could add NO KEY to the SELECT FOR UPDATE query not to prevent from other thread to be blocked

### Inconsistent Update Order

A dead lock could also occur when the order of write is not consistent.

_Example_:

Thread 1

Update Account of fromAccount 
Update Account of toAccount -> blocked bc thread 2 is trying to update toAccount

Thread 2

Update Account of toAccount 
Update Account of fromAccount -> blocked bc thread 1 is trying to update fromAccount

As a solution, we could enforce the order of updating the account. For example, we can update the account with ID smaller than the other so that the order of the update is guaranteed.


