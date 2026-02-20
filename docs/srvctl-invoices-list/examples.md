A command to list all invoices for the account:

```
srvctl invoices list
```

A command to filter invoices by status:

```
srvctl invoices list --status paid
```

A command to filter invoices by type:

```
srvctl invoices list --type invoice
```

A command to filter invoices by currency:

```
srvctl invoices list --currency USD
```

A command to filter invoices by date range:

```
srvctl invoices list --start-date 2024-01-01 --end-date 2024-12-31
```
