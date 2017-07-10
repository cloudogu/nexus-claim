# Nexus Claim

Define your [Sonatype Nexus](http://www.sonatype.org/nexus/) repository structure as code.
 

## Example

```hcl
repository "my-hosted-repo" {
  name = "My hosted Repository"
  format = "maven2"
  provider = "maven2"
  providerRole = "org.sonatype.nexus.proxy.repository.Repository"
  repoPolicy = "RELEASE"
  repoType = "hosted"
  exposed = true
  _state = "present"
}

repository "central" {
  name = "Central"
  format = "maven2"
  provider = "maven2"
  providerRole = "org.sonatype.nexus.proxy.repository.Repository"
  repoPolicy = "RELEASE"
  repoType = "proxy"
  remoteUri = "https://repo1.maven.org/maven2/"
  exposed = true
  _state = "present"
}

repository "central-m1" {
  _state = "absent"
}
```
