## shentud query ibc client status

Query client status

### Synopsis

Query client activity status. Any client without an 'Active' status is considered inactive

```
shentud query ibc client status [client-id] [flags]
```

### Examples

```
shentud query ibc client status [client-id]
```

### Options

```
  -h, --help   help for status
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data (default "~/.shentud")
      --log_format string   The logging format (json|plain) (default "plain")
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [shentud query ibc client](shentud_query_ibc_client.md)	 - IBC client query subcommands


