import org.sonatype.nexus.blobstore.api.BlobStore
import org.sonatype.nexus.blobstore.api.BlobStoreManager
import org.sonatype.nexus.common.log.LogConfigurationCustomizer
import org.sonatype.nexus.repository.config.Configuration
import org.sonatype.nexus.repository.maven.LayoutPolicy
import org.sonatype.nexus.repository.maven.VersionPolicy
import org.sonatype.nexus.repository.storage.WritePolicy

/*
repository.createMavenHosted('private-again',
  BlobStoreManager.DEFAULT_BLOBSTORE_NAME,
  true,
  VersionPolicy.RELEASE,
  WritePolicy.ALLOW_ONCE,
  LayoutPolicy.STRICT)
*/
def createMavenRepository(String name, String blobStoreName, String strictContentTypeValidation, String versionPolicy, String writePolicy, String layoutPolicy) {

  def typedBlobStoreName = getblobStoreName(blobStoreName)
  def typedStrictContentTypeValidation = getStrictContentTypeValidation(strictContentTypeValidation)
  def typedVersionPolicy = getVersionPolicy(versionPolicy)
  def typedWritePolicy = getWritePolicy(writePolicy)
  def typedLayoutPolicy = getLayoutPolicy(layoutPolicy)


  repository.createMavenHosted(name, typedBlobStoreName,typedStrictContentTypeValidation, typedVersionPolicy, typedWritePolicy, typedLayoutPolicy);
  repository.repositoryManager.update(new Configuration()

  LogConfigurationCustomizer.Configuration
}


def getStrictContentTypeValidation(String strictContentTypeValidation){
  if (strictContentTypeValidation.equals("true")) {
    return true
  }
  return false
}

def getblobStoreName(String blobStoreName){
  if(blobStoreName.equals("null") || blobStoreName == null){
    return BlobStoreManager.DEFAULT_BLOBSTORE_NAME
  }
  return blobStoreName
}

def getVersionPolicy(String versionPolicy){
  switch (versionPolicy.toLowerCase()){
    case "mixed":
      return VersionPolicy.MIXED
      break;
    case "snapshop":
      return VersionPolicy.SNAPSHOT
      break;
    default:
      return VersionPolicy.RELEASE
  }
}

def getWritePolicy(String writePolicy) {
  switch (writePolicy.toLowerCase()) {
    case "allow_once":
      return WritePolicy.ALLOW_ONCE
      break;
    case "deny":
      return WritePolicy.DENY
      break;
    default:
      return WritePolicy.ALLOW
  }
}

def getLayoutPolicy(String layoutPolicy){
  switch (layoutPolicy.toLowerCase()) {
    case "permissive":
      return LayoutPolicy.PERMISSIVE
      break;
    default:
      return LayoutPolicy.STRICT
      break;
  }
}

/*
def get

enum Repository{
  BOWER,
  DOCKER,
  GITLFS,
  MAVEN2,
  NPM,
  NUGET,
  PYPI,
  RAW,
  RUBYGEMS,
  YUM;
}

enum RepositoryType {
  HOSTED,
  GROUP,
  PROXY;
}
*/
def name
def blobStoreName
def strictContentTypeValidation
def versionPolicy
def writePolicy
def layoutPolicy


if (args != ""){
  def value = args.tokenize(' ')
  name = value[0]
  blobStoreName = value[1]
  strictContentTypeValidation = value[2]
  versionPolicy = value[3]
  writePolicy= value[4]
  layoutPolicy= value[5]

  this.createMavenRepository(name,blobStoreName,strictContentTypeValidation,versionPolicy,writePolicy,layoutPolicy)

}

