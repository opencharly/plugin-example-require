# plugin-example-require

The `plugin-example-require` reference plugin of the
[opencharly/charly](https://github.com/opencharly/charly) candy library, as a
standalone repo.

It is the live witness for charly's **declared inter-plugin dependency** mechanism
(`plugin.requires:`). Its candy manifest declares two OTHER plugins it depends on —
one that charly compiles in (`verb:externalprobe`) and one served out-of-process
(`verb:exampledispatchpeer`) — and its `examplerequire` verb peer-invokes both over
the host broker (`sdk.Executor.InvokeProvider`).

The authored `plugin.requires:` list in `candy/plugin-example-require/charly.yml` is
THE declaration; the host reads it from the resolved view and resolves every
requirement before the plugin's providers register. There is deliberately **no**
`sdk.NewMetaWithRequires` call in this plugin's Go — the manifest is the one source,
so a plugin author declares a peer ONCE and neither the core nor another plugin's Go
changes.

The Go module lives at `candy/plugin-example-require/` with module path
`github.com/opencharly/plugin-example-require/candy/plugin-example-require`; the
charly resolver fetches this repo at the pinned tag and the compiled-in wiring
imports the module at that path.
