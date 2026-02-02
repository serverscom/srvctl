Load completions in current bash session:

```
source <(srvctl completions bash)
```

Install completions permanently on Linux (system-wide):

```
srvctl completions bash > /etc/bash_completion.d/srvctl
```

Install completions permanently on Linux (user-level):

```
srvctl completions bash > ~/.bash_completion
```

Install completions permanently on macOS with Homebrew:

```
srvctl completions bash > $(brew --prefix)/etc/bash_completion.d/srvctl
```

Install completions for zsh (add to ~/.zshrc):

```
source <(srvctl completions zsh)
```

Install completions for fish:

```
srvctl completions fish > ~/.config/fish/completions/srvctl.fish
```

Install completions for PowerShell:

```
srvctl completions powershell | Out-String | Invoke-Expression
```
