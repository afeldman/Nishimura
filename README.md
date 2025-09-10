# Nishimura

Nishimura is a package manager for FANUC Karel projects.
It is named after Makoto Nishimura, the engineer of the first Japanese robot Gakutensoku (1928).

With Nishimura you can:

- Manage project manifests (nishimura.kpc)
- Add, install and resolve dependencies
- Check for conflicts
- Compile Karel code via the Gakutensoku ktrans wrapper
- Inspect your dependency graph

## 📦 Installation
```go
go install github.com/afeldman/Nishimura@latest
```

Make sure nishimura is available in your $PATH.

---

## ⚙️ Environment

Nishimura interacts with several tools and needs some environment variables:

### Nishimura Home

Global source cache and build artifacts:
```bash
export NISHIMURA_HOME=$HOME/.nishimura
```

### Makoto

Makoto provides the package configuration database (KPC).
```bash
export KPC_HOME=$HOME/.makoto
```

### Gakutensoku (Compiler Wrapper)

To use the Karel compiler via wrapper, ensure ktrans.exe is available in your $PATH.
No server is required – Nishimura calls the local wrapper directly.

## 🚀 Quick Start
1. Create a new project
    ```bash
    mkdir hello_karel && cd hello_karel
    nishimura init
    ```
    This generates a nishimura.kpc:
    ```toml
    kpc_version = "0.2.0"
    name = "hello_karel"
    version = "0.1.0"
    description = "A hello world project in FANUC Karel"
    main = "main.kl"

    [[authors]]
    name = "Your Name"
    email = "you@example.com"
    ```

2. Write your program
    main.kl:
    ```pascal
    PROGRAM hello
    BEGIN
        WRITE('Hello, Karel World!',CR)
    END hello
    ```

3. Add a dependency
    ```bash
    nishimura add motion_lib@1.2.0 https://github.com/yourname/motion_lib.git
    ```

    Dependencies are cloned into ~/.nishimura/src/ and added to your manifest.

4. Install dependencies
    ```bash
        nishimura install
    ```

    - Clones/fetches all deps
    - Registers them in the Makoto DB

5. Show dependency graph
    ```bash
    nishimura graph
    ```

    Example:
    ```graphql
    - hello_karel@0.1.0
        - motion_lib@1.2.0
            - utils@0.3.1
    ```

6. Compile project
    ```bash
    nishimura compile
    ```

    - Collects include paths from dependencies
    - Runs the Gakutensoku ktrans wrapper
    - Outputs to build/:
        ```bash
        build/hello.pc
        ```

7. Manage conflicts

    Mark incompatible versions:
    ```bash
    nishimura conflict add motion_lib@1.0.0
    nishimura conflict check
    ```
---
## 🔑 Commands Overview
|Command|	Description|
|-------|--------------|
|init|	Initialize a new project (nishimura.kpc)|
|add	|Add a dependency (name@version url)|
|install|	Install all dependencies into cache + DB|
|graph|	Print dependency graph|
|compile|	Compile project with ktrans|
|conflict|	Manage conflict rules (add, rm, list, check)|

---

## 🛠️ Roadmap

- Hoffmann integration as package server
- Deployment helper (deploy) for robots
- Smarter conflict resolution strategies
- Precompiled package cache
