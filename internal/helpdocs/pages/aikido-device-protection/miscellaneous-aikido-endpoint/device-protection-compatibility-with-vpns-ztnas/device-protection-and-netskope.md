# Device Protection & Netskope for Mac

This guide explains how to configure your Netskope Steering Configuration so that Aikido Device Protection works reliably alongside Netskope.

{% hint style="warning" %}
These instructions are for macOS only! Applying these for Windows can cause issues.
{% endhint %}

## Why This Is Needed

Device Protection inspects and protects software supply chain traffic, the connections your developers make to package registries (npm, PyPI, Maven, NuGet, Go) and developer tool marketplaces (Visual Studio, Cursor, Chrome Web Store), as well as Aikido's own protection services.

When both a corporate VPN/ZTNA and Aikido Device Protection attempt to steer or inspect the same traffic, the two can conflict, leading to unreliable behavior. To avoid this, you should configure VPN/ZTNA to **bypass (exclude) the domains that Device Protection handles**.

## Steps

1. Log in to the **Netskope Admin Console**.
2. Navigate to **Settings → Security Cloud Platform → Steering Configuration.**
3. Click New Steering Configuration (do not edit the Default tenant config) and assign a new name e.g. "macOS tenant config".
4. Under **Match Criteria**, set **OS Family** = **macOS**. Leave **User Group/OU** set to **Any**, or scope it to the same group your Default config covers. This is what restricts the bypass to Macs.

   If you don't see **Match Criteria**, your tenant is on a release older than 124.0.0. Scope the configuration to a **User Group** or **Organizational Unit (OU)** instead. This scoping applies per user, not per device, so the group must only contain users whose devices are all Macs. Otherwise their Windows devices get the bypass too. To create a new group, add it in your identity provider and sync it to Netskope (SCIM or Directory Importer) before selecting it here.
5. Under Steered Traffic, mirror your Default config so Mac steering is unchanged. This typically includes All Traffic (HTTP/HTTPS and Non-web) and All Private App Segments.
6. Click the 3 dots for settings on the relevant Steering Configuration (Config) profile you want to modify. Open Edit Configuration → Traffic Steering and make sure "**Bypass exception traffic at:**" is set to **Client**.
7. After you've checked this, enter the profile.
8. Open the **Exceptions** tab.
9. Click **New Exception → Domains**.
10. Enter the domains listed below (one per line, or as supported by your version of the interface).
11. Save the configuration. Changes typically propagate in **1 hour**.

{% hint style="info" %}
If you want to do a quick test, go to the local Netskope client, click Configuration, then check for "Configuration update available." If it's there, click Update so it fetches the latest policies.
{% endhint %}

## Domains to Exclude

These are the domains Aikido Device Protection intercepts. Package registries and developer tool marketplaces are always intercepted. AI provider domains are only intercepted when [AI Usage](/aikido-device-protection/using-aikido-endpoint/ai-monitoring.md) monitoring is enabled.

{% tabs %}
{% tab title="Package registries & marketplaces" %}

```
api.nuget.org
app.aikido.dev
endpoint-server.aikido.dev
aikido-endpoint-binaries.s3.eu-west-1.amazonaws.com
central.maven.org
chromewebstore.google.com
chromewebstore.googleapis.com
clients2.google.com
clients2.googleusercontent.com
crates.io
index.crates.io
static.crates.io
device-protection.aikido.help
files.pythonhosted.org
gallery.vsassets.io
gallerycdn.vsassets.io
github.com
globalcdn.nuget.org
marketplace.cursorapi.com
marketplace.visualstudio.com
open-vsx.org
proxy.golang.org
pypi.org
pypi.python.org
registry.npmjs.com
registry.npmjs.org
registry.yarnpkg.com
repo.maven.apache.org
repo1.maven.org
repository.apache.org
rubygems.org
index.rubygems.org
repo.packagist.org
sum.golang.org
update.googleapis.com
www.nuget.org
```

{% endtab %}

{% tab title="AI providers" %}
Only intercepted when AI Usage monitoring is enabled.

```
ai-gateway.helicone.ai
ai-gateway.vercel.sh
ai.zenifra.com
aki.io
api-inference.modelscope.cn
api-sherlock.cloudferro.com
api.302.ai
api.abliteration.ai
api.ai-router.dev
api.aiand.com
api.ambient.xyz
api.anyapi.ai
api.auriko.ai
api.berget.ai
api.cerebras.ai
api.clarifai.com
api.claudin.io
api.cline.bot
api.code.umans.ai
api.cortecs.ai
api.crossmodel.ai
api.deepinfra.com
api.deepseek.com
api.dinference.com
api.empiriolabs.ai
api.fireworks.ai
api.friendli.ai
api.getlilac.com
api.githubcopilot.com
api.gmi-serving.com
api.greenpt.ai
api.groq.com
api.hpc-ai.com
api.impossibl.com
api.inceptionlabs.ai
api.inceptron.io
api.inference.wandb.ai
api.intelligence.io.solutions
api.jiekou.ai
api.kilo.ai
api.lkeap.cloud.tencent.com
api.llama.com
api.llmgateway.io
api.longcat.chat
api.lucidquery.com
api.meganova.ai
api.meta.ai
api.mistral.ai
api.modeloracle.com
api.moonshot.ai
api.moonshot.cn
api.morphllm.com
api.neuralwatt.com
api.nova.amazon.com
api.novita.ai
api.ofox.ai
api.openai-compat.model-serving.eu01.onstackit.cloud
api.orcarouter.ai
api.perplexity.ai
api.pioneer.ai
api.poe.com
api.qhaigc.net
api.qnaigc.com
api.regolo.ai
api.routing.run
api.sakana.ai
api.sarvam.ai
api.scaleway.ai
api.scx.ai
api.siliconflow.cn
api.siliconflow.com
api.stepfun.ai
api.stepfun.com
api.synthetic.new
api.tbox.cn
api.tensorx.ai
api.thegrid.ai
api.together.xyz
api.tokenfactory.nebius.com
api.trustedrouter.com
api.unorouter.com
api.upstage.ai
api.vivgrid.com
api.vultrinference.com
api.xiaomimimo.com
api.z.ai
api.zeldoc.ai
apis.iflow.cn
app.frogbot.ai
chat.d.run
cloud-api.near.ai
coding-intl.dashscope.aliyuncs.com
coding-plan-endpoint.kuaecloud.net
coding.dashscope.aliyuncs.com
crof.ai
daoxe.com
dashscope-intl.aliyuncs.com
dashscope.aliyuncs.com
go.fastrouter.ai
hyper.charm.land
inference.baseten.co
inference.coralbricks.ai
inference.do-ai.run
inference.hetzner.com
inference.net
inference.poolside.ai
inference.tinfoil.sh
inference.us-west.modal.direct
integrate.api.nvidia.com
kenari.id
llm.chutes.ai
llm.submodel.ai
llmtr.com
maas-api.ebcloud.com
moark.com
model.inferx.net
modelishub.com
models.mixlayer.ai
models.think.evroc.com
nano-gpt.com
oai.endpoints.kepler.ai.cloud.ovh.net
ollama.com
open.bigmodel.cn
openai.blueclaw.network
opencode.ai
openrouter.ai
pass.wafer.ai
routellm.abacus.ai
router.huggingface.co
router.requesty.ai
token-plan-ams.xiaomimimo.com
token-plan-cn.xiaomimimo.com
token-plan-sgp.xiaomimimo.com
token-plan.ap-southeast-1.maas.aliyuncs.com
token-plan.cn-beijing.maas.aliyuncs.com
tokenhub.tencentmaas.com
www.xpersona.co
zenmux.ai
```

{% endtab %}
{% endtabs %}
