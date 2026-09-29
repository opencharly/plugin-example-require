# plugin-example-require

The reference plugin for charly's **declared inter-plugin dependency** mechanism
(`plugin.requires:`) — a plugin that declares the OTHER plugins it needs and
peer-invokes them over the host broker at runtime.

Its candy manifest declares two peers:

- `verb:externalprobe` — served by `candy/plugin-example-external`, compiled into
  charly (the *already-registered* / resolve branch).
- `verb:exampledispatchpeer` — served out-of-process by
  `candy/plugin-example-dispatch`, connected on demand from its declared `source`
  (the *connect* branch).

The plugin's `examplerequire` verb invokes both peers via
`sdk.Executor.InvokeProvider` and echoes their results.

## What it provides

| Capability | Surface |
|---|---|
| `verb:examplerequire` | the `examplerequire:` check verb — peer-invokes the declared `builtin_peer` and `external_peer` and echoes their results |

The authored `plugin.requires:` list in `candy/plugin-example-require/charly.yml`
is THE declaration. The host reads it from the resolved view and resolves every
requirement before this plugin's providers register, so the project's own plans
need not reference a plugin's internal peers. There is deliberately **no**
`sdk.NewMetaWithRequires` call in this plugin's Go — the manifest is the one
source, so declaring another peer needs zero core change and no Go change.

## How to use it

Compose the plugin candy in a box or check bed's `candy:` list:

```yaml
- '@github.com/opencharly/plugin-example-require/candy/plugin-example-require:<tag>'
```

Then author the verb in a plan; the peer words arrive as data:

```yaml
- check: the examplerequire verb peer-invokes both declared peers
  id: examplerequire-dispatches
  examplerequire:
    marker: plugin-example-require-ok
    builtin_peer: externalprobe
    external_peer: exampledispatchpeer
  context: [runtime]
```

## Layout

- `candy/plugin-example-require/` — the plugin module: `plugin.go` (the provider +
  `NewProvider()`/`NewMeta()`), `schema/examplerequire.cue` (the self-contained
  `#ExamplerequireInput`), `cmd/serve/main.go` (the out-of-process shim).
- `charly.yml` — the root project manifest (`discover: candy`).
- `.github/workflows/tag-on-merge.yml` — CalVer tag + `CHANGELOG/` on merge.

## Related

- Owning skill: `/charly-internals:plugin` — the plugin/provider model. This candy
  carries no `skill:` entity of its own; the gap is tracked in
  [opencharly/opencharly#291](https://github.com/opencharly/opencharly/issues/291).
- [`opencharly/charly`](https://github.com/opencharly/charly) — the charly CLI.
