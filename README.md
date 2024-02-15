# loadbot 

## Introduction

Loadbot is a workload driver designed to simulate heavy loads on systems for performance testing and benchmarking purposes. It allows users to generate various types of workloads to stress-test their systems under different scenarios.

This code, being an open-source project, is in its early development stage and likely contains various bugs. We welcome contributions from the community to help improve loadbot, make it more robust and reliable, and introduce new features.


## How to Install

Using Homebrew (MacOS/Linux)💡

```bash
brew tap kuzxnia/loadbot
brew install loadbot
```

Alternatively, you can install **loadbot** from sources or run it directly in a Docker container. For more information on these installation methods, please refer to the [documentation](https://kuzxnia.github.io/loadbot/getting_started/install/).


## Getting started
After installing loadbot, you can quickly get started by following these steps:

1. Run LoadBot agent with your desired configuration using the command:
```bash
loadbot start-agent -f config_file.json
```

2. Start the workload using the LoadBot client:
```bash
loadbot start
```

3. Monitor the progress of the workload using the command:
```bash
loadbot progress
```

4. To stop the workload, use the following command:
```bash
loadbot stop
```

For more information and detailed instructions, please refer to the [quick start guide](https://kuzxnia.github.io/loadbot/getting_started/quick-start/).


## Documentation
For detailed documentation on how to use loadbot and its available features, please refer to the [official documentation](https://kuzxnia.github.io/loadbot/).
