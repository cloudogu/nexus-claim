repository "testGroup" {
  repositoryName = "testGroup"
  online = true
  recipeName = "bower-group"
  attributes = {
    group = {
      memberNames = [
        "maven-public",
        "maven-central"]
    }
    storage = {
      blobStoreName = "default"
    }
  }
  _state = "present"
}


repository "bowerProxy" {
  repositoryName = "bowerProxy"
  online = true
  recipeName = "bower-proxy"

  attributes = {
    httpclient = {
      connection = {
        blocked = false
        autoBlock = true
      }
    }
    proxy = {
      remoteUrl = "http://www.test.de"
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
    bower = {
      rewritePackageUrls = true
    }
  }

  _state = "present"
}

repository "docker-registry2" {
  repositoryName = "docker-registry2"
  online = true
  recipeName = "docker-hosted"
  _state = "present"

  attributes = {
    docker = {
      v1Enabled = false
      forceBasicAuth = true
      httpPort = 222
      httpsPort = 122
    }

    storage = {
      blobStoreName = "default"
      writePolicy = "ALLOW"
      strictContentTypeValidation = true
    }
  }
}

repository "docker-registry1" {
  repositoryName = "docker-registry1"
  online = true
  recipeName = "docker-proxy"
  _state = "present"

  attributes = {
    docker = {
      forceBasicAuth = true
      v1Enabled = false
      httpPort = 2
      httpsPort = 1
    }

    dockerProxy = {
      indexType = "HUB"
      indexUrl = "http://www.test.de"
      useTrustStoreForIndexAccess = false
    }

    httpclient = {
      connection = {
        blocked = false
        autoBlock = true
        useTrustStore = false
      }
    }

    negativeCache = {
      enabled = true
      timeToLive = 1440
    }

    proxy = {
      contentMaxAge = 1440
      metadataMaxAge = 1440
      remoteUrl = "https://slm.zd.intranet.bund.de/nexus/repository/public-group-docker"
    }

    storage = {
      blobStoreName = "default"
      strictContentTypeValidation = true
    }
  }
}

repository "testProxy" {
  repositoryName = "testProxy"
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
      remoteUrl = "http://www.test.de"
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
  repositoryName = "testHosted"
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

repository "yumhosted" {
  repositoryName = "yumhosted"
  online = true
  recipeName = "yum-hosted"
  attributes = {
    storage = {
      blobStoreName = "default"
      writePolicy = "ALLOW"
      strictContentTypeValidation = false
    }
    yum = {
      deployPolicy = "PERMISSIVE",
      repodataDepth = 2
    }
  }
  _state = "present"
}


repository "deleteMe" {
  repositoryName = "deleteMe"
  _state = "absent"
}

repository "docker-registry3" {
  repositoryName = "docker-registry3"
  online = true
  recipeName = "docker-group"
  _state = "present"

  attributes = {
    docker = {
      v1Enabled = false
      forceBasicAuth = true
      httpPort = 33
      httpsPort = 34
    }
    group = {
      memberNames = [
        "docker-registry1",
        "docker-registry2"]
    }
    storage = {
      blobStoreName = "default"
    }

  }
}


