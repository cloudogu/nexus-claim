repository "maven-central" {
  _state = "absent"
}

repository "maven-public" {
  _state = "absent"
}

repository "maven-releases" {
  _state = "absent"
}

repository "maven-snapshots" {
  _state = "absent"
}

repository "nuget-group" {
  _state = "absent"
}

repository "nuget-hosted" {
  _state = "absent"
}

repository "nuget.org-proxy" {
  _state = "absent"
}

repository "3rdparty" {
  recipeName = "maven2-hosted"
  online = true
  attributes = {
    cleanup = {
      policyName = "None"
    },
    storage = {
      blobStoreName = "default"
      writePolicy = "ALLOW"
      strictContentTypeValidation = true
    }
    maven = {
      versionPolicy = "RELEASE"
      layoutPolicy = "STRICT"
    }
  }
  _state = "present"
}

repository "yourcompshared" {
  recipeName = "maven2-proxy"
  online = true
  attributes = {
    cleanup = {
      policyName = "None"
    },
    httpclient = {
      blocked = false
      autoBlock = true
      connection = {
        useTrustStore=false
      }
    }
    proxy = {
      remoteUrl= "https://yoururl.com/nexus/content/groups/public/"
      contentMaxAge = -1
      metadataMaxAge = 1440
    }
    negativeCache = {
      enabled = true
      timeToLive = 1440
      maven-indexer = {}
    }
    storage = {
      blobStoreName = "default"
      strictContentTypeValidation = true
    }
    maven = {
      versionPolicy = "RELEASE"
      layoutPolicy = "PERMISSIVE"
    }
  }
  _state = "present"
}

repository "public" {
  online = true
  recipeName = "maven2-group"
  attributes = {
    cleanup = {
      policyName = "None"
    },
    group = {
      memberNames = ["3rdparty"]
    }
    storage = {
      blobStoreName = "default"
      strictContentTypeValidation = true
    }
    maven = {
      versionPolicy = "RELEASE"
      layoutPolicy = "PERMISSIVE"
    }
  }
  _state = "present"
}

repository "docker-registry" {
  recipeName = "docker-proxy"
  online = true
  attributes = {
    cleanup = {
      policyName = "None"
    },
    docker = {
      forceBasicAuth = true
      v1Enabled = false
    }
    proxy = {
      remoteUrl= "https://yoururl.com/nexus/repository/public-group-docker"
      contentMaxAge = 1440
      metadataMaxAge = 1440
    }
    dockerProxy = {
      indexType = "REGISTRY"
    }
    httpclient = {
      blocked=false
      autoBlock=true
      connection = {
        useTrustStore = false
      }
    }
    storage = {
      blobStoreName = "default"
      strictContentTypeValidation = true
    }
    negativeCache = {
      enabled = true
      timeToLive = 1440
    }
  }
  _state = "present"
}
