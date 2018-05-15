repository "scm-manager-releases" {
  name = "scm-manager-releases"
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
      remoteUrl= "https://maven.scm-manager.org/nexus/content/repositories/releases/"
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

repository "scm-manager-snapshots" {
  name = "scm-manager-snapshots"
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
      remoteUrl= "https://maven.scm-manager.org/nexus/content/repositories/snapshots/"
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

repository_group "scm-manager" {
  name = "SCM-Manager"
  online = true
  recipeName = "maven-group"
  attributes = {
      group = {
        memberNames = ["scm-manager-releases","scm-manager-snapshots"]
      }
    storage = {
      blobStoreName = "default"
    }
  }
  _state = "present"
}

repository "third-party" {
  name = "3rd Party"
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

repository_group "nuget-group" {
  _state = "absent"
}

repository "nuget-hosted" {
  _state = "absent"
}

repository "nuget.org-proxy" {
  _state = "absent"
}
