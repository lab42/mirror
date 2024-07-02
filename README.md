# Mirror

![GitHub Release](https://img.shields.io/github/v/release/lab42/mirror)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/lab42/mirror/CD.yaml)
![Codecov (with branch)](https://img.shields.io/codecov/c/github/lab42/mirror/main)

Mirror is a tool for reflecting Kubernetes resources such as Secrets and ConfigMaps across namespaces.

## Why Another Reflector?

Reflectors manage sensitive data, making it crucial to minimize attack surfaces. To ensure security and reliability, it should use the native Kubernetes libraries (written in Golang) and be robust and dependable.

## Security Benefits

The security landscape in containerized environments need a minimalist approach. This is why our reflector is statically compiled into a scratch container, ensuring that only the essential binary is included. This drastically reduces the attack surface, eliminating potential vulnerabilities associated with larger base images, making our reflector a robust solution for sensitive configurations and secrets.

## Why You Should Use It

Incorporating Mirror into your Kubernetes setup brings a couple of benefits:

Enhanced Security: With a scratch container, you eliminate unnecessary components, reducing potential security risks.
High Performance: Golang's efficient concurrency model ensures that the reflector can handle high loads and scale seamlessly with your Kubernetes environment.
Reliability: Statically compiled binaries mean fewer dependencies and a lower chance of runtime failures, ensuring your configuration and secrets management is always up and running.

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

If you want to use a configuration file then you must set the environment variable `CONFIG_FILE` to the path where your config file is located. This must be an absolute path including extension(yml/yaml).

You can customize the log level, reflection settings for ConfigMaps and Secrets, and Kubernetes configuration using this YAML file.

The default configuration can be printed using the following command `mirror config`.

### Environment Variables

Mirror can also be configured using environment variables. Here are the environment variables that can be set:

- `MIRROR_LOGLEVEL`: Sets the log level (e.g., `"debug"`, `"info"`, `"warn"`, `"error"`)
- `MIRROR_REFLECT_CONFIGMAP_ENABLED`: Sets whether reflection for ConfigMaps is enabled (e.g., `"true"`, `"false"`)
- `MIRROR_REFLECT_CONFIGMAP_ANNOTATION`: Sets the annotation used for ConfigMap reflection
- `MIRROR_REFLECT_SECRET_ENABLED`: Sets whether reflection for Secrets is enabled (e.g., `"true"`, `"false"`)
- `MIRROR_REFLECT_SECRET_ANNOTATION`: Sets the annotation used for Secret reflection
- `MIRROR_KUBECONFIG_INCLUSTER`: Sets whether to use in-cluster Kubernetes configuration

Environment variables should be prefixed with `MIRROR_`. For example, to set the log level to `"debug"`, you would use `MIRROR_LOGLEVEL=debug`.

## Contributing

We thrive on community collaboration and believe in the power of collective improvement. Instead of forking the code, we encourage you to contribute directly to our project by making feature requests or pull requests. Your ideas and enhancements can benefit everyone and help us create a more robust and versatile tool.

By submitting feature requests, you ensure that your needs are considered in future updates. Making pull requests allows your improvements to be integrated seamlessly, benefiting the entire community while ensuring compatibility and avoiding fragmentation.

Join us in building a stronger, more secure ConfigMap and Secret Reflector. Let's innovate together and make a greater impact!

## License

Mirror is licensed under the [DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE](LICENSE).
