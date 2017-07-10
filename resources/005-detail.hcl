repository "releases" {
  Name = "Releases"
  Format = "maven2"
  RepoType = "hosted"
  RepoPolicy = "RELEASE"
  Provider = "maven2"
  Browseable = true
  Indexable = true
  Exposed = true
  DownloadRemoteIndexes = false
  WritePolicy = "ALLOW_WRITE_ONCE"
  ProviderRole = "org.sonatype.nexus.proxy.repository.Repository"
  NotFoundCacheTTL = 1440
}
