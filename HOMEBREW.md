# Homebrew

LangMate includes a Homebrew formula template at `Formula/langmate.rb`.

## Local Tap Test

Modern Homebrew expects formulae to live in a tap. To test the formula locally:

```bash
brew tap-new local/langmate
cp Formula/langmate.rb "$(brew --repository local/langmate)/Formula/langmate.rb"
brew install --HEAD --build-from-source local/langmate/langmate
```

Install the app bundle into `/Applications`:

```bash
langmate-install-app
```

## Publishing

1. Create and push a release tag:

   ```bash
   git tag v1.0.3
   git push origin v1.0.3
   ```

2. Download the release tarball and calculate its checksum:

   ```bash
   curl -L -o langmate-1.0.3.tar.gz \
     https://github.com/niuguy/langmate/archive/refs/tags/v1.0.3.tar.gz
   shasum -a 256 langmate-1.0.3.tar.gz
   ```

3. Replace the all-zero `sha256` in `Formula/langmate.rb`.

4. Publish the formula in a tap repository, for example `homebrew-langmate`.

5. Users install with:

   ```bash
   brew tap niuguy/langmate
   brew install langmate
   langmate-install-app
   ```

Homebrew installs the app bundle under its prefix. `langmate-install-app` copies that bundle into `/Applications`, quits any running LangMate instance, and relaunches the app.
