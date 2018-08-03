# Nishimura
Nishimura is a Karel package manager. The project is named after *Makoto Nishimura* the engineer of the first Japanese robot.

## Setup

Nishimura can interact with different other projects. Therefore the projects will search for different environment setting.

### Gakutensoku
If the user want to use a global karel compiler, because he uses LINUX or Mac, then the package can send to the **Gakutensoku** server. To set the **Gakutensoku** server connection use

```bash
export GAKUTENSOKU_HOME=127.0.0.1:1234
```
All information on **Gakutensoku** will be available in the projekt description.

### Makoto
Makoto is a project to make a Karel package configuration description. This file will contain all information of the project. The systems runs on the **KPC** shorting. All packages for Karel can be stored into a global dictionary. The path has to be set, if a user want to use dependencies or install the local libray.

```bash
export KPC_HOME=/usr/local/lib/karel/
```

The KPC-System only runs on **.karel** files, the project files.

### Hoffmann
Hoffmann is a package server. This server makes it possible to share libraies. To set the local configuration for the server a configuration will is utilized. The home can set to a path. this example is the default path if no path is set.

```bash
export HOFFMANN_HOME=/etc/hoffmann
```

The Hoffmann setting is stored into a setting.yaml

```yaml
---
servers:
    - 127.0.0.1:123
    - 192.1.2.3:143
```


## Commands

The project utilizes different methods to provide an easy work with Karel projects

|Command|Description|
|-------|-----------|
|build| build the project package. A package uses the mime **.karel**|
|install|install the project package|
|push| push data to the package server|
|pull| pull and install a package|
|get| pull to local project. not install|
|init| build a new project using a project skell|
|deploy| deploy the package to the roboter|
|compile| compile the current package and copy the result into the bin/ folder.|

### init

```bash
Nishimura init
```

The init command will ask the user question to setup the project. A default skell is set so the programmer can start writing all the information

```bash
1. package name (foldername):
2. package version (0.1.0):
3. package description:
4. main file (packagename.kl):
5. paser version (v9.10):
6. repo type (git):
7. repo address:
8. package keywords:
9. author: foo; bar
10. license (MIT):
...
11. OK (yes)?:
```

This question build the kpc file and the compiler file.
The kpc file contains all information on the project. The compiler file will contains all compiler information.
the compiler file is only for this project and will not used in the child project. The child project needs a differnent compiler file.

### build

```bash
Nishimura build <package_path>
```
To build a **.karel** package for the package system this method build a package system. This method does not compile the code. It only checks for the **kpc** file and build the package. The naming is easy *PackageName-version.karel*.

## install

```bash
Nishimura install (--local) [name<@version>]
```

If no name and version is given, then the current project is used to install. To install a project different steps are utilized:

1. build the package
2. send the package to the compiler or use the local installation
3. If the compile process was successfull, then install the package, otherwise break the installation process and give back the error message
4. install the package to a local diectory (./vendor) or global in **Makoto** global folder

If a name is given for this project, then in the first step is the local **Makoto** folder checked and if the project is not available on the local computer, then the package server **Hoffmann** is utilized.
If a version is set, then the system trys to install the requested version. Otherwise the system allways install the latest version.

### push
```bash
Nishimura push
```

This function only starts if **Hoffmann** as package server is available and only if the code is compilable. Then the package is build, compuled and send to the server.

### pull
```bash
Nishimura pull name<@version>
```
This command is an alias to download a requested package. If the package is available, then the package is installed into the **Makoto** global server.

### deploy
```bash
Nishimura deploy robot_ip <device>
```

This command will copy all data in the *bin* folder. The compile process creates the *bin* folder. The default folder on the robot is **MD:\**.

### compile

```bash
Nishimura compile <--gakutensoku> <--test>
```

This command starts the compile process. Because I use mostly Linux, the **Gakutensoku** server is mandatory. This method handed the package to the **KTrans Wrapper** project and builds the package. The test option will delete the bin folder after a successfull build.
