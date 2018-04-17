import org.sonatype.nexus.blobstore.api.BlobStoreManager
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
def createHostedRepository(String name, String recipeName, String blobStoreName, String writePolicy, String strictContentTypeValidation) {

  def configuration = createHostedConfiguration(name, recipeName, blobStoreName, writePolicy, strictContentTypeValidation)

  try {
    repository.repositoryManager.create(configuration)
  }
  catch (Exception e){
    return e
  }

  //repository.createMavenHosted(name, typedBlobStoreName,typedStrictContentTypeValidation, typedVersionPolicy, typedWritePolicy, typedLayoutPolicy);

}


def createHostedConfiguration(String name, String recipeName, String blobStoreName, String writePolicy, String strictContentTypeValidation){

  def typedBlobStoreName = getblobStoreName(blobStoreName)
  def typedStrictContentTypeValidation = getStrictContentTypeValidation(strictContentTypeValidation)
  def typedWritePolicy = getWritePolicy(writePolicy)

  def configuration = repository.createHosted(name,recipeName,typedBlobStoreName, typedWritePolicy, typedStrictContentTypeValidation)

  return configuration
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

/*
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
*/
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

/*
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
*/
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
def recipeName
def blobStoreName
def strictContentTypeValidation
def writePolicy
/*
def versionPolicy
def layoutPolicy
*/


if (args != ""){
  def value = args.tokenize(' ')
  name = value[0]
  recipeName = value[1]
  blobStoreName = value[2]
  writePolicy = value[3]
  strictContentTypeValidation = value[4]

  this.createHostedRepository(name,recipeName, blobStoreName, writePolicy, strictContentTypeValidation)

}

