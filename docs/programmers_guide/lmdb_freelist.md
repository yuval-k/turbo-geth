### LMDB Freelist

Copy-on-write - means, when you update 1 record on Page - LMDB does copy whole page. Old page marked as "free", can be
re-used later.
`mp->mp_flags & P_DIRTY`

Every TX has auto-incremental ID. If ReadOnly-Tx with ID=1 still in-progress, but Write-Tx with ID=2 finished and
produced many free pages - then all that free pages are not re-usable until tx=1 finish (because tx=1 must see world
same as it was when tx=1 started).
`mdb_find_oldest` func returns smallest tx id which still running.

FREE_DBI=0 - table where LMDB stores overall history - which tx freed which pages:
[tx_id] -> [number_of_page_ids_u64][list_of_page_ids]
Can read FREE_DBI by `mdb_cursor_open(tx, 0)` - it works only in Read transactions:
@see `mdb_cursor_open`: `if (dbi == FREE_DBI && !F_ISSET(txn->mt_flags, MDB_TXN_RDONLY))`

All pages which freed by Tx stored in `txn->mt_free_pgs`
Inside every `tx.Commit()` called method `mdb_freelist_save` which stores `txn->mt_free_pgs` to FREE_DBI.
@see `mdb_cursor_put(&mc, &key, &data, MDB_RESERVE);`

1 value can be larger than 1 page, it is good-practice to hold it on sequence of pages (avoid memory fragmentation).
`mdb_page_alloc` - is main place to hold this logic.

Page loose: If page has been pulled from the FreeDBI, but has been deleted (during same Tx) - it placed
into `txn->mt_loose_pgs`. Use these pages first before pulling again from the FreeDBI. @see `mdb_page_loose`

```
if ((mp->mp_flags & P_DIRTY) && mc->mc_dbi != FREE_DBI) {    loose = 1;    }
```


