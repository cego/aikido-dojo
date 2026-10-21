# Support for Dependency Scanning by Language (SCA)

Aikido performs SCA on every push and pull request, in real time in your IDE, and as a scheduled re-scan (default daily) so that newly disclosed CVEs in unchanged dependencies are still caught. SCA covers known CVEs and risky open-source licenses.

Below is a table of supported languages and their respective lockfiles. We recommend using lockfiles by default as they increase speed at build time, make your builds more reproducible and they are a first layer of defense against supply-chain attacks. Of course, lockfiles also help Aikido in finding vulnerable packages.

We scan for lockfiles both in the root as in all subfolders.

| Language                        | Lockfiles scanned                                                                                                                                                     |
| ------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| <p>JavaScript<br>TypeScript</p> | <p>npm-shrinkwrap.json</p><p>package-lock.json</p><p>yarn.lock</p><p>pnpm-lock.yaml</p><p>pnpm-lock.yml</p><p>bun.lock</p><p>deno.lock</p><p>libman.json</p>          |
| PHP                             | composer.lock                                                                                                                                                         |
| Java                            | <p>gradle.lockfile</p><p>build.gradle</p><p>pom.xml</p><p>.jar</p><p>.war</p><p>.ear<br>ivy.xml</p>                                                                   |
| Swift                           | <p>Package.resolved</p><p>Podfile.lock</p>                                                                                                                            |
| Go                              | go.mod                                                                                                                                                                |
| Python                          | <p>Pipfile.lock</p><p>poetry.lock</p><p>uv.lock</p><p>pdm.lock</p><p>requirements.txt</p><p>requirements.lock</p><p>Conda: requirements.yml</p>                       |
| .NET                            | <p>.csproj</p><p>.deps.json</p><p>packages.lock.json</p><p>packages.config</p><p>Packages.props<br>paket.lock</p>                                                     |
| Ruby                            | gemfile.lock                                                                                                                                                          |
| Rust                            | <p>cargo.lock</p><p>cargo.toml</p>                                                                                                                                    |
| Kotlin                          | build.gradle.kts, gradle.lockfile                                                                                                                                     |
| Dart                            | pubspec.lock                                                                                                                                                          |
| Elixir                          | mix.lock                                                                                                                                                              |
| C/C++                           | <p>conan.lock</p><p>vcpkg.json</p><p>Lockfileless C/C++ dependencies (<a href="https://help.aikido.dev/doc/cc-lockfile-less-scanning/doczSgARAloY">more info</a>)</p> |
| Scala                           | <p>build.sbt</p><p>plugins.sbt</p><p>dependencies.scala</p><p>libraries.scala</p><p>.sbt.lock</p>                                                                     |
| Clojure                         | <p>deps.edn<br>project.clj</p>                                                                                                                                        |
| Unity UPM                       | packages-lock.json (Aikido only scans and imports UPM packages fetched from NPM)                                                                                      |
