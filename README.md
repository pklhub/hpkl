# HPKL

## Overview

HPKL is an extension of the `pkl-lang` that introduces:

1. **OCI Packages Support**: Simplified management and integration of OCI-compliant packages.
2. **Secrets Reading as Resources**: Streamlined handling of secrets as first-class resources within your configurations.

---

## Installation

### macOS
```bash
brew tap hpklio/hpkl
brew install hpkl
```

### Linux (amd64)
```bash
curl -L -o hpkl.tar.gz https://github.com/hpklio/hpkl/releases/download/v0.8.0/hpkl_Linux_x86_64.tar.gz
tar -xzvf hpkl.tar.gz
chmod +x hpkl
```

### Linux (arm64)
```bash
curl -L -o hpkl.tar.gz https://github.com/hpklio/hpkl/releases/download/v0.8.0/hpkl_Linux_arm64.tar.gz
tar -xzvf hpkl.tar.gz
chmod +x hpkl
```

### Windows
```bash
curl -L -o hpkl.tar.gz https://github.com/hpklio/hpkl/releases/download/v0.8.0/hpkl_Windows_x86_64.zip
tar -xzvf hpkl.tar.gz
chmod +x hpkl
```

---

## Usage

### Adding OCI Packages
To include an `.oci` package, simply add the `.oci` suffix to your dependency declaration. For example:

```pkl
dependencies {
  ["hpkl-k8s-app.oci"] { uri = "package://ghcr.io/hpklio/hpkl-k8s-app@0.20.0" }
}
```

### Reading Secrets
Use the `read?()` function to read secrets. For example:

```pkl
read?("vals:sops:file.json#/path").text
```

This allows you to securely fetch and use secrets in your configurations.

### Detailed Example
For a more comprehensive example, check out the [HPKL Kubernetes App Example](https://github.com/hpklio/hpkl-k8s-app/blob/main/vals.pkl).

---

## Contributing
We welcome contributions! To get started:

1. Fork this repository.
2. Create a new branch for your feature or bug fix.
3. Commit your changes with a clear message.
4. Submit a pull request for review.

