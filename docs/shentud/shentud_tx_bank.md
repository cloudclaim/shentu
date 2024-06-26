## shentud tx bank

Bank transaction subcommands

```
shentud tx bank [flags]
```

### Options

```
  -h, --help   help for bank
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

* [shentud tx](shentud_tx.md)	 - Transactions subcommands
* [shentud tx bank locked-send](shentud_tx_bank_locked-send.md)	 - Send coins and have them locked (vesting).
* [shentud tx bank send](shentud_tx_bank_send.md)	 - Send funds from one account to another. 
		Note, the'--from' flag is ignored as it is implied from [from_key_or_address].
		When using '--dry-run' a key name cannot be used, only a bech32 address.


