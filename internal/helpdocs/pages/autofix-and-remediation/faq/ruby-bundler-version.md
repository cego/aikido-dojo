# Ruby / Bundler version

AutoFix uses the Ruby & Bundler versions defined in the Gemfile. If these versions are not present, AutoFix will default to Ruby 3.2.9 and Bundler 2.4.19.

Example of a minimal Gemfile with the versions defined:

{% code title="Gemfile" %}

```
source 'https://rubygems.org'

ruby '4.0.6'
gem 'bundler', '4.0.3'
```

{% endcode %}
