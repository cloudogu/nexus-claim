
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


