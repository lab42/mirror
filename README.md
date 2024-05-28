![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/lab42/mirror/CD.yaml)
![Codecov (with branch)](https://img.shields.io/codecov/c/github/lab42/mirror/main)

# Mirror

Mirror is a tool for reflecting Kubernetes resources such as Secrets and ConfigMaps across namespaces.

## Usage

Mirror provides a command-line interface (CLI) for reflecting Kubernetes resources. You can run the container or dev workflow with 'air'

By default, Mirror reflects both ConfigMaps and Secrets. You can configure the reflection behavior using a configuration file or environment variables.

## Configuration

Mirror supports configuration through YAML files. Here's an example configuration file (default settings):

```yaml
logLevel: debug

reflect:
  configMap:
    enabled: true
    annotation: reflect.lab42.io/namespaces

  secret:
    enabled: true
    annotation: reflect.lab42.io/namespaces

kubeconfig:
  inCluster: true
  path: ""
```

If you want to use a configuration file then you must set the environment variable `CONFIG_FILE` to the path where your config file is located. This must be a full path including extension(yml/yaml).

You can customize the log level, reflection settings for ConfigMaps and Secrets, and Kubernetes configuration using this YAML file.

### Environment Variables

Mirror can also be configured using environment variables. Here are the environment variables that can be set:

- `MIRROR_LOGLEVEL`: Sets the log level (e.g., `"debug"`, `"info"`, `"warn"`, `"error"`)
- `MIRROR_REFLECT_CONFIGMAP_ENABLED`: Sets whether reflection for ConfigMaps is enabled (e.g., `"true"`, `"false"`)
- `MIRROR_REFLECT_CONFIGMAP_ANNOTATION`: Sets the annotation used for ConfigMap reflection
- `MIRROR_REFLECT_SECRET_ENABLED`: Sets whether reflection for Secrets is enabled (e.g., `"true"`, `"false"`)
- `MIRROR_REFLECT_SECRET_ANNOTATION`: Sets the annotation used for Secret reflection
- `MIRROR_KUBECONFIG_INCLUSTER`: Sets whether to use in-cluster Kubernetes configuration
- `MIRROR_KUBECONFIG_PATH`: Sets the path to the Kubernetes configuration file (for out of cluster development)

Environment variables should be prefixed with `MIRROR_`. For example, to set the log level to `"debug"`, you would use `MIRROR_LOGLEVEL=debug`.

## Contributing

Contributions to Mirror are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on GitHub.

## License

Mirror is licensed under the [DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE](LICENSE).
