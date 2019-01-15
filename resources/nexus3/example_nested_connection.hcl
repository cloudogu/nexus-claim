repository "docker-registry" {
  repositoryName = "docker-registry"
  online = true
  recipeName = "docker-proxy"
  attributes = {
    docker = {
      forceBasicAuth = true
      v1Enabled = false
    }
    dockerProxy = {
      indexType = "HUB"
      indexUrl = "http://test.de"
      useTrustStoreForIndexAccess = true
    }
    proxy = {
      remoteUrl = "https://slm.zd.intranet.bund.de/nexus/repository/public-group-docker"
      contentMaxAge = 1440
      metadataMaxAge = 1440
    }
    httpclient = {
      blocked = false
      autoBlock = true
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
