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
