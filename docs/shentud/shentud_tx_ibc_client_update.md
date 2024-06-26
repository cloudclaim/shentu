## shentud tx ibc client update

update existing client with a header

### Synopsis

update existing client with a header

```
shentud tx ibc client update [client-id] [path/to/header.json] [flags]
```

### Examples

```
shentud tx ibc client update [client-id] [path/to/header.json] --from node0 --home ../node0/<app>cli --chain-id $CID
```

### Options

```
  -h, --help   help for update
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

* [shentud tx ibc client](shentud_tx_ibc_client.md)	 - IBC client transaction subcommands


