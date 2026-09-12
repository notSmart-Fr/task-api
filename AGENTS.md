# Agent Instructions

## Secrets and Environment Files

- Do not read `.env` or any other file containing runtime secrets.
- Use `.env.example`, documented configuration, or source references to determine required environment variable names.
- Do not print, log, copy, or include secret values in tool output, patches, tests, chat responses, or generated documentation.
- When configuration verification is necessary, check only whether an environment variable is set; never retrieve or display its value.
- Treat any secret value that is accidentally exposed as compromised and tell the user to rotate it.
