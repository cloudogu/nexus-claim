# Nexus Claim

Define your [Sonatype Nexus](http://www.sonatype.org/nexus/) repository structure as code.


## Example uses of nexus-claim on nexus 3
```bash
$ nexus-claim plan -i resources/nexus3/nexus3-initial-example.hcl -o resources/nexus3/nexus3-initial-example.json
$ nexus-claim apply -i resources/nexus3/nexus3-initial-example.json
```

## Things to consider for .hcl of nexus 3
- `repository "xxx"` should be the same as `name = "xxx"`
- on maven2-hosted- and maven2-proxy-repository there must be a maven sector in addition(with versionPolicy and writePolicy) 
 
## Example .hcl for nexus 3
```hcl

repository "testGroup" {
  name = "testGroup"
  online = true
  recipeName = "bower-group"
  attributes = {
      group = {
        memberNames = ["maven-public","maven-central"]
      }
    storage = {
      blobStoreName = "default"
    }
  }
  _state = "present"
}

repository "testProxy" {
  name = "testProxy"
  online = true
  recipeName = "maven2-proxy"

  attributes = {
    httpclient = {
      connection = {
        blocked = false
        autoBlock = true
      }
    }
    proxy = {
      remoteUrl= "http://www.test.de"
      contentMaxAge = 1440
      metadataMaxAge = 1440
    }
    negativeCache = {
      enabled = true
      timeToLive = 1440
    }
    storage = {
      blobStoreName = "default"
      strictContentTypeValidation = false
    }
    maven = {
      versionPolicy = "RELEASE"
      layoutPolicy = "PERMISSIVE"
    }
  }

  _state = "present"
}


repository "testHosted" {
  name = "testHosted"
  online = true
  recipeName = "maven2-hosted"
  attributes = {
    storage = {
      blobStoreName = "default"
      writePolicy = "ALLOW"
      strictContentTypeValidation = false
    }
    maven = {
      versionPolicy = "RELEASE"
      layoutPolicy = "PERMISSIVE"
    }
  }
  _state = "present"
}

repository "deleteMe" {
  name = "deleteMe"
  _state = "absent"
}

 ```

## Example .hcl for nexus 2

```hcl
repository "apache-snapshots" {
  _state = "absent"
}

repository "central-m1" {
  _state = "absent"
}

repository "thirdparty" {
  name = "Third Party"
  _state = "present"
}

repository "scm-releases" {
  name = "SCM-Manager Releases"
  format = "maven2"
  provider = "maven"
  artifactMaxAge = -1
  autoBlockActive =  true
  browseable =  true
  checksumPolicy = "WARN"
  downloadRemoteIndexes = true
  exposed = true
  fileTypeValidation = true
  indexable = true
  itemMaxAge = 1440
  metadataMaxAge = 1440
  notFoundCacheTTL = 1440
  providerRole = "org.sonatype.nexus.proxy.repository.Repository"
  remoteStorage = {
    remoteStorageUrl = "https://maven.scm-manager.org/nexus/content/repositories/releases/"
  }
  repoPolicy = "RELEASE"
  repoType = "proxy"
  writePolicy = "READ_ONLY"
  _state = "present"
}


```



